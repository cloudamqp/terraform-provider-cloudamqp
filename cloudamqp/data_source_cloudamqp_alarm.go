package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &alarmDataSource{}
	_ datasource.DataSourceWithConfigure = &alarmDataSource{}
)

type alarmDataSource struct {
	client *api.API
}

func NewAlarmDataSource() datasource.DataSource {
	return &alarmDataSource{}
}

type alarmDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	InstanceID       types.Int64  `tfsdk:"instance_id"`
	AlarmID          types.Int64  `tfsdk:"alarm_id"`
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

func (d *alarmDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_alarm"
}

func (d *alarmDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:        "Use this data source to retrieve information about default or created alarms. Either use alarm_id or type to retrieve the alarm.",
		DeprecationMessage: "Use 'cloudamqp_alarms' data source instead. This data source will be removed in a future version.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this data source",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
			},
			"alarm_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Alarm identifier",
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.MatchRoot("type")),
				},
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Type of the alarm",
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
					stringvalidator.ConflictsWith(path.MatchRoot("alarm_id")),
				},
			},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Enable or disable an alarm",
			},
			"reminder_interval": schema.Int64Attribute{
				Computed:    true,
				Description: "The reminder interval (in seconds) to resend the alarm if not resolved. Set to 0 for no reminders",
			},
			"value_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: "What value to trigger the alarm for",
			},
			"value_calculation": schema.StringAttribute{
				Computed:    true,
				Description: "Disk value threshold calculation. Fixed or percentage of disk space remaining",
			},
			"time_threshold": schema.Int64Attribute{
				Computed:    true,
				Description: "For how long (in seconds) the value_threshold should be active before trigger alarm",
			},
			"vhost_regex": schema.StringAttribute{
				Computed:    true,
				Description: "Regex for which vhost the queues are in",
			},
			"queue_regex": schema.StringAttribute{
				Computed:    true,
				Description: "Regex for which queues to check",
			},
			"message_type": schema.StringAttribute{
				Computed:    true,
				Description: "Message types (total, unacked, ready) of the queue to trigger the alarm",
			},
			"recipients": schema.ListAttribute{
				ElementType: types.Int64Type,
				Computed:    true,
				Description: "Identifiers for recipients to be notified.",
			},
		},
	}
}

func (d *alarmDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *alarmDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config alarmDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.AlarmID.IsNull() && config.Type.IsNull() {
		resp.Diagnostics.AddError(
			"Either alarm_id or type must be specified",
			"Please specify at least one of the attributes to lookup the alarm.",
		)
		return
	}

	instanceID := config.InstanceID.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	if !config.AlarmID.IsNull() {
		alarmID := fmt.Sprintf("%d", config.AlarmID.ValueInt64())
		alarm, err := d.client.ReadAlarm(timeoutCtx, instanceID, alarmID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Read Alarm by ID",
				fmt.Sprintf("Could not read alarm with ID %s for instance %d: %s", alarmID, instanceID, err.Error()),
			)
			return
		}
		d.populateResourceModel(ctx, *alarm, &config)
	} else if !config.Type.IsNull() {
		alarmType := config.Type.ValueString()
		alarms, err := d.client.ListAlarms(timeoutCtx, instanceID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Read Alarm by Type",
				fmt.Sprintf("Could not read alarm with type %s for instance %d: %s", alarmType, instanceID, err.Error()),
			)
			return
		}

		for _, alarm := range alarms {
			if alarm.Type == alarmType {
				d.populateResourceModel(ctx, alarm, &config)
				break
			}
		}
	}

	config.ID = types.StringValue(fmt.Sprintf("%d", config.AlarmID.ValueInt64()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (d *alarmDataSource) populateResourceModel(ctx context.Context, data model.AlarmResponse, config *alarmDataSourceModel) {
	config.AlarmID = types.Int64Value(int64(data.ID))
	config.Type = types.StringValue(data.Type)
	config.Enabled = types.BoolValue(data.Enabled)

	if data.ReminderInterval != nil {
		config.ReminderInterval = types.Int64Value(*data.ReminderInterval)
	} else {
		config.ReminderInterval = types.Int64Null()
	}

	if data.ValueThreshold != nil {
		config.ValueThreshold = types.Int64Value(*data.ValueThreshold)
	} else {
		config.ValueThreshold = types.Int64Null()
	}

	if data.ValueCalculation != nil {
		config.ValueCalculation = types.StringValue(*data.ValueCalculation)
	} else {
		config.ValueCalculation = types.StringNull()
	}

	if data.TimeThreshold != nil {
		config.TimeThreshold = types.Int64Value(*data.TimeThreshold)
	} else {
		config.TimeThreshold = types.Int64Null()
	}

	if data.VhostRegex != nil {
		config.VhostRegex = types.StringValue(*data.VhostRegex)
	} else {
		config.VhostRegex = types.StringNull()
	}

	if data.QueueRegex != nil {
		config.QueueRegex = types.StringValue(*data.QueueRegex)
	} else {
		config.QueueRegex = types.StringNull()
	}

	if data.MessageType != nil {
		config.MessageType = types.StringValue(*data.MessageType)
	} else {
		config.MessageType = types.StringNull()
	}

	if data.Recipients != nil {
		recipientsList, _ := types.ListValueFrom(ctx, types.Int64Type, *data.Recipients)
		config.Recipients = recipientsList
	} else {
		config.Recipients = types.ListNull(types.Int64Type)
	}
}
