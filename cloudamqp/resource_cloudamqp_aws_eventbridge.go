package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &awsEventBridgeResource{}
	_ resource.ResourceWithConfigure   = &awsEventBridgeResource{}
	_ resource.ResourceWithImportState = &awsEventBridgeResource{}
)

func NewAwsEventBridgeResource() resource.Resource {
	return &awsEventBridgeResource{}
}

type awsEventBridgeResource struct {
	client *api.API
}

type awsEventBridgeResourceModel struct {
	Id           types.String `tfsdk:"id"`
	InstanceID   types.Int64  `tfsdk:"instance_id"`
	AwsAccountId types.String `tfsdk:"aws_account_id"`
	AwsRegion    types.String `tfsdk:"aws_region"`
	Vhost        types.String `tfsdk:"vhost"`
	QueueName    types.String `tfsdk:"queue"`
	WithHeaders  types.Bool   `tfsdk:"with_headers"`
	Prefetch     types.Int64  `tfsdk:"prefetch"`
	Status       types.String `tfsdk:"status"`
}

func (r *awsEventBridgeResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*api.API)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *awsEventBridgeResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "cloudamqp_integration_aws_eventbridge"
}

func (r *awsEventBridgeResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"aws_account_id": schema.StringAttribute{
				Required:    true,
				Description: "The 12 digit AWS Account ID where you want the events to be sent to.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(12, 12),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"aws_region": schema.StringAttribute{
				Required:    true,
				Description: "The AWS region where you the events to be sent to. (e.g. us-west-1, us-west-2, ..., etc.)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Required:    true,
				Description: "The VHost the queue resides in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"queue": schema.StringAttribute{
				Required:    true,
				Description: "A (durable) queue on your RabbitMQ instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"with_headers": schema.BoolAttribute{
				Required:    true,
				Description: "Include message headers in the event data.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"prefetch": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(1),
				Computed:    true,
				Description: "Number of messages to prefetch. Default set to 1.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Always set to null, unless there is an error starting the EventBridge",
			},
		},
	}
}

func (r *awsEventBridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan awsEventBridgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	request := model.AwsEventBridgeRequest{
		AwsAccountId: plan.AwsAccountId.ValueString(),
		AwsRegion:    plan.AwsRegion.ValueString(),
		Vhost:        plan.Vhost.ValueString(),
		QueueName:    plan.QueueName.ValueString(),
		WithHeaders:  plan.WithHeaders.ValueBool(),
		Prefetch:     plan.Prefetch.ValueInt64Pointer(),
	}

	id, err := r.client.CreateAwsEventBridge(timeoutCtx, plan.InstanceID.ValueInt64(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create AWS EventBridge integration",
			fmt.Sprintf("Could not create AWS EventBridge integration: %s", err),
		)
		return
	}

	plan.Id = types.StringValue(id)
	plan.Status = types.StringNull()

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *awsEventBridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state awsEventBridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	id := state.Id.ValueString()
	instanceID := state.InstanceID.ValueInt64()

	data, err := r.client.ReadAwsEventBridge(timeoutCtx, instanceID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read AWS EventBridge integration",
			fmt.Sprintf("Could not read AWS EventBridge integration with ID %s: %s", id, err),
		)
		return
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("AWS EventBridge integration not found, resource will be recreated: %s", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	state.AwsAccountId = types.StringValue(data.AwsAccountId)
	state.AwsRegion = types.StringValue(data.AwsRegion)
	state.Vhost = types.StringValue(data.Vhost)
	state.QueueName = types.StringValue(data.QueueName)
	state.WithHeaders = types.BoolValue(data.WithHeaders)
	state.Prefetch = types.Int64Value(data.Prefetch)
	state.Status = types.StringPointerValue(data.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *awsEventBridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This resource does not implement the Update function
}

func (r *awsEventBridgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data awsEventBridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	err := r.client.DeleteAwsEventBridge(timeoutCtx, data.InstanceID.ValueInt64(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete AWS EventBridge integration",
			fmt.Sprintf("Could not delete AWS EventBridge integration with ID %s: %s", data.Id.ValueString(), err),
		)
	}
}

func (r *awsEventBridgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("ImportState: ID=%s", req.ID))
	if !strings.Contains(req.ID, ",") {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {resource_id},{instance_id}")
		return
	}

	idSplit := strings.Split(req.ID, ",")
	if len(idSplit) != 2 {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {resource_id},{instance_id}")
		return
	}
	instanceID, err := strconv.Atoi(idSplit[1])
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance_id in import ID", fmt.Sprintf("Could not convert instance_id to int: %s", err))
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), idSplit[0])
	resp.State.SetAttribute(ctx, path.Root("instance_id"), int64(instanceID))
}
