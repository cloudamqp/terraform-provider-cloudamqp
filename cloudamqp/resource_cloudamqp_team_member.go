package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/team"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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

var (
	_ resource.Resource                = &teamMemberResource{}
	_ resource.ResourceWithConfigure   = &teamMemberResource{}
	_ resource.ResourceWithImportState = &teamMemberResource{}
)

type teamMemberResource struct {
	client *api.API
}

func NewTeamMemberResource() resource.Resource {
	return &teamMemberResource{}
}

type teamMemberResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Email          types.String `tfsdk:"email"`
	Role           types.String `tfsdk:"role"`
	Tags           types.List   `tfsdk:"tags"`
	Status         types.String `tfsdk:"status"`
	InviteID       types.String `tfsdk:"invite_id"`
	UserID         types.String `tfsdk:"user_id"`
	SamlConfigured types.Bool   `tfsdk:"saml_configured"`
}

func (r *teamMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_team_member"
}

func (r *teamMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a team member invite/membership as a single resource across both states: " +
			"a pending invite before it's accepted, and an accepted member afterwards. \"status\" reports " +
			"which one currently applies, it's computed, not configured. For SAML-configured teams, " +
			"\"role\"/\"tags\" are not managed by Terraform (the IdP controls them on login); configured " +
			"values are ignored with a warning.",
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
				Description: "Tags (created via cloudamqp_team_tag) granting access to tagged instances. " +
					"Only meaningful for the \"member\"/\"monitor\" roles.",
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
			"saml_configured": schema.BoolAttribute{
				Computed:    true,
				Description: "True when the team has SAML configured, in which case role/tags are not managed by Terraform.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
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

	// POST /api/team/invite joins tags with "," and the backend splits on "," - a tag name
	// containing a comma would silently become two different tags.
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

	err := r.client.InviteTeamMember(ctx, model.InviteRequest{
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

	state, diags := r.readState(ctx, displayEmail)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state == nil {
		resp.Diagnostics.AddError(
			"Failed to Create Team Member",
			fmt.Sprintf("Invite for %q succeeded but it could not be found immediately afterwards", email),
		)
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

	newState, diags := r.readState(ctx, state.Email.ValueString())
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

	// Skip pushing role/tags changes for a SAML-configured team - readState below adds the
	// warning explaining why, once, after re-reading the actual backend state.
	if !state.SamlConfigured.ValueBool() {
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
	}

	newState, diags := r.readState(ctx, email)
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

func (r *teamMemberResource) readState(ctx context.Context, displayEmail string) (*teamMemberResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	email := normalizeEmail(displayEmail)

	data, err := r.client.LookupTeamMember(ctx, email)
	if err != nil {
		diags.AddError("Failed to Read Team Member", fmt.Sprintf("Could not look up %q: %s", email, err))
		return nil, diags
	}
	if data == nil {
		return nil, diags
	}

	tagsList, tagsDiags := types.ListValueFrom(ctx, types.StringType, data.Tags)
	diags.Append(tagsDiags...)
	if diags.HasError() {
		return nil, diags
	}

	if data.SamlConfigured {
		diags.AddWarning(
			"Role/Tags Not Managed by Terraform",
			fmt.Sprintf("Team member %q belongs to a SAML-configured team. Role and tags are controlled by "+
				"the IdP's role mapping on login and are not managed by Terraform.", email),
		)
	}

	state := &teamMemberResourceModel{
		ID:             types.StringValue(email),
		Email:          types.StringValue(displayEmail),
		Role:           types.StringValue(data.Role),
		Tags:           tagsList,
		Status:         types.StringValue(data.Status),
		SamlConfigured: types.BoolValue(data.SamlConfigured),
		InviteID:       types.StringNull(),
		UserID:         types.StringNull(),
	}
	if data.Status == "pending" {
		state.InviteID = types.StringValue(memberIDString(data.ID))
	} else {
		state.UserID = types.StringValue(memberIDString(data.ID))
	}

	return state, diags
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func memberIDString(id any) string {
	switch v := id.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return ""
	}
}
