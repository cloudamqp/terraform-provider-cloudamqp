package cloudamqp

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// teamMemberTimeout lets requests retry the API's rate limit, which needs a deadline.
const teamMemberTimeout = 5 * time.Minute

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
	Tags     types.Set    `tfsdk:"tags"`
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
				Required: true,
				Description: "Role assigned to the member. One of: " +
					strings.Join(assignableTeamRoleNames, ", ") + ".",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						pendingInviteRequiresReplaceString,
						pendingInviteReplaceDescription,
						pendingInviteReplaceDescription,
					),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(assignableTeamRoleNames...),
				},
			},
			"tags": schema.SetAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tag names granting access to tagged instances. A tag that doesn't exist yet " +
					"is created automatically. Only meaningful for the \"member\"/\"monitor\" roles, " +
					"which require at least one tag. At most 300 tags, each shorter than 500 characters.",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(300),
					setvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, 499)),
				},
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
					setplanmodifier.RequiresReplaceIf(
						pendingInviteRequiresReplaceSet,
						pendingInviteReplaceDescription,
						pendingInviteReplaceDescription,
					),
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
	role := config.Role.ValueString()

	// Tags referencing another resource's output can't be checked yet either.
	if config.Tags.IsUnknown() {
		return
	}

	hasTags := !config.Tags.IsNull() && len(config.Tags.Elements()) > 0
	tagged := slices.Contains(taggedTeamRoleNames, role)

	if tagged && !hasTags {
		resp.Diagnostics.AddAttributeError(
			path.Root("tags"),
			"Missing Required Tags",
			fmt.Sprintf("The %q role requires at least one tag to be specified.", role),
		)
	}

	// The backend clears all tags when switching to a role other than member/monitor.
	if !tagged && hasTags {
		resp.Diagnostics.AddAttributeError(
			path.Root("tags"),
			"Tags Not Allowed",
			fmt.Sprintf("The %q role can't have tags, only the %s roles can.",
				role, strings.Join(taggedTeamRoleNames, "/")),
		)
	}
}

const pendingInviteReplaceDescription = "A pending invite can't be updated. Changing role or tags " +
	"while the invite is pending revokes it and sends a new invite."

func pendingInviteRequiresReplaceString(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
	resp.RequiresReplace, resp.Diagnostics = pendingInviteRequiresReplace(ctx, req.State, req.Path)
}

func pendingInviteRequiresReplaceSet(ctx context.Context, req planmodifier.SetRequest, resp *setplanmodifier.RequiresReplaceIfFuncResponse) {
	if tagsEqualIgnoringCase(ctx, req.StateValue, req.PlanValue) {
		return
	}
	resp.RequiresReplace, resp.Diagnostics = pendingInviteRequiresReplace(ctx, req.State, req.Path)
}

