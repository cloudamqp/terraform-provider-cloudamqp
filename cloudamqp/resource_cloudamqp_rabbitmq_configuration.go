package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	ID                                    types.String  `tfsdk:"id"`
	InstanceID                            types.Int64   `tfsdk:"instance_id"`
	Heartbeat                             types.Int64   `tfsdk:"heartbeat"`
	ConnectionMax                         types.Int64   `tfsdk:"connection_max"`
	ChannelMax                            types.Int64   `tfsdk:"channel_max"`
	ConsumerTimeout                       types.Int64   `tfsdk:"consumer_timeout"`
	VmMemoryHighWatermark                 types.Float64 `tfsdk:"vm_memory_high_watermark"`
	QueueIndexEmbedMsgsBelow              types.Int64   `tfsdk:"queue_index_embed_msgs_below"`
	MaxMessageSize                        types.Int64   `tfsdk:"max_message_size"`
	LogExchangeLevel                      types.String  `tfsdk:"log_exchange_level"`
	ClusterPartitionHandling              types.String  `tfsdk:"cluster_partition_handling"`
	MessageInterceptorsTimestampOverwrite types.String  `tfsdk:"message_interceptors_timestamp_overwrite"`
	// MQTT settings
	MQTTVhost        types.String `tfsdk:"mqtt_vhost"`
	MQTTExchange     types.String `tfsdk:"mqtt_exchange"`
	MQTTSSLCertLogin types.Bool   `tfsdk:"mqtt_ssl_cert_login"`
	// SSL settings
	SSLCertLoginFrom           types.String `tfsdk:"ssl_cert_login_from"`
	SSLOptionsFailIfNoPeerCert types.Bool   `tfsdk:"ssl_options_fail_if_no_peer_cert"`
	SSLOptionsVerify           types.String `tfsdk:"ssl_options_verify"`
	// Sleep/timeout for retries
	Sleep   types.Int64 `tfsdk:"sleep"`
	Timeout types.Int64 `tfsdk:"timeout"`
}

