package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	ID         types.String          `tfsdk:"id"`
	InstanceID types.Int64           `tfsdk:"instance_id"`
	Type       types.String          `tfsdk:"type"`
	Alarms     []alarmDataSourceItem `tfsdk:"alarms"`
}

type alarmDataSourceItem struct {
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

func (d *alarmsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_alarms"
}

func (d *alarmsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve a list of pre-defined or created alarms.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this data source",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
			},
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Type of the alarm to filter by",
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
		},
		Blocks: map[string]schema.Block{
			"alarms": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"alarm_id": schema.Int64Attribute{
							Computed:    true,
							Description: "Alarm identifier",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the alarm",
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
				},
			},
		},
	}
}

func (d *alarmsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	var config alarmsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := config.InstanceID.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	alarms, err := d.client.ListAlarms(timeoutCtx, instanceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to List Alarms",
			fmt.Sprintf("Could not list alarms: %s", err),
		)
		return
	}

	if !config.Type.IsNull() {
		for _, alarm := range alarms {
			if alarm.Type == config.Type.ValueString() {
				config.Alarms = append(config.Alarms, d.populateResourceModel(ctx, alarm))
			}
		}
		config.ID = types.StringValue(fmt.Sprintf("%d.%s.alarms", instanceID, config.Type.ValueString()))
	} else {
		for _, alarm := range alarms {
			config.Alarms = append(config.Alarms, d.populateResourceModel(ctx, alarm))
		}
		config.ID = types.StringValue(fmt.Sprintf("%d.alarms", instanceID))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (d *alarmsDataSource) populateResourceModel(ctx context.Context, data model.AlarmResponse) alarmDataSourceItem {
	alarm := alarmDataSourceItem{}
	alarm.AlarmID = types.Int64Value(int64(data.ID))
	alarm.Type = types.StringValue(data.Type)
	alarm.Enabled = types.BoolValue(data.Enabled)

	if data.ReminderInterval != nil {
		alarm.ReminderInterval = types.Int64Value(*data.ReminderInterval)
	} else {
		alarm.ReminderInterval = types.Int64Null()
	}

	if data.ValueThreshold != nil {
		alarm.ValueThreshold = types.Int64Value(*data.ValueThreshold)
	} else {
		alarm.ValueThreshold = types.Int64Null()
	}

	if data.ValueCalculation != nil {
		alarm.ValueCalculation = types.StringValue(*data.ValueCalculation)
	} else {
		alarm.ValueCalculation = types.StringNull()
	}

	if data.TimeThreshold != nil {
		alarm.TimeThreshold = types.Int64Value(*data.TimeThreshold)
	} else {
		alarm.TimeThreshold = types.Int64Null()
	}

	if data.VhostRegex != nil {
		alarm.VhostRegex = types.StringValue(*data.VhostRegex)
	} else {
		alarm.VhostRegex = types.StringNull()
	}

	if data.QueueRegex != nil {
		alarm.QueueRegex = types.StringValue(*data.QueueRegex)
	} else {
		alarm.QueueRegex = types.StringNull()
	}

	if data.MessageType != nil {
		alarm.MessageType = types.StringValue(*data.MessageType)
	} else {
		alarm.MessageType = types.StringNull()
	}

	if data.Recipients != nil {
		recipientsList, _ := types.ListValueFrom(ctx, types.Int64Type, *data.Recipients)
		alarm.Recipients = recipientsList
	} else {
		alarm.Recipients = types.ListNull(types.Int64Type)
	}

	return alarm
}
