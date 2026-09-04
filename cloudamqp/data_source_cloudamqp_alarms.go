package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &alarmsDataSource{}
	_ datasource.DataSourceWithConfigure = &alarmsDataSource{}
)

type alarmsDataSource struct {
	client *api.API
}

func NewAlarmsDataSource() datasource.DataSource {
	return &alarmsDataSource{}
}

type alarmsDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	Type       types.String `tfsdk:"type"`
	Alarms     types.List   `tfsdk:"alarms"`
}

func (d *alarmsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_alarms"
}

func (d *alarmsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	alarmObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"alarm_id":          types.Int64Type,
		"type":              types.StringType,
		"enabled":           types.BoolType,
		"reminder_interval": types.Int64Type,
		"value_threshold":   types.Int64Type,
		"value_calculation": types.StringType,
		"time_threshold":    types.Int64Type,
		"vhost_regex":       types.StringType,
		"queue_regex":       types.StringType,
		"message_type":      types.StringType,
		"recipients":        types.ListType{ElemType: types.Int64Type},
	}}

	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve alarms for an instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source identifier",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Type of the alarm",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cpu", "memory", "disk", "disk_auto_resize", "queue", "connection",
						"flow", "consumer", "netsplit", "ssh", "notice", "server_unreachable",
					),
				},
			},
			"alarms": schema.ListAttribute{
				Computed:    true,
				ElementType: alarmObjectType,
				Description: "List of alarms",
			},
		},
	}
}

func (d *alarmsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *alarmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state alarmsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	data, err := d.client.ListAlarms(ctx, instanceID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to List Alarms", err.Error())
		return
	}

	filtered := data
	if !state.Type.IsNull() && !state.Type.IsUnknown() && state.Type.ValueString() != "" {
		alarmType := state.Type.ValueString()
		state.ID = types.StringValue(fmt.Sprintf("%d.%s.alarms", instanceID, alarmType))
		filtered = make([]model.AlarmResponse, 0)
		for _, alarm := range data {
			if alarm.Type == alarmType {
				filtered = append(filtered, alarm)
			}
		}
	} else {
		state.ID = types.StringValue(fmt.Sprintf("%d.alarms", instanceID))
	}

	alarmObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"alarm_id":          types.Int64Type,
		"type":              types.StringType,
		"enabled":           types.BoolType,
		"reminder_interval": types.Int64Type,
		"value_threshold":   types.Int64Type,
		"value_calculation": types.StringType,
		"time_threshold":    types.Int64Type,
		"vhost_regex":       types.StringType,
		"queue_regex":       types.StringType,
		"message_type":      types.StringType,
		"recipients":        types.ListType{ElemType: types.Int64Type},
	}}

	values := make([]attr.Value, 0, len(filtered))
	for _, a := range filtered {
		obj, diags := alarmObjectValue(ctx, alarmObjectType.AttrTypes, a)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		values = append(values, obj)
	}

	alarmsList, diags := types.ListValue(alarmObjectType, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Alarms = alarmsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func alarmObjectValue(ctx context.Context, attrTypes map[string]attr.Type, data model.AlarmResponse) (types.Object, diag.Diagnostics) {
	reminder := int64(0)
	if data.ReminderInterval != nil {
		reminder = *data.ReminderInterval
	}
	valueThreshold := int64(0)
	if data.ValueThreshold != nil {
		valueThreshold = *data.ValueThreshold
	}
	valueCalculation := ""
	if data.ValueCalculation != nil {
		valueCalculation = *data.ValueCalculation
	}
	timeThreshold := int64(0)
	if data.TimeThreshold != nil {
		timeThreshold = *data.TimeThreshold
	}
	vhostRegex := ""
	if data.VhostRegex != nil {
		vhostRegex = *data.VhostRegex
	}
	queueRegex := ""
	if data.QueueRegex != nil {
		queueRegex = *data.QueueRegex
	}
	messageType := ""
	if data.MessageType != nil {
		messageType = *data.MessageType
	}
	recipients := []int64{}
	if data.Recipients != nil {
		recipients = *data.Recipients
	}
	recipientsList, diags := types.ListValueFrom(ctx, types.Int64Type, recipients)
	if diags.HasError() {
		return types.Object{}, diags
	}

	obj, objDiags := types.ObjectValue(attrTypes, map[string]attr.Value{
		"alarm_id":          types.Int64Value(data.ID),
		"type":              types.StringValue(data.Type),
		"enabled":           types.BoolValue(data.Enabled),
		"reminder_interval": types.Int64Value(reminder),
		"value_threshold":   types.Int64Value(valueThreshold),
		"value_calculation": types.StringValue(valueCalculation),
		"time_threshold":    types.Int64Value(timeThreshold),
		"vhost_regex":       types.StringValue(vhostRegex),
		"queue_regex":       types.StringValue(queueRegex),
		"message_type":      types.StringValue(messageType),
		"recipients":        recipientsList,
	})
	return obj, objDiags
}
