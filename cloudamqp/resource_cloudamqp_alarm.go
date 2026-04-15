package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &alarmResource{}
	_ resource.ResourceWithConfigure   = &alarmResource{}
	_ resource.ResourceWithImportState = &alarmResource{}
)

type alarmResource struct {
	client *api.API
}

func NewAlarmResource() resource.Resource {
	return &alarmResource{}
}

type alarmResourceModel struct {
	ID               types.String `tfsdk:"id"`
	InstanceID       types.Int64  `tfsdk:"instance_id"`
	Type             types.String `tfsdk:"type"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	ReminderInterval types.Int64  `tfsdk:"reminder_interval"`
	ValueThreshold   types.Int64  `tfsdk:"value_threshold"`
	ValueCalculation types.String `tfsdk:"value_calculation"`
	TimeThreshold    types.Int64  `tfsdk:"time_threshold"`
	VhostRegex       types.String `tfsdk:"vhost_regex"`
	QueueRegex       types.String `tfsdk:"queue_regex"`
	MessageType      types.String `tfsdk:"message_type"`
	Recipients       types.List   `tfsdk:"recipients"`
}

func (r *alarmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_alarm"
}

func (r *alarmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to create and manage alarms to trigger based on a set of conditions. " +
			"Once triggered a notification will be sent to the assigned recipients.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "The instance identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of the alarme",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cpu",
						"memory",
						"disk",
						"queue",
						"connection",
						"flow",
						"consumer",
						"netsplit",
						"ssh",
						"notice",
						"server_unreachable",
					),
				},
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "Enable or disable an alarm",
			},
			"reminder_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The reminder interval (in seconds) to resend the alarm if not resolved. Set to 0 for no reminders.",
				Default:     int64default.StaticInt64(0),
			},
			"value_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "What value to trigger the alarm for",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"value_calculation": schema.StringAttribute{
				Optional:    true,
				Description: "Disk value threshold calculation. Fixed or percentage of disk space remaining",
				Validators: []validator.String{
					stringvalidator.OneOf("fixed", "percentage"),
				},
			},
			"time_threshold": schema.Int64Attribute{
				Optional:    true,
				Description: "For how long (in seconds) the value_threshold should be active before trigger alarm",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"vhost_regex": schema.StringAttribute{
				Optional:    true,
				Description: "Regex for which vhost the queues are in",
			},
			"queue_regex": schema.StringAttribute{
				Optional:    true,
				Description: "Regex for which queues to check",
			},
			"message_type": schema.StringAttribute{
				Optional:    true,
				Description: "Message types of the queue to trigger the alarm",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"total", "unacked", "ready", "ack", "ack_rate", "deliver", "deliver_rate", "deliver_get",
						"deliver_get_rate", "deliver_no_ack", "deliver_no_ack_rate", "get", "get_rate",
						"get_empty", "get_empty_rate", "get_no_ack", "get_no_ack_rate", "publish", "publish_rate",
						"redeliver", "redeliver_rate"),
				},
			},
			"recipients": schema.ListAttribute{
				ElementType: types.Int64Type,
				Required:    true,
				Description: "Identifiers for recipients to be notified.",
			},
		},
	}
}

func (r *alarmResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *alarmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("ImportState: ID=%s", req.ID))
	if !strings.Contains(req.ID, ",") {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {alarm_id},{instance_id}")
		return
	}

	idSplit := strings.Split(req.ID, ",")
	if len(idSplit) != 2 {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {alarm_id},{instance_id}")
		return
	}
	instanceID, err := strconv.ParseInt(idSplit[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance_id in import ID", fmt.Sprintf("Could not convert instance_id to int: %s", err))
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), idSplit[0])
	resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)
	// Default values for computed attributes
	resp.State.SetAttribute(ctx, path.Root("reminder_interval"), 0)
}

func (r *alarmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alarmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	params := r.populateRequest(ctx, plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var alarmID string
	if params.Type == "notice" {
		tflog.Info(ctx, "alarm type is 'notice', skip creation, retrieve existing alarm and update")
		alarms, err := r.client.ListAlarms(timeoutCtx, instanceID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to List Alarms",
				fmt.Sprintf("Could not list alarms to find 'notice' alarm: %s", err),
			)
			return
		}

		for _, alarm := range alarms {
			if alarm.Type == "notice" {
				alarmID = fmt.Sprintf("%d", alarm.ID)
				break
			}
		}

		if alarmID == "" {
			resp.Diagnostics.AddError(
				"Notice Alarm Not Found",
				"Could not find existing 'notice' alarm to update.",
			)
			return
		}

		err = r.client.UpdateAlarm(timeoutCtx, instanceID, alarmID, params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Update Alarm",
				fmt.Sprintf("Could not update alarm: %s", err),
			)
			return
		}
	} else {
		var err error
		alarmID, err = r.client.CreateAlarm(timeoutCtx, instanceID, params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Create Alarm",
				fmt.Sprintf("Could not create alarm: %s", err),
			)
			return
		}
	}

	plan.ID = types.StringValue(alarmID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *alarmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alarmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	alarmID := state.ID.ValueString()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	data, err := r.client.ReadAlarm(timeoutCtx, instanceID, alarmID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Alarm",
			fmt.Sprintf("Could not read alarm: %s", err),
		)
		return
	}

	if data == nil {
		tflog.Warn(ctx, fmt.Sprintf("Resource drift detected for alarm ID %s or instance ID %d", alarmID, instanceID))
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateResourceModel(ctx, *data, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *alarmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alarmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := r.populateRequest(ctx, plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	err := r.client.UpdateAlarm(timeoutCtx, plan.InstanceID.ValueInt64(), plan.ID.ValueString(), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Alarm",
			fmt.Sprintf("Could not update alarm: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *alarmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alarmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Type.ValueString() == "notice" {
		tflog.Debug(ctx, "alarm type is 'notice', skip deletion and just remove from state")
		resp.State.RemoveResource(ctx)
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	alarmID := state.ID.ValueString()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	err := r.client.DeleteAlarm(timeoutCtx, instanceID, alarmID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Alarm",
			fmt.Sprintf("Could not delete alarm: %s", err),
		)
		return
	}
}

func (r *alarmResource) populateRequest(ctx context.Context, plan alarmResourceModel) model.AlarmRequest {
	params := model.AlarmRequest{}
	params.Type = plan.Type.ValueString()
	params.Enabled = plan.Enabled.ValueBool()

	if !plan.ReminderInterval.IsUnknown() && !plan.ReminderInterval.IsNull() {
		params.ReminderInterval = plan.ReminderInterval.ValueInt64Pointer()
	}

	if !plan.ValueThreshold.IsUnknown() && !plan.ValueThreshold.IsNull() {
		params.ValueThreshold = plan.ValueThreshold.ValueInt64Pointer()
	}
	if !plan.ValueCalculation.IsNull() {
		params.ValueCalculation = plan.ValueCalculation.ValueString()
	}
	if !plan.TimeThreshold.IsUnknown() && !plan.TimeThreshold.IsNull() {
		params.TimeThreshold = plan.TimeThreshold.ValueInt64Pointer()
	}
	if !plan.VhostRegex.IsUnknown() && !plan.VhostRegex.IsNull() {
		params.VhostRegex = plan.VhostRegex.ValueString()
	}
	if !plan.QueueRegex.IsUnknown() && !plan.QueueRegex.IsNull() {
		params.QueueRegex = plan.QueueRegex.ValueString()
	}
	if !plan.MessageType.IsUnknown() && !plan.MessageType.IsNull() {
		params.MessageType = plan.MessageType.ValueString()
	}
	if !plan.Recipients.IsUnknown() && !plan.Recipients.IsNull() {
		var recipients = []int64{}
		plan.Recipients.ElementsAs(ctx, &recipients, false)
		params.Recipients = utils.Pointer(recipients)
	}

	return params
}

func (a *alarmResource) populateResourceModel(ctx context.Context, data model.AlarmResponse, state *alarmResourceModel) {
	state.Type = types.StringValue(data.Type)
	state.Enabled = types.BoolValue(data.Enabled)

	if data.ReminderInterval != nil {
		state.ReminderInterval = types.Int64Value(*data.ReminderInterval)
	} else {
		state.ReminderInterval = types.Int64Null()
	}

	if data.ValueThreshold != nil {
		state.ValueThreshold = types.Int64Value(*data.ValueThreshold)
	} else {
		state.ValueThreshold = types.Int64Null()
	}

	if data.ValueCalculation != nil {
		state.ValueCalculation = types.StringValue(*data.ValueCalculation)
	} else {
		state.ValueCalculation = types.StringNull()
	}

	if data.TimeThreshold != nil {
		state.TimeThreshold = types.Int64Value(*data.TimeThreshold)
	} else {
		state.TimeThreshold = types.Int64Null()
	}

	if data.VhostRegex != nil {
		state.VhostRegex = types.StringValue(*data.VhostRegex)
	} else {
		state.VhostRegex = types.StringNull()
	}

	if data.QueueRegex != nil {
		state.QueueRegex = types.StringValue(*data.QueueRegex)
	} else {
		state.QueueRegex = types.StringNull()
	}

	if data.MessageType != nil {
		state.MessageType = types.StringValue(*data.MessageType)
	} else {
		state.MessageType = types.StringNull()
	}

	if data.Recipients != nil {
		recipientsList, _ := types.ListValueFrom(ctx, types.Int64Type, *data.Recipients)
		state.Recipients = recipientsList
	} else {
		state.Recipients = types.ListNull(types.Int64Type)
	}
}
