package cloudamqp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &rabbitMqConfigurationResource{}
	_ resource.ResourceWithConfigure   = &rabbitMqConfigurationResource{}
	_ resource.ResourceWithImportState = &rabbitMqConfigurationResource{}
)

type rabbitMqConfigurationResource struct {
	client *api.API
}

func NewRabbitMqConfigurationResource() resource.Resource {
	return &rabbitMqConfigurationResource{}
}

type rabbitMqConfigurationResourceModel struct {
	ID                       types.String  `tfsdk:"id"`
	InstanceID               types.Int64   `tfsdk:"instance_id"`
	Heartbeat                types.Int64   `tfsdk:"heartbeat"`
	ConnectionMax            types.Int64   `tfsdk:"connection_max"`
	ChannelMax               types.Int64   `tfsdk:"channel_max"`
	ConsumerTimeout          types.Int64   `tfsdk:"consumer_timeout"`
	VmMemoryHighWatermark    types.Float64 `tfsdk:"vm_memory_high_watermark"`
	QueueIndexEmbedMsgsBelow types.Int64   `tfsdk:"queue_index_embed_msgs_below"`
	MaxMessageSize           types.Int64   `tfsdk:"max_message_size"`
	LogExchangeLevel         types.String  `tfsdk:"log_exchange_level"`
	ClusterPartitionHandling types.String  `tfsdk:"cluster_partition_handling"`
	Sleep                    types.Int64   `tfsdk:"sleep"`
	Timeout                  types.Int64   `tfsdk:"timeout"`
}

func (r *rabbitMqConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest,
	resp *resource.MetadataResponse) {

	resp.TypeName = "cloudamqp_rabbitmq_configuration"
}

