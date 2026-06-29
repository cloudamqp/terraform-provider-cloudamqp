package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	ID        types.String           `tfsdk:"id"`
	Instances []accountInstanceModel `tfsdk:"instances"`
}

type accountInstanceModel struct {
	ID     types.Int64  `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Plan   types.String `tfsdk:"plan"`
	Region types.String `tfsdk:"region"`
	Tags   types.List   `tfsdk:"tags"`
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_account"
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about all instances associated with the account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The account identifier",
			},
		},
		Blocks: map[string]schema.Block{
			"instances": schema.SetNestedBlock{
				Description: "List of instances for the account.",
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
						"plan": schema.StringAttribute{
							Computed:    true,
							Description: "The subscription plan used for the instance",
						},
						"region": schema.StringAttribute{
							Computed:    true,
							Description: "The region where the instance is located",
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Description: "Tag for the instance",
						},
					},
				},
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

	state.Instances = make([]accountInstanceModel, 0, len(instances))
	for _, instance := range instances {
		instanceState := accountInstanceModel{}
		instanceState.ID = types.Int64Value(instance.ID)
		instanceState.Name = types.StringValue(instance.Name)
		instanceState.Plan = types.StringValue(instance.Plan)
		instanceState.Region = types.StringValue(instance.Region)
		instanceState.Tags, _ = types.ListValueFrom(ctx, types.StringType, instance.Tags)
		state.Instances = append(state.Instances, instanceState)
	}

	state.ID = types.StringValue("account")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
