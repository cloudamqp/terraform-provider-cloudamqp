package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &trustStoreConfigurationResource{}
	_ resource.ResourceWithConfigure   = &trustStoreConfigurationResource{}
	_ resource.ResourceWithImportState = &trustStoreConfigurationResource{}
)

type trustStoreConfigurationResource struct {
	client *api.API
}

func NewTrustStoreConfigurationResource() resource.Resource {
	return &trustStoreConfigurationResource{}
}

type trustStoreConfigurationResourceModel struct {
	ID              types.String `tfsdk:"id"`
	InstanceID      types.Int64  `tfsdk:"instance_id"`
	Url             types.String `tfsdk:"url"`
	RefreshInterval types.Int64  `tfsdk:"refresh_interval"`
	Provider        types.String `tfsdk:"provider"`
	Sleep           types.Int64  `tfsdk:"sleep"`
	Timeout         types.Int64  `tfsdk:"timeout"`
}

func (r *trustStoreConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_trust_store_configuration"
}

func (r *trustStoreConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Description: "Resource ID",
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
			"url": schema.StringAttribute{
				Required:    true,
				Description: "URL to fetch trust store certificates from",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"refresh_interval": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(30),
				Computed:    true,
				Description: "Interval in seconds to refresh the trust store",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"provider": schema.StringAttribute{
				Optional:    true,
				Default:     stringdefault.StaticString("http"),
				Computed:    true,
				Description: "Trust store provider (currently only 'http' is supported)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(30),
				Computed:    true,
				Description: "Configurable sleep time in seconds between retries for trust store configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(1800),
				Computed:    true,
				Description: "Configurable timeout time in seconds for trust store configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *trustStoreConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *trustStoreConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("ImportState: ID=%s", req.ID))
	instanceID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected numeric instance_id, got: %q", req.ID))
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), req.ID)
	resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)
	// Set default values for optional/computed attributes
	resp.State.SetAttribute(ctx, path.Root("refresh_interval"), 30)
	resp.State.SetAttribute(ctx, path.Root("provider"), "http")
	resp.State.SetAttribute(ctx, path.Root("sleep"), 30)
	resp.State.SetAttribute(ctx, path.Root("timeout"), 1800)
}

func (r *trustStoreConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan trustStoreConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	request := model.TrustStoreRequest{
		Url:             plan.Url.ValueString(),
		RefreshInterval: plan.RefreshInterval.ValueInt64(),
		Provider:        plan.Provider.ValueString(),
	}

	job, err := r.client.CreateTrustStoreConfiguration(timeoutCtx, instanceID, sleep, request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating trust store configuration", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *job.ID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for trust store configuration", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.InstanceID.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *trustStoreConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state trustStoreConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep := time.Duration(state.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(state.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := r.client.ReadTrustStoreConfiguration(timeoutCtx, instanceID, sleep)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			tflog.Info(ctx, "Trust store configuration not found, removing resource")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading trust store configuration", err.Error())
		return
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("trust store configuration not found, resource will be recreated: %s", state.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Url = types.StringValue(data.Url)
	state.RefreshInterval = types.Int64Value(data.RefreshInterval)
	state.Provider = types.StringValue(data.Provider)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *trustStoreConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state trustStoreConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	changed := false
	params := model.TrustStoreRequest{}
	if plan.Url.ValueString() != state.Url.ValueString() {
		params.Url = plan.Url.ValueString()
		changed = true
	}
	if plan.RefreshInterval.ValueInt64() != state.RefreshInterval.ValueInt64() {
		params.RefreshInterval = plan.RefreshInterval.ValueInt64()
		changed = true
	}
	if plan.Provider.ValueString() != state.Provider.ValueString() {
		params.Provider = plan.Provider.ValueString()
		changed = true
	}

	if !changed {
		tflog.Info(ctx, "No changes detected for trust store configuration, skipping update")
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	job, err := r.client.UpdateTrustStoreConfiguration(timeoutCtx, instanceID, sleep, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating trust store configuration", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *job.ID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for trust store configuration", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *trustStoreConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state trustStoreConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep := time.Duration(state.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(state.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	job, err := r.client.DeleteTrustStoreConfiguration(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting trust store configuration", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *job.ID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for deleted trust store configuration", err.Error())
		return
	}
}
