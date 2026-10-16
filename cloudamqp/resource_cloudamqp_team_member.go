package cloudamqp

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// assignableTeamRoleNames are the base roles a member can be assigned.
var assignableTeamRoleNames = []string{
	"admin",
	"billing manager",
	"member",
	"monitor",
	"devops",
}

// taggedTeamRoleNames are the roles that require at least one tag
var taggedTeamRoleNames = []string{"member", "monitor"}

var (
	_ resource.Resource                   = &teamMemberResource{}
	_ resource.ResourceWithConfigure      = &teamMemberResource{}
	_ resource.ResourceWithImportState    = &teamMemberResource{}
	_ resource.ResourceWithValidateConfig = &teamMemberResource{}
)

type teamMemberResource struct {
	client *api.API
}

func NewTeamMemberResource() resource.Resource {
	return &teamMemberResource{}
}

type teamMemberResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Email    types.String `tfsdk:"email"`
	Role     types.String `tfsdk:"role"`
	Tags     types.List   `tfsdk:"tags"`
	Status   types.String `tfsdk:"status"`
	InviteID types.String `tfsdk:"invite_id"`
	UserID   types.String `tfsdk:"user_id"`
}

func (r *teamMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_team_member"
}

func (r *teamMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a team member invite/membership as a single resource across both states: " +
			"a pending invite before it's accepted, and an accepted member afterwards. \"status\" reports " +
			"which one currently applies, it's computed, not configured.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The member's normalized (lowercased) email, used as this resource's identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "Email address to invite. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Role assigned to the member. Defaults to \"member\" if omitted. One of: " +
					strings.Join(assignableTeamRoleNames, ", ") + ". For SAML-configured teams, leave " +
					"this unset to let it track whatever the IdP assigns without producing a perpetual diff.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(assignableTeamRoleNames...),
				},
			},
			"tags": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tag names granting access to tagged instances. A tag that doesn't exist yet " +
					"is created automatically. Only meaningful for the \"member\"/\"monitor\" roles, " +
					"which require at least one tag.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "\"pending\" while the invite hasn't been accepted yet, \"member\" afterwards",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invite_id": schema.StringAttribute{
				Computed:    true,
				Description: "The pending invite's id. Null once status is \"member\".",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Computed:    true,
				Description: "The accepted member's user id. Null while status is \"pending\".",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *teamMemberResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config teamMemberResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Role.IsUnknown() {
		return
	}

	// A null role in config defaults to "member" server-side.
	role := "member"
	if !config.Role.IsNull() {
		role = config.Role.ValueString()
	}
	if !slices.Contains(taggedTeamRoleNames, role) {
		return
	}

	// Tags referencing another resource's output can't be checked yet either.
	if config.Tags.IsUnknown() {
		return
	}

	if config.Tags.IsNull() || len(config.Tags.Elements()) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("tags"),
			"Missing Required Tags",
			fmt.Sprintf("The %q role requires at least one tag to be specified.", role),
		)
	}
}

func (r *teamMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *teamMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), normalizeEmail(req.ID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), req.ID)...)
}

