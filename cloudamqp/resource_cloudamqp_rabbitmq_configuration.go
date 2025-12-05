package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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
	// MQTT settings
	MQTTVhost        types.String `tfsdk:"mqtt_vhost"`
	MQTTExchange     types.String `tfsdk:"mqtt_exchange"`
	MQTTSSLCertLogin types.Bool   `tfsdk:"mqtt_ssl_cert_login"`
	// SSL settings
	SSLCertLoginFrom           types.String `tfsdk:"ssl_cert_login_from"`
	SSLOptionsFailIfNoPeerCert types.Bool   `tfsdk:"ssl_options_fail_if_no_peer_cert"`
	SSLOptionsVerify           types.String `tfsdk:"ssl_options_verify"`
	// Message interceptor settings
	MessageInterceptorsTimestampOverwrite types.String `tfsdk:"message_interceptors_timestamp_overwrite"`
	// Sleep/timeout for retries
	Sleep   types.Int64 `tfsdk:"sleep"`
	Timeout types.Int64 `tfsdk:"timeout"`
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
			"message_interceptors_timestamp_overwrite": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Sets a timestamp header on incoming messages. enabled_with_overwrite will " +
					"overwrite any existing timestamps in the header.",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive("enabled", "enabled_with_overwrite", "disabled"),
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
	data := r.populateCreateRequest(ctx, &plan)

	err := r.client.UpdateRabbitMqConfiguration(timeoutCtx, instanceID, data, sleep)
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

	r.populateRabbitMqConfigModel(&plan, dataResp, instanceID, sleep, timeout)
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
	if sleep == 0 {
		sleep = 60
	}
	timeout := state.Timeout.ValueInt64()
	if timeout == 0 {
		timeout = 3600
	}
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

	r.populateRabbitMqConfigModel(&state, data, instanceID, sleep, timeout)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rabbitMqConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan rabbitMqConfigurationResourceModel
	var state rabbitMqConfigurationResourceModel
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
	data, changed := r.populateUpdateRequest(ctx, plan, state)

	if !changed {
		// No changes detected, skip update
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	err := r.client.UpdateRabbitMqConfiguration(timeoutCtx, instanceID, data, sleep)
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

	r.populateRabbitMqConfigModel(&plan, dataResp, instanceID, sleep, timeout)
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
}

