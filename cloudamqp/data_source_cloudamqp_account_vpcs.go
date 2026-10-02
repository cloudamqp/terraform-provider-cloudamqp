package cloudamqp

import (
	"context"
	"fmt"
	"sort"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	networkModel "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &accountVpcsDataSource{}
var _ datasource.DataSourceWithConfigure = &accountVpcsDataSource{}

type accountVpcsDataSource struct {
	client *api.API
}

func NewAccountVpcsDataSource() datasource.DataSource {
	return &accountVpcsDataSource{}
}

type accountVpcsDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	VPCs types.List   `tfsdk:"vpcs"`
}

func (d *accountVpcsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_account_vpcs"
}

func (d *accountVpcsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	accountVpcObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":       types.Int64Type,
		"name":     types.StringType,
		"region":   types.StringType,
		"subnet":   types.StringType,
		"tags":     types.ListType{ElemType: types.StringType},
		"vpc_name": types.StringType,
	}}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The account identifier",
			},
			"vpcs": schema.ListAttribute{
				Computed:    true,
				ElementType: accountVpcObjectType,
				Description: "List of VPCs for the account.",
			},
		},
	}
}

func (d *accountVpcsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *api.API, got: %T", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *accountVpcsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state accountVpcsDataSourceModel

	vpcs, err := d.client.ListVpcs(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to list VPCs: %s", err.Error()))
		return
	}

	accountVpcObjectType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"id":       types.Int64Type,
		"name":     types.StringType,
		"region":   types.StringType,
		"subnet":   types.StringType,
		"tags":     types.ListType{ElemType: types.StringType},
		"vpc_name": types.StringType,
	}}

	sort.Slice(vpcs, func(i, j int) bool { return vpcs[i].ID < vpcs[j].ID })

	values := make([]attr.Value, 0, len(vpcs))
	for _, vpc := range vpcs {
		obj, diags := accountVpcObjectValue(ctx, accountVpcObjectType.AttrTypes, vpc)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		values = append(values, obj)
	}

	vpcsList, diags := types.ListValue(accountVpcObjectType, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.VPCs = vpcsList

	state.ID = types.StringValue("account_vpcs")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func accountVpcObjectValue(ctx context.Context, attrTypes map[string]attr.Type, vpc networkModel.VpcResponse) (types.Object, diag.Diagnostics) {
	tags, diags := types.ListValueFrom(ctx, types.StringType, vpc.Tags)
	if diags.HasError() {
		return types.Object{}, diags
	}

	return types.ObjectValue(attrTypes, map[string]attr.Value{
		"id":       types.Int64Value(vpc.ID),
		"name":     types.StringValue(vpc.Name),
		"region":   types.StringValue(vpc.Region),
		"subnet":   types.StringValue(vpc.Subnet),
		"tags":     tags,
		"vpc_name": types.StringValue(vpc.VpcName),
	})
}
