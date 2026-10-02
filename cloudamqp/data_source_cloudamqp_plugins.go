package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &pluginsDataSource{}
	_ datasource.DataSourceWithConfigure = &pluginsDataSource{}
)

type pluginsDataSource struct {
	client *api.API
}

func NewPluginsDataSource() datasource.DataSource {
	return &pluginsDataSource{}
}

type pluginsDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	InstanceID  types.Int64  `tfsdk:"instance_id"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Recommended types.Bool   `tfsdk:"recommended"`
	Required    types.Bool   `tfsdk:"required"`
	Plugins     types.List   `tfsdk:"plugins"`
	Sleep       types.Int64  `tfsdk:"sleep"`
	Timeout     types.Int64  `tfsdk:"timeout"`
}

func (d *pluginsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_plugins"
}

func (d *pluginsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	pluginObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"version":     types.StringType,
		"description": types.StringType,
		"enabled":     types.BoolType,
		"recommended": types.StringType,
		"required":    types.StringType,
	}}

	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve plugins and their status for an instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source identifier",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Only include enabled plugins in state",
			},
			"recommended": schema.BoolAttribute{
				Optional:    true,
				Description: "Only include recommended plugins in state",
			},
			"required": schema.BoolAttribute{
				Optional:    true,
				Description: "Only include required plugins in state",
			},
			"plugins": schema.ListAttribute{
				Computed:    true,
				ElementType: pluginObjectType,
				Description: "List of plugins",
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable sleep time in seconds between retries for plugins",
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable timeout time in seconds for plugins",
			},
		},
	}
}

func (d *pluginsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *pluginsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pluginsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	enabledFilter := !state.Enabled.IsNull() && state.Enabled.ValueBool()
	recommendedFilter := !state.Recommended.IsNull() && state.Recommended.ValueBool()
	requiredFilter := !state.Required.IsNull() && state.Required.ValueBool()
	sleep := int64(10)
	if !state.Sleep.IsNull() && !state.Sleep.IsUnknown() {
		sleep = state.Sleep.ValueInt64()
	}
	timeout := int64(1800)
	if !state.Timeout.IsNull() && !state.Timeout.IsUnknown() {
		timeout = state.Timeout.ValueInt64()
	}

	tflog.Info(ctx, fmt.Sprintf("Reading plugins for instanceID=%d with filters enabled=%t, recommended=%t, required=%t, sleep=%d, timeout=%d", instanceID, enabledFilter, recommendedFilter, requiredFilter, sleep, timeout))

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	data, err := d.client.ListPlugins(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Read Plugins", err.Error())
		return
	}

	pluginObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"version":     types.StringType,
		"description": types.StringType,
		"enabled":     types.BoolType,
		"recommended": types.StringType,
		"required":    types.StringType,
	}}

	values := make([]attr.Value, 0, len(data))
	for _, p := range data {
		recommendedBool := p.Recommended != nil && *p.Recommended != "" && *p.Recommended != "false"
		requiredBool := p.Required != nil && *p.Required != "" && *p.Required != "false"

		if enabledFilter && !p.Enabled {
			continue
		}
		if recommendedFilter && !recommendedBool {
			continue
		}
		if requiredFilter && !requiredBool {
			continue
		}

		recommendedVal := ""
		if p.Recommended != nil {
			recommendedVal = *p.Recommended
		}
		requiredVal := ""
		if p.Required != nil {
			requiredVal = *p.Required
		}

		obj, diags := types.ObjectValue(pluginObjectType.AttrTypes, map[string]attr.Value{
			"name":        types.StringValue(p.Name),
			"version":     types.StringValue(p.Version),
			"description": types.StringValue(p.Description),
			"enabled":     types.BoolValue(p.Enabled),
			"recommended": types.StringValue(recommendedVal),
			"required":    types.StringValue(requiredVal),
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		values = append(values, obj)
	}

	pluginsList, diags := types.ListValue(pluginObjectType, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Sleep = types.Int64Value(sleep)
	state.Timeout = types.Int64Value(timeout)
	state.ID = types.StringValue(fmt.Sprintf("%d.plugins", instanceID))
	state.Plugins = pluginsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