// Handle data conversion from API response to resource model
func (r *rabbitMqConfigurationResource) populateRabbitMqConfigModel(resourceModel *rabbitMqConfigurationResourceModel,
	data *model.RabbitMqConfigResponse, instanceID, sleep, timeout int64) {

	resourceModel.ID = types.StringValue(strconv.Itoa(int(instanceID)))
	resourceModel.InstanceID = types.Int64Value(instanceID)
	resourceModel.Sleep = types.Int64Value(sleep)
	resourceModel.Timeout = types.Int64Value(timeout)
	resourceModel.Heartbeat = types.Int64Value(data.Heartbeat)
	resourceModel.ChannelMax = types.Int64Value(data.ChannelMax)
	resourceModel.MaxMessageSize = types.Int64Value(data.MaxMessageSize)
	resourceModel.LogExchangeLevel = types.StringValue(data.LogExchangeLevel)
	resourceModel.ClusterPartitionHandling = types.StringValue(data.ClusterPartitionHandling)
	resourceModel.VmMemoryHighWatermark = types.Float64Value(data.VmMemoryHighWatermark)
	// MQTT settings
	resourceModel.MQTTVhost = types.StringValue(data.MQTTVhost)
	resourceModel.MQTTExchange = types.StringValue(data.MQTTExchange)
	resourceModel.MQTTSSLCertLogin = types.BoolValue(bool(data.MQTTSSLCertLogin))
	// SSL settings
	resourceModel.SSLCertLoginFrom = types.StringValue(data.SSLCertLoginFrom)
	resourceModel.SSLOptionsFailIfNoPeerCert = types.BoolValue(bool(data.SSLOptionsFailIfNoPeerCert))
	if data.SSLOptionsVerify == nil {
		resourceModel.SSLOptionsVerify = types.StringValue("verify_none")
	} else {
		resourceModel.SSLOptionsVerify = types.StringValue(*data.SSLOptionsVerify)
	}
	// Message interceptor settings
	if data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite == nil {
		resourceModel.MessageInterceptorsTimestampOverwrite = types.StringValue("disabled")
	} else if *data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite == "true" {
		resourceModel.MessageInterceptorsTimestampOverwrite = types.StringValue("enabled_with_overwrite")
	} else if *data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite == "false" {
		resourceModel.MessageInterceptorsTimestampOverwrite = types.StringValue("enabled")
	}

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

// Handle data conversion from resource model to create API request
func (r *rabbitMqConfigurationResource) populateCreateRequest(ctx context.Context, plan *rabbitMqConfigurationResourceModel) model.RabbitMqConfigRequest {
	data := model.RabbitMqConfigRequest{}

	tflog.Info(ctx, fmt.Sprintf("Populating RabbitMQ configuration create request, plan: %+v", plan))

	// TODO: Can plan.Heartbeat.ValueInt64Pointer() be used instead? Same goes for the rest of the fields.
	if !plan.Heartbeat.IsUnknown() {
		data.Heartbeat = utils.Pointer(plan.Heartbeat.ValueInt64())
	}

	if !plan.ConnectionMax.IsUnknown() {
		if plan.ConnectionMax.ValueInt64() == -1 {
			data.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: true})
		} else {
			data.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: false, Value: plan.ConnectionMax.ValueInt64()})
		}
	}

	if !plan.ChannelMax.IsUnknown() {
		data.ChannelMax = utils.Pointer(plan.ChannelMax.ValueInt64())
	}

	if !plan.ConsumerTimeout.IsUnknown() {
		if plan.ConsumerTimeout.ValueInt64() == -1 {
			data.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: false})
		} else {
			data.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: true, Value: plan.ConsumerTimeout.ValueInt64()})
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
		data.LogExchangeLevel = plan.LogExchangeLevel.ValueString()
	}

	if !plan.ClusterPartitionHandling.IsUnknown() {
		data.ClusterPartitionHandling = plan.ClusterPartitionHandling.ValueString()
	}

	if !plan.MQTTVhost.IsUnknown() {
		data.MQTTVhost = utils.Pointer(plan.MQTTVhost.ValueString())
	}

	if !plan.MQTTExchange.IsUnknown() {
		data.MQTTExchange = utils.Pointer(plan.MQTTExchange.ValueString())
	}

	if !plan.MQTTSSLCertLogin.IsUnknown() {
		data.MQTTSSLCertLogin = utils.Pointer(bool(plan.MQTTSSLCertLogin.ValueBool()))
	}

	if !plan.SSLCertLoginFrom.IsUnknown() {
		data.SSLCertLoginFrom = utils.Pointer(plan.SSLCertLoginFrom.ValueString())
	}

	if !plan.SSLOptionsFailIfNoPeerCert.IsUnknown() {
		data.SSLOptionsFailIfNoPeerCert = utils.Pointer(bool(plan.SSLOptionsFailIfNoPeerCert.ValueBool()))
	}

	if !plan.SSLOptionsVerify.IsUnknown() {
		data.SSLOptionsVerify = utils.Pointer(plan.SSLOptionsVerify.ValueString())
	}

	if !plan.MessageInterceptorsTimestampOverwrite.IsUnknown() {
		value := plan.MessageInterceptorsTimestampOverwrite.ValueString()
		switch strings.ToLower(value) {
		case "disabled":
			data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite = nil
		case "enabled_with_overwrite":
			data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite = utils.Pointer("true")
		case "enabled":
			data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite = utils.Pointer("false")
		}
	}

	return data
}

