package cloudamqp

import (
	"context"
	"fmt"
	"regexp"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &teamTagResource{}
	_ resource.ResourceWithConfigure   = &teamTagResource{}
	_ resource.ResourceWithImportState = &teamTagResource{}
)

type teamTagResource struct {
	client *api.API
}

func NewTeamTagResource() resource.Resource {
	return &teamTagResource{}
}

type teamTagResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	TagID types.Int64  `tfsdk:"tag_id"`
}

func (r *teamTagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_team_tag"
}

func (r *teamTagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a team tag, used to control which team members can access which " +
			"instances. Not supported for SAML-enabled teams.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The tag's name, used as its identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the tag. Changing this forces a new resource, since tags can't be renamed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 500),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^,]*$`), "must not contain a comma"),
				},
			},
			"tag_id": schema.Int64Attribute{
				Computed:    true,
				Description: "The tag's underlying numeric id, used internally to delete it.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *teamTagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *teamTagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *teamTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan teamTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.CreateTeamTag(ctx, model.TagRequest{Name: plan.Name.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Team Tag",
			fmt.Sprintf("Could not create team tag: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(data.Name)
	plan.TagID = types.Int64Value(int64(data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *teamTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags, err := r.client.ListTeamTags(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Team Tag",
			fmt.Sprintf("Could not list team tags: %s", err),
		)
		return
	}

	// No single-tag GET endpoint - find this tag by id (== name) in the list.
	id := state.ID.ValueString()
	var found *model.TagResponse
	for i := range tags {
		if tags[i].Name == id {
			found = &tags[i]
			break
		}
	}

	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(found.Name)
	state.Name = types.StringValue(found.Name)
	state.TagID = types.Int64Value(int64(found.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update only runs for in-place changes, and "name" is RequiresReplace - so
// there is nothing to do here in practice, but Framework requires the method.
func (r *teamTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan teamTagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *teamTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamTagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeamTag(ctx, int(state.TagID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Team Tag",
			fmt.Sprintf("Could not delete team tag %q: %s", state.ID.ValueString(), err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}