func (r *teamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan teamMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	displayEmail := plan.Email.ValueString()
	email := normalizeEmail(displayEmail)

	tags := make([]string, 0)
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// POST /api/team/invite joins tags with "," and the backend splits on ",", so a tag
	// name containing a comma would silently become two different tags.
	for _, tag := range tags {
		if strings.Contains(tag, ",") {
			resp.Diagnostics.AddError(
				"Invalid Tag Name",
				fmt.Sprintf("Tag %q contains a comma, which isn't supported when inviting a member.", tag),
			)
			return
		}
	}

	role := ""
	if !plan.Role.IsNull() && !plan.Role.IsUnknown() {
		role = plan.Role.ValueString()
	}

	invite, err := r.client.InviteTeamMember(ctx, model.InviteRequest{
		Email: email,
		Role:  role,
		Tags:  strings.Join(tags, ","),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Team Member",
			fmt.Sprintf("Could not invite %q: %s", email, err),
		)
		return
	}

	// POST /api/team/invite defaults an omitted role to "member" server-side.
	if role == "" {
		role = "member"
	}

	state, diags := inviteState(ctx, displayEmail, &model.InviteByEmailResponse{
		ID:   invite.ID,
		Role: role,
		Tags: tags,
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *teamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, diags := r.readState(ctx, state.Email.ValueString(), &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *teamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state teamMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	email := state.Email.ValueString()

	tags := make([]string, 0)
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	role := ""
	if !plan.Role.IsNull() && !plan.Role.IsUnknown() {
		role = plan.Role.ValueString()
	}

	var err error
	if state.Status.ValueString() == "pending" {
		err = r.client.UpdateTeamInvite(ctx, state.InviteID.ValueString(), model.InviteUpdateRequest{
			Role: role,
			Tags: tags,
		})
	} else {
		err = r.client.UpdateTeamMember(ctx, state.UserID.ValueString(), model.MemberUpdateRequest{
			Role: role,
			Tags: tags,
		})
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Team Member",
			fmt.Sprintf("Could not update %q: %s", email, err),
		)
		return
	}

	newState, diags := r.readState(ctx, email, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if newState == nil {
		resp.Diagnostics.AddError(
			"Failed to Update Team Member",
			fmt.Sprintf("%q could not be found after updating", email),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *teamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if state.Status.ValueString() == "pending" {
		err = r.client.DeleteTeamInvite(ctx, state.InviteID.ValueString())
	} else {
		err = r.client.RemoveTeamMember(ctx, state.Email.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Team Member",
			fmt.Sprintf("Could not remove %q: %s", state.Email.ValueString(), err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

// readState uses the previously known status, if any, to do a single direct
// lookup by id instead of searching by email:
//   - status "member": GET /api/team/:user_id. A member can't go back to
//     pending, so a 404 here means genuinely removed, no fallback needed.
//   - status "pending": GET /api/team/invites/:id. A 404 here means the
//     invite was accepted since the last read, so we fall back to a member
//     lookup by email, the only remaining possibility, to pick up the new
//     user_id.
//   - no previous status (Create, or anything unexpected): cold start, check
//     pending invites by email first (the state right after inviting), then
//     fall back to an accepted member.
//
// previous is nil on Create, where there's no id yet to look up by.
func (r *teamMemberResource) readState(
	ctx context.Context, displayEmail string, previous *teamMemberResourceModel,
) (*teamMemberResourceModel, diag.Diagnostics) {
	email := normalizeEmail(displayEmail)

	if previous != nil {
		switch previous.Status.ValueString() {
		case "member":
			return r.readMemberByID(ctx, displayEmail, email, previous.UserID.ValueString())
		case "pending":
			state, diags := r.readInviteByID(ctx, displayEmail, previous.InviteID.ValueString())
			if diags.HasError() || state != nil {
				return state, diags
			}
			return r.readMemberByEmail(ctx, displayEmail, email)
		}
	}

	state, diags := r.readInviteByEmail(ctx, displayEmail, email)
	if diags.HasError() || state != nil {
		return state, diags
	}
	return r.readMemberByEmail(ctx, displayEmail, email)
}

func (r *teamMemberResource) readMemberByID(
	ctx context.Context, displayEmail, email, userID string,
) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	member, err := r.client.GetTeamMember(ctx, userID)
	if err != nil {
		diags.AddError("Failed to Read Team Member", fmt.Sprintf("Could not look up %q: %s", email, err))
		return nil, diags
	}
	if member == nil {
		return nil, diags
	}
	return memberState(ctx, displayEmail, email, member)
}

func (r *teamMemberResource) readMemberByEmail(
	ctx context.Context, displayEmail, email string,
) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	member, err := r.client.GetTeamMemberByEmail(ctx, email)
	if err != nil {
		diags.AddError("Failed to Read Team Member", fmt.Sprintf("Could not look up %q: %s", email, err))
		return nil, diags
	}
	if member == nil {
		return nil, diags
	}
	return memberState(ctx, displayEmail, email, member)
}

func (r *teamMemberResource) readInviteByID(
	ctx context.Context, displayEmail, inviteID string,
) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	invite, err := r.client.GetTeamInvite(ctx, inviteID)
	if err != nil {
		diags.AddError("Failed to Read Team Member", fmt.Sprintf("Could not look up pending invite for %q: %s", displayEmail, err))
		return nil, diags
	}
	if invite == nil {
		return nil, diags
	}
	return inviteState(ctx, displayEmail, invite)
}

func (r *teamMemberResource) readInviteByEmail(
	ctx context.Context, displayEmail, email string,
) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	invite, err := r.client.GetTeamInviteByEmail(ctx, email)
	if err != nil {
		diags.AddError("Failed to Read Team Member", fmt.Sprintf("Could not look up pending invite for %q: %s", email, err))
		return nil, diags
	}
	if invite == nil {
		return nil, diags
	}
	return inviteState(ctx, displayEmail, invite)
}

func memberState(
	ctx context.Context, displayEmail, email string, member *model.MemberByEmailResponse,
) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	tagsList, tagsDiags := types.ListValueFrom(ctx, types.StringType, member.Tags)
	diags.Append(tagsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	return &teamMemberResourceModel{
		ID:       types.StringValue(email),
		Email:    types.StringValue(displayEmail),
		Role:     types.StringValue(member.Role),
		Tags:     tagsList,
		Status:   types.StringValue("member"),
		InviteID: types.StringNull(),
		UserID:   types.StringValue(member.ID),
	}, diags
}

func inviteState(
	ctx context.Context, displayEmail string, invite *model.InviteByEmailResponse,
) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	tagsList, tagsDiags := types.ListValueFrom(ctx, types.StringType, invite.Tags)
	diags.Append(tagsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	return &teamMemberResourceModel{
		ID:       types.StringValue(normalizeEmail(displayEmail)),
		Email:    types.StringValue(displayEmail),
		Role:     types.StringValue(invite.Role),
		Tags:     tagsList,
		Status:   types.StringValue("pending"),
		InviteID: types.StringValue(strconv.Itoa(invite.ID)),
		UserID:   types.StringNull(),
	}, diags
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
