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
	_ datasource.DataSource              = &vpcGcpInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &vpcGcpInfoDataSource{}
)

type vpcGcpInfoDataSource struct {
	client *api.API
}

func NewVpcGcpInfoDataSource() datasource.DataSource {
	return &vpcGcpInfoDataSource{}
}

type vpcGcpInfoDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	VpcID      types.String `tfsdk:"vpc_id"`
	Name       types.String `tfsdk:"name"`
	VpcSubnet  types.String `tfsdk:"vpc_subnet"`
	Network    types.String `tfsdk:"network"`
	Sleep      types.Int64  `tfsdk:"sleep"`
	Timeout    types.Int64  `tfsdk:"timeout"`
}

func (d *vpcGcpInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_vpc_gcp_info"
}

func (d *vpcGcpInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve GCP VPC peering information. Either use instance_id or vpc_id to retrieve the VPC info.",
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
			"network": schema.StringAttribute{
				Computed:    true,
				Description: "VPC network uri",
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable sleep in seconds between retries when reading peering",
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable timeout time (seconds) before retries times out",
			},
		},
	}
}

func (d *vpcGcpInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vpcGcpInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config vpcGcpInfoDataSourceModel
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

	sleep := int64(10)
	if !config.Sleep.IsNull() && !config.Sleep.IsUnknown() {
		sleep = config.Sleep.ValueInt64()
	}
	timeout := int64(1800)
	if !config.Timeout.IsNull() && !config.Timeout.IsUnknown() {
		timeout = config.Timeout.ValueInt64()
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	var (
		data *model.VpcGcpInfoResponse
		err  error
	)

	if !config.InstanceID.IsNull() {
		data, err = d.client.ReadVpcGcpInfo(timeoutCtx, config.InstanceID.ValueInt64(), sleep)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read VPC GCP info",
				fmt.Sprintf("Could not read VPC GCP info for instance %d: %s", config.InstanceID.ValueInt64(), err.Error()),
			)
			return
		}
	} else {
		data, err = d.client.ReadVpcGcpInfoWithVpcId(timeoutCtx, config.VpcID.ValueString(), sleep)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to read VPC GCP info",
				fmt.Sprintf("Could not read VPC GCP info for VPC %s: %s", config.VpcID.ValueString(), err.Error()),
			)
			return
		}
	}

	config.ID = types.StringValue(data.Name)
	config.Name = types.StringValue(data.Name)
	config.Network = types.StringValue(data.Network)
	config.VpcSubnet = types.StringValue(data.Subnet)
	config.Sleep = types.Int64Value(sleep)
	config.Timeout = types.Int64Value(timeout)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