func (r *rabbitMqConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
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
			"message_interceptors_timestamp_overwrite": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Sets a timestamp header on incoming messages. enabled_with_overwrite will " +
					"overwrite any existing timestamps in the header.",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("enabled_with_overwrite", "enabled", "disabled"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mqtt_vhost": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Virtual host for MQTT connections.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mqtt_exchange": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The exchange option determines which exchange messages from MQTT clients are published to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mqtt_ssl_cert_login": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable SSL certificate-based authentication for MQTT connections.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cert_login_from": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Determines which certificate field to use as the username for TLS-based authentication.",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("common_name", "distinguished_name"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_options_fail_if_no_peer_cert": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "When set to true, TLS connections will fail if the client does not provide a certificate.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_options_verify": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Controls peer certificate verification for TLS connections.",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("verify_none", "verify_peer"),
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
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(3600),
				Computed:    true,
				Description: "Configurable timeout time in seconds for RabbitMQ configuration",
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

	sleep := plan.Sleep.ValueInt64()
	timeout := plan.Timeout.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	instanceID := plan.InstanceID.ValueInt64()
	request := r.populateCreateRequest(plan)

	err := r.client.UpdateRabbitMqConfiguration(timeoutCtx, instanceID, request, sleep)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to create RabbitMQ configuration: %s", err.Error()))
		return
	}

	dataResp, err := r.client.ReadRabbitMqConfiguration(ctx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read RabbitMQ configuration: %s", err.Error()))
		return
	}

	if dataResp == nil {
		resp.Diagnostics.AddError("API Error", "Failed to read RabbitMQ configuration: received nil response")
		return
	}

	r.populateResourceModel(&plan, dataResp, instanceID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rabbitMqConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state rabbitMqConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Sleep/timeout with default values
	sleep := state.Sleep.ValueInt64()
	timeout := state.Timeout.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	instanceID := state.InstanceID.ValueInt64()

	data, err := r.client.ReadRabbitMqConfiguration(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}
	if data == nil {
		tflog.Warn(ctx, fmt.Sprintf("RabbitMQ configuration resource drift for instance ID %d, trigger re-create", instanceID))
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateResourceModel(&state, data, instanceID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rabbitMqConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rabbitMqConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sleep := plan.Sleep.ValueInt64()
	timeout := plan.Timeout.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	instanceID := plan.InstanceID.ValueInt64()
	request, changed := r.populateUpdateRequest(plan, state)

	if !changed {
		// No rabbitmq configuration changes detected, only save the state
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	err := r.client.UpdateRabbitMqConfiguration(timeoutCtx, instanceID, request, sleep)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to update RabbitMQ configuration: %s", err.Error()))
		return
	}

	dataResp, err := r.client.ReadRabbitMqConfiguration(ctx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to read RabbitMQ configuration: %s", err.Error()))
		return
	}

	if dataResp == nil {
		resp.Diagnostics.AddError("API Error", "Failed to read RabbitMQ configuration: received nil response")
		return
	}

	r.populateResourceModel(&plan, dataResp, instanceID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *rabbitMqConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op: configuration is not deleted from the server only removed from the state.
	resp.State.RemoveResource(ctx)
}

func (r *rabbitMqConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	instanceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected numeric instance_id, got: %q", id))
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), id)
	resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)
	resp.State.SetAttribute(ctx, path.Root("sleep"), 60)     // default value
	resp.State.SetAttribute(ctx, path.Root("timeout"), 3600) // default value
}

// Convert API response to resource model
func (r *rabbitMqConfigurationResource) populateResourceModel(resourceModel *rabbitMqConfigurationResourceModel, data *model.RabbitMqConfigResponse, instanceID int64) {

	resourceModel.ID = types.StringValue(strconv.Itoa(int(instanceID)))
	resourceModel.InstanceID = types.Int64Value(instanceID)
	resourceModel.Heartbeat = types.Int64Value(data.Heartbeat)
	resourceModel.ChannelMax = types.Int64Value(data.ChannelMax)
	resourceModel.MaxMessageSize = types.Int64Value(data.MaxMessageSize)
	resourceModel.LogExchangeLevel = types.StringValue(data.LogExchangeLevel)
	resourceModel.ClusterPartitionHandling = types.StringValue(data.ClusterPartitionHandling)
	resourceModel.VmMemoryHighWatermark = types.Float64Value(data.VmMemoryHighWatermark)
	resourceModel.MessageInterceptorsTimestampOverwrite = types.StringValue(data.MessageInterceptorsTimestampOverwrite)
	// MQTT settings
	resourceModel.MQTTVhost = types.StringValue(data.MQTTVhost)
	resourceModel.MQTTExchange = types.StringValue(data.MQTTExchange)
	resourceModel.MQTTSSLCertLogin = types.BoolValue(bool(data.MQTTSSLCertLogin))
	// SSL settings
	resourceModel.SSLCertLoginFrom = types.StringValue(data.SSLCertLoginFrom)
	resourceModel.SSLOptionsFailIfNoPeerCert = types.BoolValue(bool(data.SSLOptionsFailIfNoPeerCert))
	resourceModel.SSLOptionsVerify = types.StringValue(data.SSLOptionsVerify)

	// Handle special cases for pointer and custom types
	if data.ConnectionMax == nil {
		resourceModel.ConnectionMax = types.Int64Value(-1) // Default value when not set
	} else if data.ConnectionMax.IsInfinity {
		resourceModel.ConnectionMax = types.Int64Value(-1)
	} else {
		resourceModel.ConnectionMax = types.Int64Value(data.ConnectionMax.Value)
	}

	if !data.ConsumerTimeout.IsEnabled {
		resourceModel.ConsumerTimeout = types.Int64Value(-1)
	} else {
		resourceModel.ConsumerTimeout = types.Int64Value(data.ConsumerTimeout.Value)
	}

	if data.QueueIndexEmbedMsgsBelow == nil {
		resourceModel.QueueIndexEmbedMsgsBelow = types.Int64Value(4096) // Default value when not set
	} else {
		resourceModel.QueueIndexEmbedMsgsBelow = types.Int64Value(*data.QueueIndexEmbedMsgsBelow)
	}
}

// Populate API create request from resource model
func (r *rabbitMqConfigurationResource) populateCreateRequest(plan rabbitMqConfigurationResourceModel) model.RabbitMqConfigRequest {
	request := model.RabbitMqConfigRequest{}

	if !plan.Heartbeat.IsUnknown() {
		request.Heartbeat = plan.Heartbeat.ValueInt64Pointer()
	}

	if !plan.ConnectionMax.IsUnknown() {
		if plan.ConnectionMax.ValueInt64() == -1 {
			request.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: true})
		} else {
			request.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: false, Value: plan.ConnectionMax.ValueInt64()})
		}
	}

	if !plan.ChannelMax.IsUnknown() {
		request.ChannelMax = plan.ChannelMax.ValueInt64Pointer()
	}

	if !plan.ConsumerTimeout.IsUnknown() {
		if plan.ConsumerTimeout.ValueInt64() == -1 {
			request.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: false})
		} else {
			request.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: true, Value: plan.ConsumerTimeout.ValueInt64()})
		}
	}

	if !plan.VmMemoryHighWatermark.IsUnknown() {
		request.VmMemoryHighWatermark = plan.VmMemoryHighWatermark.ValueFloat64Pointer()
	}

	if !plan.QueueIndexEmbedMsgsBelow.IsUnknown() {
		request.QueueIndexEmbedMsgsBelow = plan.QueueIndexEmbedMsgsBelow.ValueInt64Pointer()
	}

	if !plan.MaxMessageSize.IsUnknown() {
		request.MaxMessageSize = plan.MaxMessageSize.ValueInt64Pointer()
	}

	if !plan.LogExchangeLevel.IsUnknown() {
		request.LogExchangeLevel = plan.LogExchangeLevel.ValueString()
	}

	if !plan.ClusterPartitionHandling.IsUnknown() {
		request.ClusterPartitionHandling = plan.ClusterPartitionHandling.ValueString()
	}

	if !plan.MessageInterceptorsTimestampOverwrite.IsUnknown() {
		request.MessageInterceptorsTimestampOverwrite = plan.MessageInterceptorsTimestampOverwrite.ValueStringPointer()
	}

	// MQTT settings
	if !plan.MQTTVhost.IsUnknown() {
		request.MQTTVhost = plan.MQTTVhost.ValueStringPointer()
	}

	if !plan.MQTTExchange.IsUnknown() {
		request.MQTTExchange = plan.MQTTExchange.ValueStringPointer()
	}

	if !plan.MQTTSSLCertLogin.IsUnknown() {
		request.MQTTSSLCertLogin = plan.MQTTSSLCertLogin.ValueBoolPointer()
	}

	// SSL settings
	if !plan.SSLCertLoginFrom.IsUnknown() {
		request.SSLCertLoginFrom = plan.SSLCertLoginFrom.ValueStringPointer()
	}

	if !plan.SSLOptionsFailIfNoPeerCert.IsUnknown() {
		request.SSLOptionsFailIfNoPeerCert = plan.SSLOptionsFailIfNoPeerCert.ValueBoolPointer()
	}

	if !plan.SSLOptionsVerify.IsUnknown() {
		request.SSLOptionsVerify = plan.SSLOptionsVerify.ValueStringPointer()
	}

	return request
}

