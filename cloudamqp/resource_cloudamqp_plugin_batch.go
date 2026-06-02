package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	instancemodel "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource              = &pluginBatchResource{}
	_ resource.ResourceWithConfigure = &pluginBatchResource{}
)

type pluginBatchResource struct {
	client *api.API
}

func NewPluginBatchResource() resource.Resource {
	return &pluginBatchResource{}
}

type pluginBatchResourceModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	Plugins    types.Map    `tfsdk:"plugins"`
	Sleep      types.Int64  `tfsdk:"sleep"`
	Timeout    types.Int64  `tfsdk:"timeout"`
}

func (r *pluginBatchResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_plugin_batch"
}

func (r *pluginBatchResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage multiple RabbitMQ plugins in a single batch operation.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "The CloudAMQP instance identifier.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"plugins": schema.MapAttribute{
				Required:    true,
				ElementType: types.BoolType,
				Description: "Map of plugin name to enabled state (true = enabled, false = disabled).",
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
				Description: "Configurable sleep time in seconds between retries for plugin operations (default: 10).",
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1800),
				Description: "Configurable timeout time in seconds for plugin operations (default: 1800).",
			},
		},
	}
}

func (r *pluginBatchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *pluginBatchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pluginBatchResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep := plan.Sleep.ValueInt64()
	timeout := plan.Timeout.ValueInt64()

	planPlugins, diags := pluginsMapToBoolMap(ctx, plan.Plugins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var enableList []string
	for name, enabled := range planPlugins {
		if enabled {
			enableList = append(enableList, name)
		}
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	jobResp, err := r.client.CreatePluginBatch(timeoutCtx, instanceID, instancemodel.PluginBatchRequest{
		Enable: enableList,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating plugin batch", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *jobResp.ID, time.Duration(sleep)*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for plugin batch job", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d", instanceID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pluginBatchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pluginBatchResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(state.InstanceID.ValueInt64())
	sleep := int(state.Sleep.ValueInt64())
	timeout := int(state.Timeout.ValueInt64())

	statePlugins, diags := pluginsMapToBoolMap(ctx, state.Plugins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiPlugins, err := r.client.ListPlugins(ctx, instanceID, sleep, timeout)
	if err != nil {
		resp.Diagnostics.AddError("Error reading plugins", err.Error())
		return
	}

	// Index API response by plugin name.
	apiIndex := make(map[string]bool, len(apiPlugins))
	for _, p := range apiPlugins {
		name, ok := p["name"].(string)
		if !ok {
			continue
		}
		enabled, _ := p["enabled"].(bool)
		apiIndex[name] = enabled
	}

	// Rebuild the plugins map using only the keys tracked in state.
	updatedPlugins := make(map[string]attr.Value, len(statePlugins))
	for name := range statePlugins {
		if apiEnabled, exists := apiIndex[name]; exists {
			updatedPlugins[name] = types.BoolValue(apiEnabled)
		} else {
			// Plugin no longer exists in API; remove from state by not adding it.
			continue
		}
	}

	refreshed, diags := types.MapValue(types.BoolType, updatedPlugins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Plugins = refreshed
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *pluginBatchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state pluginBatchResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep := plan.Sleep.ValueInt64()
	timeout := plan.Timeout.ValueInt64()

	planPlugins, diags := pluginsMapToBoolMap(ctx, plan.Plugins)
	resp.Diagnostics.Append(diags...)
	statePlugins, diags := pluginsMapToBoolMap(ctx, state.Plugins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var enableList, disableList []string

	// Enable: in plan as true, and either absent from state or previously false.
	for name, planEnabled := range planPlugins {
		if planEnabled {
			if stateEnabled, exists := statePlugins[name]; !exists || !stateEnabled {
				enableList = append(enableList, name)
			}
		}
	}

	// Disable: absent from plan or set to false in plan, and was true in state.
	for name, stateEnabled := range statePlugins {
		if !stateEnabled {
			continue
		}
		planEnabled, inPlan := planPlugins[name]
		if !inPlan || !planEnabled {
			disableList = append(disableList, name)
		}
	}

	// Nothing changed — write plan to state and return.
	if len(enableList) == 0 && len(disableList) == 0 {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	jobResp, err := r.client.UpdatePluginBatch(timeoutCtx, instanceID, instancemodel.PluginBatchRequest{
		Enable:  enableList,
		Disable: disableList,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating plugin batch", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *jobResp.ID, time.Duration(sleep)*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for plugin batch update job", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d", instanceID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pluginBatchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state pluginBatchResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep := state.Sleep.ValueInt64()
	timeout := state.Timeout.ValueInt64()

	statePlugins, diags := pluginsMapToBoolMap(ctx, state.Plugins)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var disableList []string
	for name, enabled := range statePlugins {
		if enabled {
			disableList = append(disableList, name)
		}
	}

	// Nothing to disable.
	if len(disableList) == 0 {
		return
	}

	if enableFasterInstanceDestroy {
		tflog.Debug(ctx, "cloudamqp::resource::plugin_batch::delete skip calling backend.")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	jobResp, err := r.client.DeletePluginBatch(timeoutCtx, instanceID, instancemodel.PluginBatchRequest{
		Disable: disableList,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting plugin batch", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *jobResp.ID, time.Duration(sleep)*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for plugin batch delete job", err.Error())
		return
	}
}

// pluginsMapToBoolMap converts a types.Map (BoolType elements) to map[string]bool.
func pluginsMapToBoolMap(ctx context.Context, m types.Map) (map[string]bool, diag.Diagnostics) {
	elements := make(map[string]types.Bool, len(m.Elements()))
	diags := m.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make(map[string]bool, len(elements))
	for k, v := range elements {
		result[k] = v.ValueBool()
	}
	return result, diags
}
