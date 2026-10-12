package cloudamqp

import (
	"context"
	"fmt"
	"sort"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &teamTagsDataSource{}
var _ datasource.DataSourceWithConfigure = &teamTagsDataSource{}

type teamTagsDataSource struct {
	client *api.API
}

func NewTeamTagsDataSource() datasource.DataSource {
	return &teamTagsDataSource{}
}

type teamTagsDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Tags types.List   `tfsdk:"tags"`
}

func teamTagObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"name": types.StringType,
	}}
}

func (d *teamTagsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_team_tags"
}

func (d *teamTagsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all tags for the team, e.g. to reference existing tags without managing them " +
			"via cloudamqp_team_tag.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The team identifier",
			},
			"tags": schema.ListAttribute{
				Computed:    true,
				ElementType: teamTagObjectType(),
				Description: "List of tags for the team.",
			},
		},
	}
}

func (d *teamTagsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *teamTagsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state teamTagsDataSourceModel

	tags, err := d.client.ListTeamTags(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Failed to list team tags: %s", err.Error()))
		return
	}

	sort.Slice(tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })

	objectType := teamTagObjectType()
	values := make([]attr.Value, 0, len(tags))
	for _, tag := range tags {
		obj, diags := teamTagObjectValue(objectType.AttrTypes, tag)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		values = append(values, obj)
	}

	tagsList, diags := types.ListValue(objectType, values)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Tags = tagsList

	state.ID = types.StringValue("team_tags")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func teamTagObjectValue(attrTypes map[string]attr.Type, tag model.TagResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(attrTypes, map[string]attr.Value{
		"name": types.StringValue(tag.Name),
	})
}
