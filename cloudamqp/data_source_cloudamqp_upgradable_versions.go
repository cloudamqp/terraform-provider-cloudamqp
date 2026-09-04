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
	_ datasource.DataSource              = &upgradableVersionsDataSource{}
	_ datasource.DataSourceWithConfigure = &upgradableVersionsDataSource{}
)

type upgradableVersionsDataSource struct {
	client *api.API
}

func NewUpgradableVersionsDataSource() datasource.DataSource {
	return &upgradableVersionsDataSource{}
}

type upgradableVersionsDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	InstanceID         types.Int64  `tfsdk:"instance_id"`
	NewRabbitMQVersion types.String `tfsdk:"new_rabbitmq_version"`
	NewErlangVersion   types.String `tfsdk:"new_erlang_version"`
}

func (d *upgradableVersionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_upgradable_versions"
}

func (d *upgradableVersionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about available RabbitMQ and Erlang upgradable versions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this data source",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
			},
			"new_rabbitmq_version": schema.StringAttribute{
				Computed:    true,
				Description: "Latest possible upgradable RabbitMQ version",
			},
			"new_erlang_version": schema.StringAttribute{
				Computed:    true,
				Description: "Latest possible upgradable Erlang version",
			},
		},
	}
}

func (d *upgradableVersionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *upgradableVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config upgradableVersionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := d.client.ReadVersions(ctx, config.InstanceID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read upgradable versions", err.Error())
		return
	}

	config.ID = types.StringValue(strconv.FormatInt(config.InstanceID.ValueInt64(), 10))

	if v, ok := data["new_rabbitmq_version"]; ok {
		config.NewRabbitMQVersion = types.StringValue(fmt.Sprintf("%v", v))
	}
	if v, ok := data["new_erlang_version"]; ok {
		config.NewErlangVersion = types.StringValue(fmt.Sprintf("%v", v))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
