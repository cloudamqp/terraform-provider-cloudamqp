package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &customDomainResource{}
	_ resource.ResourceWithConfigure   = &customDomainResource{}
	_ resource.ResourceWithImportState = &customDomainResource{}
)

type customDomainResource struct {
	client *api.API
}

func NewCustomDomainResource() resource.Resource {
	return &customDomainResource{}
}

type customDomainResourceModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	Hostname   types.String `tfsdk:"hostname"`
	Sleep      types.Int64  `tfsdk:"sleep"`
	Timeout    types.Int64  `tfsdk:"timeout"`
}

func (r *customDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_custom_domain"
}

func (r *customDomainResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this resource",
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
			"hostname": schema.StringAttribute{
				Required:    true,
				Description: "The custom hostname.",
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
				Description: "Configurable sleep time in seconds between retries for custom domain configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1800),
				Description: "Configurable timeout time in seconds for custom domain configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *customDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *customDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *customDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(plan.Timeout.ValueInt64())*time.Second)
	defer cancel()

	data, err := r.client.CreateCustomDomain(timeoutCtx, instanceID, plan.Hostname.ValueString(), sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error creating custom domain", err.Error())
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(instanceID, 10))
	plan.Hostname = types.StringValue(data.Hostname)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state customDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID, err := strconv.ParseInt(state.ID.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not convert ID to integer: %s", err))
		return
	}

	// Apply defaults when values are absent (e.g. during import)
	if state.Sleep.IsNull() || state.Sleep.IsUnknown() || state.Sleep.ValueInt64() == 0 {
		state.Sleep = types.Int64Value(10)
	}
	if state.Timeout.IsNull() || state.Timeout.IsUnknown() || state.Timeout.ValueInt64() == 0 {
		state.Timeout = types.Int64Value(1800)
	}
	if state.InstanceID.IsNull() || state.InstanceID.IsUnknown() {
		state.InstanceID = types.Int64Value(instanceID)
	}

	sleep := time.Duration(state.Sleep.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(state.Timeout.ValueInt64())*time.Second)
	defer cancel()

	data, err := r.client.ReadCustomDomain(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error reading custom domain", err.Error())
		return
	}

	// Resource drift: resource not found, trigger re-creation
	if data == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Hostname = types.StringValue(data.Hostname)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *customDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan customDomainResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(plan.Timeout.ValueInt64())*time.Second)
	defer cancel()

	data, err := r.client.UpdateCustomDomain(timeoutCtx, instanceID, plan.Hostname.ValueString(), sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error updating custom domain", err.Error())
		return
	}

	plan.Hostname = types.StringValue(data.Hostname)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customDomainResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep := time.Duration(state.Sleep.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(state.Timeout.ValueInt64())*time.Second)
	defer cancel()

	err := r.client.DeleteCustomDomain(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting custom domain", err.Error())
	}
}