func pendingInviteRequiresReplace(ctx context.Context, state tfsdk.State, attr path.Path) (bool, diag.Diagnostics) {
	var status, email types.String
	diags := state.GetAttribute(ctx, path.Root("status"), &status)
	diags.Append(state.GetAttribute(ctx, path.Root("email"), &email)...)
	if diags.HasError() || status.ValueString() != "pending" {
		return false, diags
	}

	diags.AddAttributeWarning(
		attr,
		"Pending Invite Will Be Replaced",
		fmt.Sprintf("The invite for %q is still pending and can't be updated. Applying this change "+
			"revokes the invite and sends a new one.", email.ValueString()),
	)
	return true, diags
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
	ctx, cancel := context.WithTimeout(ctx, teamMemberTimeout)
	defer cancel()

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

	invite, err := r.client.CreateTeamInvite(ctx, model.InviteRequest{
		Email: email,
		Role:  plan.Role.ValueString(),
		Tags:  tags,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Team Member",
			fmt.Sprintf("Could not invite %q: %s", email, err),
		)
		return
	}

	state, diags := inviteState(ctx, displayEmail, invite)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(keepConfiguredTagCase(ctx, state, tags)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *teamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, cancel := context.WithTimeout(ctx, teamMemberTimeout)
	defer cancel()

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

	var known []string
	if !state.Tags.IsNull() && !state.Tags.IsUnknown() {
		resp.Diagnostics.Append(state.Tags.ElementsAs(ctx, &known, false)...)
	}
	_, diags = tagsWithKnownCase(ctx, newState, known)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *teamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, cancel := context.WithTimeout(ctx, teamMemberTimeout)
	defer cancel()

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

	// Other role/tag changes on a pending invite replace it, see pendingInviteRequiresReplace.
	if state.Status.ValueString() == "pending" {
		state.Tags = plan.Tags
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		return
	}

	err := r.client.UpdateTeamMember(ctx, state.UserID.ValueString(), model.MemberUpdateRequest{
		Role: plan.Role.ValueString(),
		Tags: tags,
	})
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
	resp.Diagnostics.Append(keepConfiguredTagCase(ctx, newState, tags)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *teamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, cancel := context.WithTimeout(ctx, teamMemberTimeout)
	defer cancel()

	var state teamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	email := state.Email.ValueString()

	var err error
	if state.Status.ValueString() == "pending" {
		var revoked bool
		revoked, err = r.client.DeleteTeamInvite(ctx, state.InviteID.ValueString())
		if err == nil && !revoked {
			resp.Diagnostics.Append(r.checkInviteNotAccepted(ctx, email)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	} else {
		err = r.client.RemoveTeamMember(ctx, normalizeEmail(email))
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Team Member",
			fmt.Sprintf("Could not remove %q: %s", email, err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

// checkInviteNotAccepted gives an error if a pending invite that's already gone was accepted
// since the last refresh, e.g. between plan and apply. Removing the resource from state then
// would make the next invite fail with "already in this team" on every apply. Keeping the state
// lets the next plan pick up the member and update it in place.
func (r *teamMemberResource) checkInviteNotAccepted(ctx context.Context, email string) diag.Diagnostics {
	var diags diag.Diagnostics
	member, err := r.client.GetTeamMemberByEmail(ctx, normalizeEmail(email))
	if err != nil {
		diags.AddError("Failed to Delete Team Member", fmt.Sprintf("Could not look up %q: %s", email, err))
		return diags
	}
	if member != nil {
		diags.AddError(
			"Invite Already Accepted",
			fmt.Sprintf("The invite for %q was accepted after the plan was made, so the planned "+
				"change no longer applies. Run plan again to get a plan for the member.", email),
		)
	}
	return diags
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
	tags, tagsDiags := types.SetValueFrom(ctx, types.StringType, member.Tags)
	diags.Append(tagsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	return &teamMemberResourceModel{
		ID:       types.StringValue(email),
		Email:    types.StringValue(displayEmail),
		Role:     types.StringValue(member.Role),
		Tags:     tags,
		Status:   types.StringValue("member"),
		InviteID: types.StringNull(),
		UserID:   types.StringValue(member.ID),
	}, diags
}

func inviteState(
	ctx context.Context, displayEmail string, invite *model.InviteResponse,
) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	tags, tagsDiags := types.SetValueFrom(ctx, types.StringType, invite.Tags)
	diags.Append(tagsDiags...)
	if diags.HasError() {
		return nil, diags
	}
	return &teamMemberResourceModel{
		ID:       types.StringValue(normalizeEmail(displayEmail)),
		Email:    types.StringValue(displayEmail),
		Role:     types.StringValue(invite.Role),
		Tags:     tags,
		Status:   types.StringValue("pending"),
		InviteID: types.StringValue(strconv.Itoa(invite.ID)),
		UserID:   types.StringNull(),
	}, diags
}

// tagsWithKnownCase keeps the spelling from known for tags the backend matched to a team tag
// that only differs in letter case. Returns the replaced tags as known => assigned.
func tagsWithKnownCase(
	ctx context.Context, state *teamMemberResourceModel, known []string,
) (map[string]string, diag.Diagnostics) {
	var assigned []string
	diags := state.Tags.ElementsAs(ctx, &assigned, false)
	if diags.HasError() {
		return nil, diags
	}

	replaced := map[string]string{}
	tags := make([]string, 0, len(assigned))
	for _, tag := range assigned {
		i := slices.IndexFunc(known, func(k string) bool { return strings.EqualFold(k, tag) })
		if i < 0 || slices.Contains(known, tag) {
			tags = append(tags, tag)
			continue
		}
		tags = append(tags, known[i])
		replaced[known[i]] = tag
	}

	state.Tags, diags = types.SetValueFrom(ctx, types.StringType, tags)
	return replaced, diags
}

// keepConfiguredTagCase is tagsWithKnownCase for the configured tags, with a warning.
func keepConfiguredTagCase(ctx context.Context, state *teamMemberResourceModel, configured []string) diag.Diagnostics {
	replaced, diags := tagsWithKnownCase(ctx, state, configured)
	if diags.HasError() || len(replaced) == 0 {
		return diags
	}

	names := make([]string, 0, len(replaced))
	for tag, teamTag := range replaced {
		names = append(names, fmt.Sprintf("%q was assigned the team's tag %q", tag, teamTag))
	}
	slices.Sort(names)
	diags.AddAttributeWarning(
		path.Root("tags"),
		"Tag Assigned With Different Letter Case",
		fmt.Sprintf("The team already has tags that only differ in letter case, and tags are case-sensitive "+
			"for access, so those were assigned instead: %s. Consider using the team's spelling.",
			strings.Join(names, ", ")),
	)
	return diags
}

// tagsEqualIgnoringCase reports whether two tag sets only differ in letter case.
func tagsEqualIgnoringCase(ctx context.Context, a, b types.Set) bool {
	if a.IsNull() || a.IsUnknown() || b.IsNull() || b.IsUnknown() {
		return false
	}
	var aTags, bTags []string
	if a.ElementsAs(ctx, &aTags, false).HasError() || b.ElementsAs(ctx, &bTags, false).HasError() {
		return false
	}
	lower := func(tags []string) []string {
		for i, tag := range tags {
			tags[i] = strings.ToLower(tag)
		}
		slices.Sort(tags)
		return slices.Compact(tags)
	}
	return slices.Equal(lower(aTags), lower(bTags))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
