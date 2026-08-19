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
)

var (
	_ datasource.DataSource              = &pluginsCommunityDataSource{}
	_ datasource.DataSourceWithConfigure = &pluginsCommunityDataSource{}
)

type pluginsCommunityDataSource struct {
	client *api.API
}

func NewPluginsCommunityDataSource() datasource.DataSource {
	return &pluginsCommunityDataSource{}
}

type pluginsCommunityDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	Plugins    types.List   `tfsdk:"plugins"`
	Sleep      types.Int64  `tfsdk:"sleep"`
	Timeout    types.Int64  `tfsdk:"timeout"`
}

func (d *pluginsCommunityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_plugins_community"
}

func (d *pluginsCommunityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	pluginObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"require":     types.StringType,
		"description": types.StringType,
	}}

	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve community plugins for an instance.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The data source identifier",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
			},
			"plugins": schema.ListAttribute{
				Computed:    true,
				ElementType: pluginObjectType,
				Description: "List of community plugins",
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

func (d *pluginsCommunityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *pluginsCommunityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state pluginsCommunityDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep := int64(10)
	if !state.Sleep.IsNull() && !state.Sleep.IsUnknown() {
		sleep = state.Sleep.ValueInt64()
	}
	timeout := int64(1800)
	if !state.Timeout.IsNull() && !state.Timeout.IsUnknown() {
		timeout = state.Timeout.ValueInt64()
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	data, err := d.client.ListPluginsCommunity(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Read Community Plugins", err.Error())
		return
	}

	pluginObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"require":     types.StringType,
		"description": types.StringType,
	}}

	values := make([]attr.Value, 0, len(data))
	for _, p := range data {
		obj, diags := types.ObjectValue(pluginObjectType.AttrTypes, map[string]attr.Value{
			"name":        types.StringValue(p.Name),
			"require":     types.StringValue(p.Require),
			"description": types.StringValue(p.Description),
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

	state.ID = types.StringValue(fmt.Sprintf("%d.plugins_community", instanceID))
	state.Plugins = pluginsList
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
