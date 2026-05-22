package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &vpcInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcInfoDataSource{}
)

type vpcInfoDataSource struct {
	client *api.API
}

func NewVpcInfoDataSource() datasource.DataSource {
	return &vpcInfoDataSource{}
}

type vpcInfoDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	InstanceID      types.Int64  `tfsdk:"instance_id"`
	VpcID           types.String `tfsdk:"vpc_id"`
	Name            types.String `tfsdk:"name"`
	VpcSubnet       types.String `tfsdk:"vpc_subnet"`
	OwnerID         types.String `tfsdk:"owner_id"`
	SecurityGroupID types.String `tfsdk:"security_group_id"`
}

func (d *vpcInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_vpc_info"
}

func (d *vpcInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve AWS VPC peering information. Either use instance_id or vpc_id to retrieve the VPC info.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this data source",
			},
			"instance_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Instance identifier",
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "VPC instance identifier",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "VPC name",
			},
			"vpc_subnet": schema.StringAttribute{
				Computed:    true,
				Description: "VPC subnet",
			},
			"owner_id": schema.StringAttribute{
				Computed:    true,
				Description: "Owner identifier",
			},
			"security_group_id": schema.StringAttribute{
				Computed:    true,
				Description: "The security group identifier",
			},
		},
	}
}

func (d *vpcInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpcInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config vpcInfoDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.InstanceID.IsNull() && config.VpcID.IsNull() {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"You need to specify either instance_id or vpc_id.",
		)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(1800)*time.Second)
	defer cancel()

	var (
		data *model.VpcInfoResponse
		err  error
	)

	if !config.InstanceID.IsNull() {
		data, err = d.client.ReadVpcInfo(timeoutCtx, config.InstanceID.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read VPC info",
				fmt.Sprintf("Could not read VPC info for instance %d: %s", config.InstanceID.ValueInt64(), err.Error()),
			)
			return
		}
	} else {
		data, err = d.client.ReadVpcInfoWithVpcId(timeoutCtx, config.VpcID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read VPC info",
				fmt.Sprintf("Could not read VPC info for VPC %s: %s", config.VpcID.ValueString(), err.Error()),
			)
			return
		}
	}

	if data == nil || data.Id == "" {
		resp.Diagnostics.AddError(
			"Failed to find VPC identifier",
			"Data source is used for AWS VPC peering.",
		)
		return
	}

	config.ID = types.StringValue(data.Id)
	config.Name = types.StringValue(data.Name)
	config.VpcSubnet = types.StringValue(data.Subnet)
	config.OwnerID = types.StringValue(data.OwnerId)
	config.SecurityGroupID = types.StringValue(data.SecurityGroupId.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