// Handle data conversion from resource model to update API request
func (r *rabbitMqConfigurationResource) populateUpdateRequest(ctx context.Context, plan, state rabbitMqConfigurationResourceModel) (model.RabbitMqConfigRequest, bool) {
	data := model.RabbitMqConfigRequest{}
	changed := false
	tflog.Info(ctx, fmt.Sprintf("Populating RabbitMQ configuration update request, plan: %+v", plan))

	if !plan.Heartbeat.IsNull() {
		if plan.Heartbeat.ValueInt64() != state.Heartbeat.ValueInt64() {
			data.Heartbeat = utils.Pointer(plan.Heartbeat.ValueInt64())
		}
	}

	if !plan.ConnectionMax.IsNull() {
		if plan.ConnectionMax.ValueInt64() != state.ConnectionMax.ValueInt64() {
			if plan.ConnectionMax.ValueInt64() == -1 {
				data.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: true})
			} else {
				data.ConnectionMax = utils.Pointer(model.ConnectionMaxValue{IsInfinity: false, Value: plan.ConnectionMax.ValueInt64()})
			}
			changed = true
		}
	}

	if !plan.ChannelMax.IsNull() {
		if plan.ChannelMax.ValueInt64() != state.ChannelMax.ValueInt64() {
			data.ChannelMax = utils.Pointer(plan.ChannelMax.ValueInt64())
			changed = true
		}
	}

	if !plan.ConsumerTimeout.IsNull() {
		if plan.ConsumerTimeout.ValueInt64() != state.ConsumerTimeout.ValueInt64() {
			if plan.ConsumerTimeout.ValueInt64() == -1 {
				data.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: false})
			} else {
				data.ConsumerTimeout = utils.Pointer(model.ConsumerTimeoutValue{IsEnabled: true, Value: plan.ConsumerTimeout.ValueInt64()})
			}
			changed = true
		}
	}

	if !plan.VmMemoryHighWatermark.IsNull() {
		if plan.VmMemoryHighWatermark.ValueFloat64() != state.VmMemoryHighWatermark.ValueFloat64() {
			data.VmMemoryHighWatermark = utils.Pointer(plan.VmMemoryHighWatermark.ValueFloat64())
			changed = true
		}
	}

	if !plan.QueueIndexEmbedMsgsBelow.IsNull() {
		if plan.QueueIndexEmbedMsgsBelow.ValueInt64() != state.QueueIndexEmbedMsgsBelow.ValueInt64() {
			data.QueueIndexEmbedMsgsBelow = utils.Pointer(plan.QueueIndexEmbedMsgsBelow.ValueInt64())
			changed = true
		}
	}

	if !plan.MaxMessageSize.IsNull() {
		if plan.MaxMessageSize.ValueInt64() != state.MaxMessageSize.ValueInt64() {
			data.MaxMessageSize = utils.Pointer(plan.MaxMessageSize.ValueInt64())
			changed = true
		}
	}

	if !plan.LogExchangeLevel.IsNull() {
		if plan.LogExchangeLevel.ValueString() != state.LogExchangeLevel.ValueString() {
			data.LogExchangeLevel = plan.LogExchangeLevel.ValueString()
			changed = true
		}
	}

	if !plan.ClusterPartitionHandling.IsNull() {
		if plan.ClusterPartitionHandling.ValueString() != state.ClusterPartitionHandling.ValueString() {
			data.ClusterPartitionHandling = plan.ClusterPartitionHandling.ValueString()
			changed = true
		}
	}

	if !plan.MQTTVhost.IsNull() {
		if plan.MQTTVhost.ValueString() != state.MQTTVhost.ValueString() {
			data.MQTTVhost = utils.Pointer(plan.MQTTVhost.ValueString())
			changed = true
		}
	}

	if !plan.MQTTExchange.IsNull() {
		if plan.MQTTExchange.ValueString() != state.MQTTExchange.ValueString() {
			data.MQTTExchange = utils.Pointer(plan.MQTTExchange.ValueString())
			changed = true
		}
	}

	if !plan.MQTTSSLCertLogin.IsNull() {
		if plan.MQTTSSLCertLogin.ValueBool() != state.MQTTSSLCertLogin.ValueBool() {
			data.MQTTSSLCertLogin = utils.Pointer(bool(plan.MQTTSSLCertLogin.ValueBool()))
			changed = true
		}
	}

	if !plan.SSLCertLoginFrom.IsNull() {
		if plan.SSLCertLoginFrom.ValueString() != state.SSLCertLoginFrom.ValueString() {
			data.SSLCertLoginFrom = utils.Pointer(plan.SSLCertLoginFrom.ValueString())
			changed = true
		}
	}

	if !plan.SSLOptionsFailIfNoPeerCert.IsNull() {
		if plan.SSLOptionsFailIfNoPeerCert.ValueBool() != state.SSLOptionsFailIfNoPeerCert.ValueBool() {
			data.SSLOptionsFailIfNoPeerCert = utils.Pointer(bool(plan.SSLOptionsFailIfNoPeerCert.ValueBool()))
			changed = true
		}
	}

	if !plan.SSLOptionsVerify.IsNull() {
		if plan.SSLOptionsVerify.ValueString() != state.SSLOptionsVerify.ValueString() {
			data.SSLOptionsVerify = utils.Pointer(plan.SSLOptionsVerify.ValueString())
			changed = true
		}
	}

	if !plan.MessageInterceptorsTimestampOverwrite.IsNull() {
		if plan.MessageInterceptorsTimestampOverwrite.ValueString() != state.MessageInterceptorsTimestampOverwrite.ValueString() {
			value := plan.MessageInterceptorsTimestampOverwrite.ValueString()
			switch strings.ToLower(value) {
			case "enabled_with_overwrite":
				data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite = utils.Pointer("true")
			case "enabled":
				data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite = utils.Pointer("false")
			case "disabled":
				data.MessageInterceptorsIncomingSetHeaderTimestampOverwrite = nil
			}
			changed = true
		}
	}

	return data, changed
}
