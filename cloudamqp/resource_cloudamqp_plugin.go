package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &pluginResource{}
	_ resource.ResourceWithConfigure   = &pluginResource{}
	_ resource.ResourceWithImportState = &pluginResource{}
)

type pluginResource struct {
	client *api.API
}

func NewPluginResource() resource.Resource {
	return &pluginResource{}
}

type pluginResourceModel struct {
	ID          types.String   `tfsdk:"id"`
	InstanceID  types.Int64    `tfsdk:"instance_id"`
	Name        types.String   `tfsdk:"name"`
	Enabled     types.Bool     `tfsdk:"enabled"`
	Description types.String   `tfsdk:"description"`
	Version     types.String   `tfsdk:"version"`
	Sleep       types.Int64    `tfsdk:"sleep"`
	Timeout     types.Int64    `tfsdk:"timeout"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

func (r *pluginResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_plugin"
}

func (r *pluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage RabbitMQ plugins for an instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (plugin name)",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the plugin",
			},
			"enabled": schema.BoolAttribute{
				Required:    true,
				Description: "If the plugin is enabled",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "The description of the plugin",
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "The version of the plugin",
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
				Description: "Configurable sleep time in seconds between retries for plugins",
			},
			"timeout": schema.Int64Attribute{
				Optional:           true,
				Computed:           true,
				Default:            int64default.StaticInt64(1800),
				Description:        "Configurable timeout time in seconds for plugins",
				DeprecationMessage: "The timeout attribute is deprecated and will be removed in next 2.0.0 version. Use the timeouts block instead.",
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create:            true,
				Read:              true,
				Update:            true,
				Delete:            true,
				CreateDescription: "Timeout for creating a plugin. Default is 60 minutes.",
				ReadDescription:   "Timeout for reading a plugin. Default is 60 minutes.",
				UpdateDescription: "Timeout for updating a plugin. Default is 60 minutes.",
				DeleteDescription: "Timeout for deleting a plugin. Default is 60 minutes.",
			}),
		},
	}
}

func (r *pluginResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *pluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("import of resource with identifiers: %s", req.ID))

	s := strings.Split(req.ID, ",")
	if len(s) != 2 {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {resource_id},{instance_id}")
		return
	}

	instanceID, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance_id in import ID", fmt.Sprintf("Could not convert instance_id to int64: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), s[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), s[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("sleep"), int64(10))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("timeout"), int64(1800))...)
}

func (r *pluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pluginResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	params := model.PluginRequest{
		Name:    plan.Name.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
	}
	sleep := plan.Sleep.ValueInt64()

	createTimeout, diags := plan.Timeouts.Create(ctx, 60*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	data, err := r.client.EnablePlugin(timeoutCtx, instanceID, params, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Create Plugin", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.Name.ValueString())
	plan.Description = types.StringValue(data.Description)
	plan.Version = types.StringValue(data.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pluginResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	name := state.Name.ValueString()
	sleep := state.Sleep.ValueInt64()

	createTimeout, diags := state.Timeouts.Read(ctx, 60*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	data, err := r.client.ReadPlugin(timeoutCtx, instanceID, name, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Read Plugin", err.Error())
		return
	}

	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("plugin not found, %s", state.Name.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Description = types.StringValue(data.Description)
	state.Version = types.StringValue(data.Version)
	state.Enabled = types.BoolValue(data.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *pluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pluginResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	params := model.PluginRequest{
		Name:    plan.Name.ValueString(),
		Enabled: plan.Enabled.ValueBool(),
	}
	sleep := plan.Sleep.ValueInt64()

	updateTimeout, diags := plan.Timeouts.Update(ctx, 60*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	data, err := r.client.UpdatePlugin(timeoutCtx, instanceID, params, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Update Plugin", err.Error())
		return
	}

	plan.Description = types.StringValue(data.Description)
	plan.Version = types.StringValue(data.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pluginResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if enableFasterInstanceDestroy {
		tflog.Debug(ctx, "cloudamqp::resource::plugin::delete skip calling backend.")
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	params := model.PluginRequest{
		Name:    state.Name.ValueString(),
		Enabled: state.Enabled.ValueBool(),
	}
	sleep := state.Sleep.ValueInt64()

	updateTimeout, diags := state.Timeouts.Update(ctx, 60*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	_, err := r.client.DeletePlugin(timeoutCtx, instanceID, params, sleep)
	if err != nil {
		if strings.Contains(err.Error(), "instance not found") || strings.Contains(err.Error(), "status=404") {
			tflog.Info(ctx, fmt.Sprintf("instance not found during plugin deletion, considering successful: %s", state.Name.ValueString()))
			return
		}
		resp.Diagnostics.AddError("Failed to Delete Plugin", err.Error())
	}
}
