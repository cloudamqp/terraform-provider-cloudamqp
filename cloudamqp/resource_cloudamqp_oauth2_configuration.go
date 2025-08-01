package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &oauth2ConfigurationResource{}
	_ resource.ResourceWithConfigure   = &oauth2ConfigurationResource{}
	_ resource.ResourceWithImportState = &oauth2ConfigurationResource{}
)

type oauth2ConfigurationResource struct {
	client *api.API
}

func NewOAuth2ConfigurationResource() resource.Resource {
	return &oauth2ConfigurationResource{}
}

type oauth2ConfigurationResourceModel struct {
	ID                      types.String `tfsdk:"id"`
	InstanceID              types.Int64  `tfsdk:"instance_id"`
	ResourceServerId        types.String `tfsdk:"resource_server_id"`
	Issuer                  types.String `tfsdk:"issuer"`
	PreferredUsernameClaims types.List   `tfsdk:"preferred_username_claims"`
	AdditionalScopesKeys    types.List   `tfsdk:"additional_scopes_keys"`
	ScopePrefix             types.String `tfsdk:"scope_prefix"`
	ScopeAliases            types.Map    `tfsdk:"scope_aliases"`
	VerifyAud               types.Bool   `tfsdk:"verify_aud"`
	OauthClientId           types.String `tfsdk:"oauth_client_id"`
	OauthScopes             types.List   `tfsdk:"oauth_scopes"`
	Configured              types.Bool   `tfsdk:"configured"`
	Sleep                   types.Int64  `tfsdk:"sleep"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
}

func (r *oauth2ConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(plan.InstanceID.ValueInt64())
	sleep, timeout := extractSleepAndTimeout(&plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	request := model.OAuth2ConfigRequest{}
	populateOAuth2ConfigRequestModel(ctx, &plan, &request)

	data, err := r.client.CreateOAuth2Configuration(timeoutCtx, instanceID, sleep, request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OAuth2 configuration", err.Error())
		return
	}

	err = r.client.PollForOauth2Configured(timeoutCtx, instanceID, *data.ConfigurationId, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for OAuth2 configuration", err.Error())
		return
	}

	data, err = r.client.ReadOAuth2Configuration(timeoutCtx, instanceID, sleep, data.ConfigurationId)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OAuth2 configuration", err.Error())
		return
	}

	populateOAuth2ConfigurationStateModel(ctx, &plan, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oauth2ConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(state.InstanceID.ValueInt64())
	sleep, timeout := extractSleepAndTimeout(&state)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := r.client.DeleteOAuth2Configuration(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting OAuth2 configuration", err.Error())
		return
	}

	settingID := state.ID.ValueString()
	err = r.client.PollForConfigured(timeoutCtx, instanceID, settingID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for deleted OAuth2 configuration", err.Error())
		return
	}
}

func (r *oauth2ConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(state.InstanceID.ValueInt64())
	settingID := state.ID.ValueString()
	sleep, timeout := extractSleepAndTimeout(&state)

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := r.client.ReadOAuth2Configuration(timeoutCtx, instanceID, sleep, &settingID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OAuth2 configuration", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read OAuth2 configuration data: %v", data))

	populateOAuth2ConfigurationStateModel(ctx, &state, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *oauth2ConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := int(plan.InstanceID.ValueInt64())

	preferredUsernameClaims := make([]string, 0)
	plan.PreferredUsernameClaims.ElementsAs(ctx, &preferredUsernameClaims, false)

	scopeAliases := make(map[string]string)
	plan.ScopeAliases.ElementsAs(ctx, &scopeAliases, false)

	oauthScopes := make([]string, 0)
	plan.OauthScopes.ElementsAs(ctx, &oauthScopes, false)

	additionalScopesKeys := make([]string, 0)
	plan.AdditionalScopesKeys.ElementsAs(ctx, &additionalScopesKeys, false)

	params := model.OAuth2ConfigRequest{
		ResourceServerId:        plan.ResourceServerId.ValueString(),
		Issuer:                  plan.Issuer.ValueString(),
		PreferredUsernameClaims: preferredUsernameClaims,
		ScopeAliases:            scopeAliases,
		VerifyAud:               utils.Pointer(plan.VerifyAud.ValueBool()),
		OauthClientId:           plan.OauthClientId.ValueString(),
		OauthScopes:             oauthScopes,
		AdditionalScopesKeys:    additionalScopesKeys,
		ScopePrefix:             plan.ScopePrefix.ValueString(),
	}

	sleep, timeout := extractSleepAndTimeout(&plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := r.client.UpdateOAuth2Configuration(timeoutCtx, instanceID, sleep, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating OAuth2 configuration", err.Error())
		return
	}

	settingID := plan.ID.ValueString()
	err = r.client.PollForConfigured(timeoutCtx, instanceID, settingID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for OAuth2 configuration", err.Error())
		return
	}

	data, err := r.client.ReadOAuth2Configuration(timeoutCtx, instanceID, sleep, &settingID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OAuth2 configuration", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read OAuth2 configuration data: %v", data))

	populateOAuth2ConfigurationStateModel(ctx, &plan, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oauth2ConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var state oauth2ConfigurationResourceModel

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the import ID as a combination of resource ID and instance ID
	parts := strings.Split(req.ID, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: {resource_id},{instance_id}")
		return
	}
	state.ID = types.StringValue(parts[0])
	instanceID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance ID", fmt.Sprintf("Expected numeric instance ID, got: %s", parts[1]))
		return
	}
	state.InstanceID = types.Int64Value(instanceID)
	state.Sleep = types.Int64Value(60)
	state.Timeout = types.Int64Value(3600)

	timeoutCtx, cancel := context.WithTimeout(ctx, 3600*time.Second)
	defer cancel()

	data, err := r.client.ReadOAuth2Configuration(timeoutCtx, int(instanceID), 60*time.Second, state.ID.ValueStringPointer())
	if err != nil {
		resp.Diagnostics.AddError("Error reading OAuth2 configuration", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read OAuth2 configuration data: %v", data))

	populateOAuth2ConfigurationStateModel(ctx, &state, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *oauth2ConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *oauth2ConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_oauth2_configuration"
}

func (r *oauth2ConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Description: "Resource ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"resource_server_id": schema.StringAttribute{
				Required:    true,
				Description: "Resource server identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"issuer": schema.StringAttribute{
				Required:    true,
				Description: "Issuer",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"preferred_username_claims": schema.ListAttribute{
				Optional:    true,
				Description: "Preferred username claims",
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_scopes_keys": schema.ListAttribute{
				Optional:    true,
				Description: "Additional scopes keys",
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"scope_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Scope prefix",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope_aliases": schema.MapAttribute{
				Optional:    true,
				Description: "Scope aliases",
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"verify_aud": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Verify aud",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"oauth_client_id": schema.StringAttribute{
				Optional:    true,
				Description: "Oauth client id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oauth_scopes": schema.ListAttribute{
				Optional:    true,
				Description: "Oauth scopes",
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"configured": schema.BoolAttribute{
				Optional:    false,
				Computed:    true,
				Description: "Configured",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(60),
				Computed:    true,
				Description: "Configurable sleep time in seconds between retries for OAuth2 configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Default:     int64default.StaticInt64(3600),
				Computed:    true,
				Description: "Configurable timeout time in seconds for OAuth2 configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func populateOAuth2ConfigurationStateModel(ctx context.Context, state *oauth2ConfigurationResourceModel, data *model.OAuth2ConfigResponse) {
	state.ID = types.StringValue(*data.ConfigurationId)
	state.ResourceServerId = types.StringValue(*data.ResourceServerId)
	state.Issuer = types.StringValue(*data.Issuer)

	state.PreferredUsernameClaims, _ = types.ListValueFrom(ctx, types.StringType, data.PreferredUsernameClaims)
	state.ScopeAliases, _ = types.MapValueFrom(ctx, types.StringType, data.ScopeAliases)

	if data.ScopePrefix != nil {
		state.ScopePrefix = types.StringValue(*data.ScopePrefix)
	}

	state.AdditionalScopesKeys, _ = types.ListValueFrom(ctx, types.StringType, data.AdditionalScopesKeys)
	state.OauthScopes, _ = types.ListValueFrom(ctx, types.StringType, data.OauthScopes)

	state.VerifyAud = types.BoolValue(*data.VerifyAud)
	state.Configured = types.BoolValue(*data.Configured)
}

func populateOAuth2ConfigRequestModel(ctx context.Context, plan *oauth2ConfigurationResourceModel, data *model.OAuth2ConfigRequest) {
	plan.PreferredUsernameClaims.ElementsAs(ctx, &data.PreferredUsernameClaims, false)
	plan.ScopeAliases.ElementsAs(ctx, &data.ScopeAliases, false)

	data.ResourceServerId = plan.ResourceServerId.ValueString()
	data.Issuer = plan.Issuer.ValueString()
	data.VerifyAud = utils.Pointer(plan.VerifyAud.ValueBool())
	data.OauthClientId = plan.OauthClientId.ValueString()
	plan.OauthScopes.ElementsAs(ctx, &data.OauthScopes, false)
	plan.AdditionalScopesKeys.ElementsAs(ctx, &data.AdditionalScopesKeys, false)
	data.ScopePrefix = plan.ScopePrefix.ValueString()
}

func extractSleepAndTimeout(plan *oauth2ConfigurationResourceModel) (time.Duration, time.Duration) {
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second

	if sleep == 0 {
		sleep = time.Duration(60) * time.Second
	}

	if timeout == 0 {
		timeout = time.Duration(3600) * time.Second
	}

	return sleep, timeout
}