func (r *rabbitMqConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest,
	resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID (instance_id as string)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"heartbeat": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Set the server AMQP 0-9-1 heartbeat timeout in seconds.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"connection_max": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(-1),
				Description: "Set the maximum permissible number of connections, -1 means infinity.",
				Validators: []validator.Int64{
					int64validator.AtLeast(-1),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"channel_max": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Set the maximum permissible number of channels per connection. 0 means no limit.",
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"consumer_timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "A consumer that has received a message and does not acknowledge that " +
					"message within the timeout in milliseconds.",
				Validators: []validator.Int64{
					int64validator.Any(
						int64validator.OneOf(-1),
						int64validator.Between(10000, 86400000), // 10 seconds to 24 hours
					),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"vm_memory_high_watermark": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				Description: "When the server will enter memory based flow-control as relative to the " +
					"maximum available memory.",
				Validators: []validator.Float64{
					float64validator.Between(0.4, 0.9),
				},
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"queue_index_embed_msgs_below": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Description: "Size in bytes below which to embed messages in the queue index. 0 will " +
					"turn off payload embedding in the queue index.",
				Validators: []validator.Int64{
					int64validator.Between(0, 10485760),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_message_size": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The largest allowed message payload size in bytes.",
				Validators: []validator.Int64{
					int64validator.Between(1, 536870912),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"log_exchange_level": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Log level for the logger used for log integrations and the CloudAMQP " +
					"Console log view.",
				Validators: []validator.String{
					stringvalidator.OneOf("debug", "info", "warning", "error", "critical", "none"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_partition_handling": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set how the cluster should handle network partition.",
				Validators: []validator.String{
					stringvalidator.OneOf("autoheal", "pause_minority", "ignore"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(60),
				Computed:    true,
				Description: "Configurable sleep time in seconds between retries for RabbitMQ configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(3600),
				Computed:    true,
				Description: "Configurable timeout time in seconds for RabbitMQ configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *rabbitMqConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *rabbitMqConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan rabbitMqConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(plan.InstanceID.ValueInt64())
	sleep := int(plan.Sleep.ValueInt64())
	timeout := int(plan.Timeout.ValueInt64())

	data := model.RabbitMqConfigRequest{}
	if !plan.Heartbeat.IsUnknown() {
		data.Heartbeat = utils.Pointer(plan.Heartbeat.ValueInt64())
	}
	if !plan.ConnectionMax.IsUnknown() {
		if plan.ConnectionMax.ValueInt64() == -1 {
			data.ConnectionMax = &model.ConnectionMaxValue{IsInfinity: true}
		} else {
			data.ConnectionMax = &model.ConnectionMaxValue{IsInfinity: false, Value: plan.ConnectionMax.ValueInt64()}
		}
	}
	if !plan.ChannelMax.IsUnknown() {
		data.ChannelMax = utils.Pointer(plan.ChannelMax.ValueInt64())
	}
	if !plan.ConsumerTimeout.IsUnknown() {
		if plan.ConsumerTimeout.ValueInt64() == -1 {
			data.ConsumerTimeout = &model.ConsumerTimeoutValue{IsEnabled: false}
		} else {
			data.ConsumerTimeout = &model.ConsumerTimeoutValue{IsEnabled: true, Value: plan.ConsumerTimeout.ValueInt64()}
		}
	}
	if !plan.VmMemoryHighWatermark.IsUnknown() {
		data.VmMemoryHighWatermark = utils.Pointer(plan.VmMemoryHighWatermark.ValueFloat64())
	}
	if !plan.QueueIndexEmbedMsgsBelow.IsUnknown() {
		data.QueueIndexEmbedMsgsBelow = utils.Pointer(plan.QueueIndexEmbedMsgsBelow.ValueInt64())
	}
	if !plan.MaxMessageSize.IsUnknown() {
		data.MaxMessageSize = utils.Pointer(plan.MaxMessageSize.ValueInt64())
	}
	if !plan.LogExchangeLevel.IsUnknown() {
		data.LogExchangeLevel = utils.Pointer(plan.LogExchangeLevel.ValueString())
	}
	if !plan.ClusterPartitionHandling.IsUnknown() {
		data.ClusterPartitionHandling = utils.Pointer(plan.ClusterPartitionHandling.ValueString())
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating RabbitMQ configuration for instance ID %d with data: %v", instanceID, data))
	err := r.client.UpdateRabbitMqConfiguration(ctx, instanceID, data, sleep, timeout)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to create RabbitMQ configuration: %s", err.Error()))
		return
	}

	dataResp, err := r.client.ReadRabbitMqConfiguration(ctx, instanceID, sleep, timeout)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read RabbitMQ configuration: %s", err.Error()))
		return
	}

	if dataResp == nil {
		resp.Diagnostics.AddError("API Error", "Failed to read RabbitMQ configuration: received nil response")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read RabbitMQ configuration data: %v", dataResp))
	populateRabbitMqConfigModel(&plan, dataResp, instanceID, sleep, timeout)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rabbitMqConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rabbitMqConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(state.InstanceID.ValueInt64())
	sleep := int(state.Sleep.ValueInt64())
	timeout := int(state.Timeout.ValueInt64())

	data, err := r.client.ReadRabbitMqConfiguration(ctx, instanceID, sleep, timeout)
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}

	// Handle resource drift and trigger re-creation if resource been deleted outside the provider
	if data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read RabbitMQ configuration data: %v", data))

	populateRabbitMqConfigModel(&state, data, instanceID, sleep, timeout)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rabbitMqConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rabbitMqConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(plan.InstanceID.ValueInt64())
	sleep := int(plan.Sleep.ValueInt64())
	timeout := int(plan.Timeout.ValueInt64())

	data := model.RabbitMqConfigRequest{}
	if !plan.Heartbeat.IsNull() {
		data.Heartbeat = utils.Pointer(plan.Heartbeat.ValueInt64())
	}
	if !plan.ConnectionMax.IsNull() {
		if plan.ConnectionMax.ValueInt64() == -1 {
			data.ConnectionMax = &model.ConnectionMaxValue{IsInfinity: true}
		} else {
			data.ConnectionMax = &model.ConnectionMaxValue{IsInfinity: false, Value: plan.ConnectionMax.ValueInt64()}
		}
	}
	if !plan.ChannelMax.IsNull() {
		data.ChannelMax = utils.Pointer(plan.ChannelMax.ValueInt64())
	}
	if !plan.ConsumerTimeout.IsNull() {
		if plan.ConsumerTimeout.ValueInt64() == -1 {
			data.ConsumerTimeout = &model.ConsumerTimeoutValue{IsEnabled: false}
		} else {
			data.ConsumerTimeout = &model.ConsumerTimeoutValue{IsEnabled: true, Value: plan.ConsumerTimeout.ValueInt64()}
		}
	}
	if !plan.VmMemoryHighWatermark.IsNull() {
		data.VmMemoryHighWatermark = utils.Pointer(plan.VmMemoryHighWatermark.ValueFloat64())
	}
	if !plan.QueueIndexEmbedMsgsBelow.IsNull() {
		data.QueueIndexEmbedMsgsBelow = utils.Pointer(plan.QueueIndexEmbedMsgsBelow.ValueInt64())
	}
	if !plan.MaxMessageSize.IsNull() {
		data.MaxMessageSize = utils.Pointer(plan.MaxMessageSize.ValueInt64())
	}
	if !plan.LogExchangeLevel.IsNull() {
		data.LogExchangeLevel = utils.Pointer(plan.LogExchangeLevel.ValueString())
	}
	if !plan.ClusterPartitionHandling.IsNull() {
		data.ClusterPartitionHandling = utils.Pointer(plan.ClusterPartitionHandling.ValueString())
	}

	err := r.client.UpdateRabbitMqConfiguration(ctx, instanceID, data, sleep, timeout)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to update RabbitMQ configuration: %s", err.Error()))
		return
	}

	dataResp, err := r.client.ReadRabbitMqConfiguration(ctx, instanceID, sleep, timeout)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read RabbitMQ configuration: %s", err.Error()))
		return
	}

	if dataResp == nil {
		resp.Diagnostics.AddError("API Error", "Failed to read RabbitMQ configuration: received nil response")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Read RabbitMQ configuration data: %v", dataResp))
	populateRabbitMqConfigModel(&plan, dataResp, instanceID, sleep, timeout)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rabbitMqConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op: configuration is not deleted from the server only removed from the state.
}

func (r *rabbitMqConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID // This is a string
	instanceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected numeric instance_id, got: %q", id))
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)
}

// Helper to populate the resource model from the API response
func populateRabbitMqConfigModel(model *rabbitMqConfigurationResourceModel, data *model.RabbitMqConfigResponse, instanceID, sleep, timeout int) {
	model.ID = types.StringValue(strconv.Itoa(instanceID))
	model.InstanceID = types.Int64Value(int64(instanceID))
	model.Sleep = types.Int64Value(int64(sleep))
	model.Timeout = types.Int64Value(int64(timeout))
	model.Heartbeat = types.Int64Value(data.Heartbeat)
	if data.ConnectionMax.IsInfinity {
		model.ConnectionMax = types.Int64Value(-1)
	} else {
		model.ConnectionMax = types.Int64Value(data.ConnectionMax.Value)
	}
	model.ChannelMax = types.Int64Value(data.ChannelMax)
	if !data.ConsumerTimeout.IsEnabled {
		model.ConsumerTimeout = types.Int64Value(-1)
	} else {
		model.ConsumerTimeout = types.Int64Value(data.ConsumerTimeout.Value)
	}
	model.VmMemoryHighWatermark = types.Float64Value(data.VmMemoryHighWatermark)
	if data.QueueIndexEmbedMsgsBelow == nil {
		model.QueueIndexEmbedMsgsBelow = types.Int64Null()
	} else {
		model.QueueIndexEmbedMsgsBelow = types.Int64Value(*data.QueueIndexEmbedMsgsBelow)
	}
	model.MaxMessageSize = types.Int64Value(data.MaxMessageSize)
	model.LogExchangeLevel = types.StringValue(data.LogExchangeLevel)
	model.ClusterPartitionHandling = types.StringValue(data.ClusterPartitionHandling)
}
