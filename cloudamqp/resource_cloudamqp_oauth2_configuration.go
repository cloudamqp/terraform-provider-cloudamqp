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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	AdditionalScopesKey     types.List   `tfsdk:"additional_scopes_key"`
	ScopePrefix             types.String `tfsdk:"scope_prefix"`
	ScopeAliases            types.Map    `tfsdk:"scope_aliases"`
	VerifyAud               types.Bool   `tfsdk:"verify_aud"`
	OauthClientId           types.String `tfsdk:"oauth_client_id"`
	OauthScopes             types.List   `tfsdk:"oauth_scopes"`
	Audience                types.String `tfsdk:"audience"`
	DisableBasicAuth        types.Bool   `tfsdk:"disable_basic_auth"`
	JwksCacheTtl            types.Int64  `tfsdk:"jwks_cache_ttl"`
	MgmtScopes              types.String `tfsdk:"mgmt_scopes"`
	Sleep                   types.Int64  `tfsdk:"sleep"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
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
			"additional_scopes_key": schema.ListAttribute{
				Optional:    true,
				Description: "Additional scopes key",
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
			"audience": schema.StringAttribute{
				Optional:    true,
				Description: "Audience",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disable_basic_auth": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Disable basic auth for OAuth2 configuration",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"jwks_cache_ttl": schema.Int64Attribute{
				Optional:    true,
				Description: "LavinMQ only. How long JWKS signing keys are cached before being re-fetched, in seconds",
				Validators: []validator.Int64{
					int64validator.Between(0, 86400),
				},
			},
			"mgmt_scopes": schema.StringAttribute{
				Optional:    true,
				Description: "LavinMQ only. Space-separated scopes requested during management UI SSO login",
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

func (r *oauth2ConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("ImportState: ID=%s", req.ID))
	instanceID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected numeric instance_id, got: %q", req.ID))
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), req.ID)
	resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)
}

func (r *oauth2ConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep, timeout := extractSleepAndTimeout(&plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var config oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backend, diags := r.readOAuth2Backend(timeoutCtx, instanceID, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := model.OAuth2ConfigRequest{}
	populateOAuth2ConfigRequestModel(ctx, &plan, &request)
	tflog.Info(ctx, fmt.Sprintf("CreateOAuth2Configuration with request: %+v", request))

	job, err := r.client.CreateOAuth2Configuration(timeoutCtx, instanceID, backend, sleep, request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating OAuth2 configuration", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *job.ID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for OAuth2 configuration", err.Error())
		return
	}

	data, err := r.client.ReadOAuth2Configuration(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OAuth2 configuration", err.Error())
		return
	}

	populateOAuth2ConfigurationStateModel(ctx, &plan, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oauth2ConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep, timeout := extractSleepAndTimeout(&state)

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := r.client.ReadOAuth2Configuration(timeoutCtx, instanceID, sleep)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			tflog.Info(ctx, "OAuth2 configuration not found, removing resource")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading OAuth2 configuration", err.Error())
		return
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("oauth2 configuration not found, resource will be recreated: %s", state.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read OAuth2 configuration data: %+v", data))

	populateOAuth2ConfigurationStateModel(ctx, &state, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *oauth2ConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep, timeout := extractSleepAndTimeout(&plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	backend, diags := r.readOAuth2Backend(timeoutCtx, instanceID, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := model.OAuth2ConfigRequest{}
	populateOAuth2ConfigRequestModel(ctx, &plan, &params)
	tflog.Info(ctx, fmt.Sprintf("Updating OAuth2 configuration with params: %+v", params))
	job, err := r.client.UpdateOAuth2Configuration(timeoutCtx, instanceID, backend, sleep, params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating OAuth2 configuration", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *job.ID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for OAuth2 configuration", err.Error())
		return
	}

	data, err := r.client.ReadOAuth2Configuration(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error reading OAuth2 configuration", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read OAuth2 configuration data: %+v", data))

	populateOAuth2ConfigurationStateModel(ctx, &plan, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *oauth2ConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state oauth2ConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep, timeout := extractSleepAndTimeout(&state)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	job, err := r.client.DeleteOAuth2Configuration(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting OAuth2 configuration", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *job.ID, sleep)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for deleted OAuth2 configuration", err.Error())
		return
	}
}

// Private methods

// oauth2Attribute is a resource attribute that is only supported by one broker
type oauth2Attribute struct {
	name  string
	isSet bool
}

// readOAuth2Backend looks up which broker the instance runs and checks that the configuration
// only sets attributes that broker supports.
func (r *oauth2ConfigurationResource) readOAuth2Backend(ctx context.Context, instanceID int64, config *oauth2ConfigurationResourceModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	backend, err := r.client.ReadInstanceBackend(ctx, instanceID)
	if err != nil {
		diags.AddError("Error reading instance backend", err.Error())
		return "", diags
	}
	tflog.Debug(ctx, fmt.Sprintf("OAuth2 configuration for instance %d with backend %q", instanceID, backend))

	if backend == "lavinmq" {
		diags.Append(validateOAuth2BrokerAttributes(instanceID, "LavinMQ", []oauth2Attribute{
			{"scope_aliases", !config.ScopeAliases.IsNull()},
			{"oauth_scopes", !config.OauthScopes.IsNull()},
			{"disable_basic_auth", !config.DisableBasicAuth.IsNull()},
		})...)
	} else {
		diags.Append(validateOAuth2BrokerAttributes(instanceID, "RabbitMQ", []oauth2Attribute{
			{"jwks_cache_ttl", !config.JwksCacheTtl.IsNull()},
			{"mgmt_scopes", !config.MgmtScopes.IsNull()},
		})...)
	}
	return backend, diags
}

func validateOAuth2BrokerAttributes(instanceID int64, broker string, attributes []oauth2Attribute) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, attribute := range attributes {
		if attribute.isSet {
			diags.AddAttributeError(
				path.Root(attribute.name),
				"Unsupported OAuth2 configuration attribute",
				fmt.Sprintf("Attribute %q is not supported for instance %d, which runs %s. Remove it from the configuration.",
					attribute.name, instanceID, broker),
			)
		}
	}
	return diags
}

func populateOAuth2ConfigurationStateModel(ctx context.Context, state *oauth2ConfigurationResourceModel, data *model.OAuth2ConfigResponse) {
	state.ID = types.StringValue(data.ConfigurationId)
	state.ResourceServerId = types.StringValue(data.ResourceServerId)
	state.Issuer = types.StringValue(data.Issuer)

	state.PreferredUsernameClaims, _ = types.ListValueFrom(ctx, types.StringType, data.PreferredUsernameClaims)
	state.ScopeAliases, _ = types.MapValueFrom(ctx, types.StringType, data.ScopeAliases)

	if data.ScopePrefix != nil {
		state.ScopePrefix = types.StringValue(*data.ScopePrefix)
	}

	state.AdditionalScopesKey, _ = types.ListValueFrom(ctx, types.StringType, data.AdditionalScopesKeys)
	state.OauthScopes, _ = types.ListValueFrom(ctx, types.StringType, data.OauthScopes)
	if data.ClientId != nil {
		state.OauthClientId = types.StringValue(*data.ClientId)
	}

	if data.Audience != nil {
		state.Audience = types.StringValue(*data.Audience)
	}

	if data.VerifyAud != nil {
		state.VerifyAud = types.BoolValue(*data.VerifyAud)
	}
	// LavinMQ has no oauth_disable_basic_auth, default to false to keep the computed attribute stable
	state.DisableBasicAuth = types.BoolValue(data.DisableBasicAuth != nil && *data.DisableBasicAuth)

	state.JwksCacheTtl = types.Int64PointerValue(data.JwksCacheTtl)
	state.MgmtScopes = types.StringPointerValue(data.MgmtScopes)
}

// populateOAuth2ConfigRequestModel only sets known values, the API client applies broker defaults
func populateOAuth2ConfigRequestModel(ctx context.Context, plan *oauth2ConfigurationResourceModel, data *model.OAuth2ConfigRequest) {
	plan.PreferredUsernameClaims.ElementsAs(ctx, &data.PreferredUsernameClaims, false)
	plan.AdditionalScopesKey.ElementsAs(ctx, &data.AdditionalScopesKeys, false)
	plan.ScopeAliases.ElementsAs(ctx, &data.ScopeAliases, false)
	plan.OauthScopes.ElementsAs(ctx, &data.OauthScopes, false)

	data.ResourceServerId = plan.ResourceServerId.ValueString()
	data.Issuer = plan.Issuer.ValueString()
	data.VerifyAud = knownBoolPointer(plan.VerifyAud)
	data.ClientId = plan.OauthClientId.ValueString()
	data.ScopePrefix = plan.ScopePrefix.ValueString()
	data.Audience = plan.Audience.ValueString()
	data.DisableBasicAuth = knownBoolPointer(plan.DisableBasicAuth)
	data.JwksCacheTtl = plan.JwksCacheTtl.ValueInt64Pointer()
	data.MgmtScopes = plan.MgmtScopes.ValueString()
}

func knownBoolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	return utils.Pointer(value.ValueBool())
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
