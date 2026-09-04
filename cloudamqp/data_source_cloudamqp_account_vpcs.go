package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	ID   types.String                `tfsdk:"id"`
	VPCs []accountVpcDataSourceModel `tfsdk:"vpcs"`
}

type accountVpcDataSourceModel struct {
	ID      types.Int64  `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Region  types.String `tfsdk:"region"`
	Subnet  types.String `tfsdk:"subnet"`
	Tags    types.List   `tfsdk:"tags"`
	VpcName types.String `tfsdk:"vpc_name"`
}

func (d *accountVpcsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_account_vpcs"
}

func (d *accountVpcsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The account identifier",
			},
		},
		Blocks: map[string]schema.Block{
			"vpcs": schema.SetNestedBlock{
				Description: "List of VPCs for the account.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:    true,
							Description: "The instance identifier",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the instance",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "The region where the instance is located",
						},
						"subnet": schema.StringAttribute{
							Computed:    true,
							Description: "The VPC subnet",
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Optional tags to associate with the VPC instance",
						},
						"vpc_name": schema.StringAttribute{
							Computed:    true,
							Description: "VPC name given when hosted at the cloud provider",
						},
					},
				},
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

	state.VPCs = make([]accountVpcDataSourceModel, 0, len(vpcs))
	for _, vpc := range vpcs {
		vpcState := accountVpcDataSourceModel{}
		vpcState.ID = types.Int64Value(vpc.ID)
		vpcState.Name = types.StringValue(vpc.Name)
		vpcState.Region = types.StringValue(vpc.Region)
		vpcState.Subnet = types.StringValue(vpc.Subnet)
		vpcState.Tags, _ = types.ListValueFrom(ctx, types.StringType, vpc.Tags)
		vpcState.VpcName = types.StringValue(vpc.VpcName)
		state.VPCs = append(state.VPCs, vpcState)
	}

	state.ID = types.StringValue("account_vpcs")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
