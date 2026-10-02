package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	instanceModel "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &accountDataSource{}
	_ datasource.DataSourceWithConfigure = &accountDataSource{}
)

type accountDataSource struct {
	client *api.API
}

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Instances types.List   `tfsdk:"instances"`
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_account"
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	accountInstanceObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":     types.Int64Type,
		"name":   types.StringType,
		"plan":   types.StringType,
		"region": types.StringType,
		"tags":   types.ListType{ElemType: types.StringType},
	}}

	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about all instances associated with the account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The account identifier",
			},
			"instances": schema.ListAttribute{
				Computed:    true,
				ElementType: accountInstanceObjectType,
				Description: "List of instances for the account.",
			},
		},
	}
}

func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountDataSourceModel

	instances, err := d.client.ListInstances(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list instances", err.Error())
		return
	}

	accountInstanceObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":     types.Int64Type,
		"name":   types.StringType,
		"plan":   types.StringType,
		"region": types.StringType,
		"tags":   types.ListType{ElemType: types.StringType},
	}}

	values := make([]attr.Value, 0, len(instances))
	for _, instance := range instances {
		obj, diags := accountInstanceObjectValue(ctx, accountInstanceObjectType.AttrTypes, instance)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		values = append(values, obj)
	}

	instancesList, diags := types.ListValue(accountInstanceObjectType, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Instances = instancesList

	state.ID = types.StringValue("account")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func accountInstanceObjectValue(ctx context.Context, attrTypes map[string]attr.Type, instance instanceModel.InstanceResponse) (types.Object, diag.Diagnostics) {
	tags, diags := types.ListValueFrom(ctx, types.StringType, instance.Tags)
	if diags.HasError() {
		return types.Object{}, diags
	}

	return types.ObjectValue(attrTypes, map[string]attr.Value{
		"id":     types.Int64Value(instance.ID),
		"name":   types.StringValue(instance.Name),
		"plan":   types.StringValue(instance.Plan),
		"region": types.StringValue(instance.Region),
		"tags":   tags,
	})
}
