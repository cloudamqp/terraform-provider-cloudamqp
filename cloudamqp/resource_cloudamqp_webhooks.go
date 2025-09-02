package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &webhookResource{}
	_ resource.ResourceWithConfigure   = &webhookResource{}
	_ resource.ResourceWithImportState = &webhookResource{}
)

type webhookResource struct {
	client *api.API
}

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

type webhookResourceModel struct {
	ID          types.String `tfsdk:"id"`
	InstanceID  types.Int64  `tfsdk:"instance_id"`
	Vhost       types.String `tfsdk:"vhost"`
	Queue       types.String `tfsdk:"queue"`
	WebhookURI  types.String `tfsdk:"webhook_uri"`
	Concurrency types.Int64  `tfsdk:"concurrency"`
	Sleep       types.Int64  `tfsdk:"sleep"`
	Timeout     types.Int64  `tfsdk:"timeout"`
}

func (r *webhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_webhook"
}

func (r *webhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID of the webhook",
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
			"vhost": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual host",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"queue": schema.StringAttribute{
				Required:    true,
				Description: "The queue that should be forwarded, must be a durable queue!",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"webhook_uri": schema.StringAttribute{
				Required:    true,
				Description: "A POST request will be made for each message in the queue to this endpoint",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"concurrency": schema.Int64Attribute{
				Required:    true,
				Description: "How many times the request will be made if previous call fails",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable sleep time in seconds between retries for webhook",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable timeout time in seconds for webhook",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *webhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sleep, timeout := r.extractSleepAndTimeout(&plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	instanceID := plan.InstanceID.ValueInt64()

	request := model.WebhookCreateRequest{
		WebhookURI:  plan.WebhookURI.ValueString(),
		Vhost:       plan.Vhost.ValueString(),
		Queue:       plan.Queue.ValueString(),
		Concurrency: plan.Concurrency.ValueInt64(),
	}

	id, err := r.client.CreateWebhook(timeoutCtx, instanceID, request, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Webhook",
			fmt.Sprintf("Could not create webhook: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(id)
	plan.Sleep = types.Int64Value(int64(sleep.Seconds()))
	plan.Timeout = types.Int64Value(int64(timeout.Seconds()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sleep, timeout := r.extractSleepAndTimeout(&state)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	id := state.ID.ValueString()
	instanceID := state.InstanceID.ValueInt64()

	data, err := r.client.ReadWebhook(timeoutCtx, instanceID, id, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Webhook",
			fmt.Sprintf("Could not read webhook with ID %s: %s", id, err),
		)
		return
	}

	state.ID = types.StringValue(strconv.FormatInt(data.ID, 10))
	state.Concurrency = types.Int64Value(data.Concurrency)
	state.Queue = types.StringValue(data.Queue)
	state.Vhost = types.StringValue(data.Vhost)
	state.WebhookURI = types.StringValue(data.WebhookURI)
	state.Sleep = types.Int64Value(int64(sleep.Seconds()))
	state.Timeout = types.Int64Value(int64(timeout.Seconds()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sleep, timeout := r.extractSleepAndTimeout(&plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	id := plan.ID.ValueString()
	instanceID := plan.InstanceID.ValueInt64()
	webhookID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Convert Webhook ID",
			fmt.Sprintf("Could not convert webhook ID %s to int: %s", id, err),
		)
		return
	}

	request := model.WebhookUpdateRequest{
		WebhookID:   webhookID,
		WebhookURI:  plan.WebhookURI.ValueString(),
		Vhost:       plan.Vhost.ValueString(),
		Queue:       plan.Queue.ValueString(),
		Concurrency: plan.Concurrency.ValueInt64(),
	}

	err = r.client.UpdateWebhook(timeoutCtx, instanceID, id, request, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Webhook",
			fmt.Sprintf("Could not update webhook with ID %s: %s", id, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sleep, timeout := r.extractSleepAndTimeout(&state)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	id := state.ID.ValueString()
	instanceID := state.InstanceID.ValueInt64()

	err := r.client.DeleteWebhook(timeoutCtx, instanceID, id, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Webhook",
			fmt.Sprintf("Could not delete webhook with ID %s: %s", id, err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *webhookResource) extractSleepAndTimeout(plan *webhookResourceModel) (time.Duration, time.Duration) {
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second

	// Default values
	if sleep == 0 {
		sleep = time.Duration(10) * time.Second
	}
	if timeout == 0 {
		timeout = time.Duration(1800) * time.Second
	}

	return sleep, timeout
}