// Populate API update request from resource model
func (r *rabbitMqConfigurationResource) populateUpdateRequest(plan, state rabbitMqConfigurationResourceModel) (model.RabbitMqConfigRequest, bool) {
	request := model.RabbitMqConfigRequest{}
	changed := false

	if !plan.Heartbeat.IsNull() && !plan.Heartbeat.Equal(state.Heartbeat) {
		request.Heartbeat = plan.Heartbeat.ValueInt64Pointer()
		changed = true
	}

	if !plan.ConnectionMax.IsNull() && !plan.ConnectionMax.Equal(state.ConnectionMax) {
		if plan.ConnectionMax.ValueInt64() == -1 {
			request.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: true})
		} else {
			request.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: false, Value: plan.ConnectionMax.ValueInt64()})
		}
		changed = true
	}

	if !plan.ChannelMax.IsNull() && !plan.ChannelMax.Equal(state.ChannelMax) {
		request.ChannelMax = plan.ChannelMax.ValueInt64Pointer()
		changed = true
	}

	if !plan.ConsumerTimeout.IsNull() && !plan.ConsumerTimeout.Equal(state.ConsumerTimeout) {
		if plan.ConsumerTimeout.ValueInt64() == -1 {
			request.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: false})
		} else {
			request.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: true, Value: plan.ConsumerTimeout.ValueInt64()})
		}
		changed = true
	}

	if !plan.VmMemoryHighWatermark.IsNull() && !plan.VmMemoryHighWatermark.Equal(state.VmMemoryHighWatermark) {
		request.VmMemoryHighWatermark = plan.VmMemoryHighWatermark.ValueFloat64Pointer()
		changed = true
	}

	if !plan.QueueIndexEmbedMsgsBelow.IsNull() && !plan.QueueIndexEmbedMsgsBelow.Equal(state.QueueIndexEmbedMsgsBelow) {
		request.QueueIndexEmbedMsgsBelow = plan.QueueIndexEmbedMsgsBelow.ValueInt64Pointer()
		changed = true
	}

	if !plan.MaxMessageSize.IsNull() && !plan.MaxMessageSize.Equal(state.MaxMessageSize) {
		request.MaxMessageSize = plan.MaxMessageSize.ValueInt64Pointer()
		changed = true
	}

	if !plan.LogExchangeLevel.IsNull() && !plan.LogExchangeLevel.Equal(state.LogExchangeLevel) {
		request.LogExchangeLevel = plan.LogExchangeLevel.ValueString()
		changed = true
	}

	if !plan.ClusterPartitionHandling.IsNull() && !plan.ClusterPartitionHandling.Equal(state.ClusterPartitionHandling) {
		request.ClusterPartitionHandling = plan.ClusterPartitionHandling.ValueString()
		changed = true
	}

	if !plan.MessageInterceptorsTimestampOverwrite.IsNull() && !plan.MessageInterceptorsTimestampOverwrite.Equal(state.MessageInterceptorsTimestampOverwrite) {
		request.MessageInterceptorsTimestampOverwrite = plan.MessageInterceptorsTimestampOverwrite.ValueStringPointer()
		changed = true
	}

	// MQTT settings
	if !plan.MQTTVhost.IsNull() && !plan.MQTTVhost.Equal(state.MQTTVhost) {
		request.MQTTVhost = plan.MQTTVhost.ValueStringPointer()
		changed = true
	}

	if !plan.MQTTExchange.IsNull() && !plan.MQTTExchange.Equal(state.MQTTExchange) {
		request.MQTTExchange = plan.MQTTExchange.ValueStringPointer()
		changed = true
	}

	if !plan.MQTTSSLCertLogin.IsNull() && !plan.MQTTSSLCertLogin.Equal(state.MQTTSSLCertLogin) {
		request.MQTTSSLCertLogin = plan.MQTTSSLCertLogin.ValueBoolPointer()
		changed = true
	}

	// SSL settings
	if !plan.SSLCertLoginFrom.IsNull() && !plan.SSLCertLoginFrom.Equal(state.SSLCertLoginFrom) {
		request.SSLCertLoginFrom = plan.SSLCertLoginFrom.ValueStringPointer()
		changed = true
	}

	if !plan.SSLOptionsFailIfNoPeerCert.IsNull() && !plan.SSLOptionsFailIfNoPeerCert.Equal(state.SSLOptionsFailIfNoPeerCert) {
		request.SSLOptionsFailIfNoPeerCert = plan.SSLOptionsFailIfNoPeerCert.ValueBoolPointer()
		changed = true
	}

	if !plan.SSLOptionsVerify.IsNull() && !plan.SSLOptionsVerify.Equal(state.SSLOptionsVerify) {
		request.SSLOptionsVerify = plan.SSLOptionsVerify.ValueStringPointer()
		changed = true
	}

	return request, changed
}
