package cloudamqp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &instanceDataSource{}
	_ datasource.DataSourceWithConfigure = &instanceDataSource{}
)

type instanceDataSource struct {
	client *api.API
}

func NewInstanceDataSource() datasource.DataSource {
	return &instanceDataSource{}
}

type instanceDataSourceModel struct {
	InstanceID      types.Int64  `tfsdk:"instance_id"`
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Plan            types.String `tfsdk:"plan"`
	Region          types.String `tfsdk:"region"`
	VpcID           types.Int64  `tfsdk:"vpc_id"`
	VpcSubnet       types.String `tfsdk:"vpc_subnet"`
	Nodes           types.Int64  `tfsdk:"nodes"`
	RmqVersion      types.String `tfsdk:"rmq_version"`
	Tags            types.List   `tfsdk:"tags"`
	Url             types.String `tfsdk:"url"`
	ApiKey          types.String `tfsdk:"apikey"`
	Host            types.String `tfsdk:"host"`
	HostInternal    types.String `tfsdk:"host_internal"`
	Vhost           types.String `tfsdk:"vhost"`
	Ready           types.Bool   `tfsdk:"ready"`
	Dedicated       types.Bool   `tfsdk:"dedicated"`
	NoDefaultAlarms types.Bool   `tfsdk:"no_default_alarms"`
	Backend         types.String `tfsdk:"backend"`
}

func (d *instanceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_instance"
}

func (d *instanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about an existing CloudAMQP instance.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this data source instance",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the instance",
			},
			"plan": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the plan",
			},
			"region": schema.StringAttribute{
				Computed:    true,
				Description: "Name of the region the instance is located in",
			},
			"vpc_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The VPC ID",
			},
			"vpc_subnet": schema.StringAttribute{
				Computed:    true,
				Description: "The VPC subnet",
			},
			"nodes": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of nodes in cluster",
			},
			"rmq_version": schema.StringAttribute{
				Computed:    true,
				Description: "RabbitMQ/LavinMQ version",
			},
			"tags": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tags associated with the instance",
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "URL of the CloudAMQP instance",
			},
			"apikey": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "API key for the CloudAMQP instance",
			},
			"host": schema.StringAttribute{
				Computed:    true,
				Description: "External hostname for the CloudAMQP instance",
			},
			"host_internal": schema.StringAttribute{
				Computed:    true,
				Description: "Internal hostname for the CloudAMQP instance",
			},
			"vhost": schema.StringAttribute{
				Computed:    true,
				Description: "The virtual host",
			},
			"ready": schema.BoolAttribute{
				Computed:    true,
				Description: "Flag describing if the resource is ready",
			},
			"dedicated": schema.BoolAttribute{
				Computed:    true,
				Description: "Is the instance hosted on a dedicated server",
			},
			"no_default_alarms": schema.BoolAttribute{
				Computed:    true,
				Description: "If default alarms are set or not for the instance",
			},
			"backend": schema.StringAttribute{
				Computed:    true,
				Description: "Software backend used, either 'rabbitmq' or 'lavinmq'",
			},
		},
	}
}

func (d *instanceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *instanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config instanceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := strconv.FormatInt(config.InstanceID.ValueInt64(), 10)

	data, err := d.client.ReadInstance(ctx, instanceID, 0)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Instance",
			fmt.Sprintf("Could not read instance with ID %s: %s", instanceID, err),
		)
		return
	}

	if data == nil {
		resp.Diagnostics.AddError(
			"Instance Not Found",
			fmt.Sprintf("Instance with ID %s does not exist", instanceID),
		)
		return
	}

	config.ID = types.StringValue(instanceID)
	config.Name = types.StringValue(data.Name)
	config.Plan = types.StringValue(data.Plan)
	config.Region = types.StringValue(data.Region)
	config.Nodes = types.Int64Value(data.Nodes)
	config.RmqVersion = types.StringValue(data.RmqVersion)
	config.Url = types.StringValue(data.Url)
	config.ApiKey = types.StringValue(data.ApiKey)
	config.Ready = types.BoolValue(data.Ready)
	config.Backend = types.StringValue(data.Backend)
	config.Dedicated = types.BoolValue(data.Nodes > 0)
	config.Host = types.StringValue(data.HostnameExternal)
	config.HostInternal = types.StringValue(data.HostnameInternal)
	config.NoDefaultAlarms = types.BoolValue(false)

	if data.Url != "" {
		urlInfo := d.client.UrlInformation(data.Url)
		if vhost, ok := urlInfo["vhost"].(string); ok {
			config.Vhost = types.StringValue(vhost)
		}
	}

	if data.VPC != nil {
		config.VpcID = types.Int64Value(data.VPC.ID)
		if data.VPC.Subnet != "" {
			config.VpcSubnet = types.StringValue(data.VPC.Subnet)
		}
	} else {
		config.VpcID = types.Int64Null()
		config.VpcSubnet = types.StringNull()
	}

	if len(data.Tags) > 0 {
		tags, diags := types.ListValueFrom(ctx, types.StringType, data.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		config.Tags = tags
	} else {
		config.Tags = types.ListNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
