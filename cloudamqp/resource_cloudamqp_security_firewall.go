package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &firewallResource{}
	_ resource.ResourceWithConfigure   = &firewallResource{}
	_ resource.ResourceWithImportState = &firewallResource{}
)

type firewallResource struct {
	client *api.API
}

func NewFirewallResource() resource.Resource {
	return &firewallResource{}
}

type firewallResourceModel struct {
	ID         types.String        `tfsdk:"id"`
	InstanceID types.Int64         `tfsdk:"instance_id"`
	Rules      []ruleResourceModel `tfsdk:"rules"`
	Sleep      types.Int64         `tfsdk:"sleep"`
	Timeout    types.Int64         `tfsdk:"timeout"`
}

type ruleResourceModel struct {
	IP          types.String `tfsdk:"ip"`
	Services    types.List   `tfsdk:"services"`
	Ports       types.List   `tfsdk:"ports"`
	Description types.String `tfsdk:"description"`
}

func (r *firewallResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_security_firewall"
}

func (r *firewallResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure firewall rules for CloudAMQP instances",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID (same as instance_id)",
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
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(10),
				Description: "Configurable sleep time in seconds between retries for firewall configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(600),
				Description: "Configurable timeout time in seconds for firewall configuration",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"rules": schema.SetNestedBlock{
				Description: "Firewall rules for the instance",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"ip": schema.StringAttribute{
							Required:    true,
							Description: "CIDR address: IP address with CIDR notation (e.g. 10.56.72.0/24)",
							Validators: []validator.String{
								validators.CidrValidator{},
							},
						},
						"services": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,

							Description: "Pre-defined services",
							Validators: []validator.List{
								listvalidator.ValueStringsAre(
									stringvalidator.OneOf("AMQP", "AMQPS", "HTTPS", "MQTT", "MQTTS", "STOMP", "STOMPS", "STREAM", "STREAM_SSL"),
								),
							},
						},
						"ports": schema.ListAttribute{
							ElementType: types.Int64Type,
							Optional:    true,
							Computed:    true,

							Description: "Custom ports between 0 - 65554",
							Validators: []validator.List{
								listvalidator.ValueInt64sAre(
									int64validator.Between(0, 65554),
									&portNotServiceValidator{},
								),
							},
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "Naming description e.g. 'Default'",
						},
					},
				},
			},
		},
	}
}

func (r *firewallResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *firewallResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("ImportState: ID=%s", req.ID))

	instanceID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Could not parse instance_id: %s", err))
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), req.ID)
	resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)
	// Default values
	resp.State.SetAttribute(ctx, path.Root("sleep"), int64(10))
	resp.State.SetAttribute(ctx, path.Root("timeout"), int64(600))
}

func (r *firewallResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan firewallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := make([]model.FirewallRuleRequest, len(plan.Rules))

	for i, rule := range plan.Rules {
		params[i].Ip = rule.IP.ValueString()

		if !rule.Description.IsNull() && !rule.Description.IsUnknown() {
			params[i].Description = rule.Description.ValueString()
		}

		if !rule.Services.IsNull() && !rule.Services.IsUnknown() {
			var services []string
			rule.Services.ElementsAs(ctx, &services, false)
			params[i].Services = services
		} else {
			params[i].Services = []string{}
		}

		if !rule.Ports.IsNull() && !rule.Ports.IsUnknown() {
			var ports []int64
			rule.Ports.ElementsAs(ctx, &ports, false)
			params[i].Ports = ports
		} else {
			params[i].Ports = []int64{}
		}
	}

	err := r.client.CreateFirewallSettings(timeoutCtx, instanceID, params, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Security Firewall",
			fmt.Sprintf("Could not create security firewall for instance %d: %s", instanceID, err),
		)
		return
	}

	err = r.client.PollForFirewallConfigured(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Waiting for Security Firewall Configuration",
			fmt.Sprintf("Error while waiting for security firewall to be configured for instance %d: %s", instanceID, err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(instanceID, 10))
	normalizeRules(plan.Rules)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *firewallResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state firewallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep := time.Duration(state.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(state.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := r.client.ReadFirewallSettings(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Security Firewall",
			fmt.Sprintf("Could not read security firewall for instance %d: %s", instanceID, err),
		)
		return
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("firewall settings not found, resource will be recreated: %d", instanceID))
		resp.State.RemoveResource(ctx)
		return
	}

	rulesModel := []ruleResourceModel{}
	for _, rule := range *data {
		services := rule.Services
		if services == nil {
			services = []string{}
		}
		ports := rule.Ports
		if ports == nil {
			ports = []int64{}
		}
		ruleModel := ruleResourceModel{}
		ruleModel.IP = types.StringValue(rule.Ip)
		ruleModel.Description = types.StringPointerValue(rule.Description)
		ruleModel.Services, _ = types.ListValueFrom(ctx, types.StringType, services)
		ruleModel.Ports, _ = types.ListValueFrom(ctx, types.Int64Type, ports)
		rulesModel = append(rulesModel, ruleModel)
	}

	state.Rules = rulesModel
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *firewallResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state firewallResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating firewall for instance with plan: %+v", plan))

	instanceID := plan.InstanceID.ValueInt64()
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	params := make([]model.FirewallRuleRequest, len(plan.Rules))

	for i, rule := range plan.Rules {
		params[i].Ip = rule.IP.ValueString()

		if !rule.Description.IsNull() && !rule.Description.IsUnknown() {
			params[i].Description = rule.Description.ValueString()
		}

		if !rule.Services.IsNull() && !rule.Services.IsUnknown() {
			var services []string
			rule.Services.ElementsAs(ctx, &services, false)
			params[i].Services = services
		} else {
			params[i].Services = []string{}
		}

		if !rule.Ports.IsNull() && !rule.Ports.IsUnknown() {
			var ports []int64
			rule.Ports.ElementsAs(ctx, &ports, false)
			params[i].Ports = ports
		} else {
			params[i].Ports = []int64{}
		}
	}

	// Check if rules changed (order-independent set comparison by hashing each rule).
	// If only sleep/timeout changed, skip API call and just persist new state.
	rulesChanged := len(plan.Rules) != len(state.Rules)
	if !rulesChanged {
		ruleKey := func(r ruleResourceModel) string {
			var services []string
			if !r.Services.IsNull() && !r.Services.IsUnknown() {
				r.Services.ElementsAs(ctx, &services, false)
			}
			var ports []int64
			if !r.Ports.IsNull() && !r.Ports.IsUnknown() {
				r.Ports.ElementsAs(ctx, &ports, false)
			}
			return fmt.Sprintf("%s|%v|%v|%s", r.IP.ValueString(), services, ports, r.Description.ValueString())
		}
		stateKeys := make(map[string]struct{}, len(state.Rules))
		for _, r := range state.Rules {
			stateKeys[ruleKey(r)] = struct{}{}
		}
		for _, r := range plan.Rules {
			if _, ok := stateKeys[ruleKey(r)]; !ok {
				rulesChanged = true
				break
			}
		}
	}

	if !rulesChanged {
		tflog.Info(ctx, fmt.Sprintf("only sleep/timeout changed for instance %d, skipping API call", instanceID))
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	// Update firewall settings
	err := r.client.UpdateFirewallSettings(timeoutCtx, instanceID, params, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Security Firewall",
			fmt.Sprintf("Could not update security firewall for instance %d: %s", instanceID, err),
		)
		return
	}

	err = r.client.PollForFirewallConfigured(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Waiting for Security Firewall Configuration",
			fmt.Sprintf("Error while waiting for security firewall to be configured for instance %d: %s", instanceID, err),
		)
		return
	}

	normalizeRules(plan.Rules)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *firewallResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state firewallResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if enableFasterInstanceDestroy {
		tflog.Info(ctx, "delete being skipped and no call to backend")
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	sleep := time.Duration(state.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(state.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Delete firewall settings (send empty rules array)
	err := r.client.DeleteFirewallSettings(timeoutCtx, instanceID, []model.FirewallRuleRequest{}, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Security Firewall",
			fmt.Sprintf("Could not delete security firewall for instance %d: %s", instanceID, err),
		)
		return
	}

	err = r.client.PollForFirewallConfigured(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Waiting for Security Firewall Configuration",
			fmt.Sprintf("Error while waiting for security firewall to be configured for instance %d: %s", instanceID, err),
		)
		return
	}
}

// normalizeRules ensures services and ports are always known (empty list) in state,
// preventing "unknown value after apply" errors when optional attributes are omitted.
func normalizeRules(rules []ruleResourceModel) {
	emptyStrings := types.ListValueMust(types.StringType, []attr.Value{})
	emptyInt64s := types.ListValueMust(types.Int64Type, []attr.Value{})
	for i := range rules {
		if rules[i].Services.IsNull() || rules[i].Services.IsUnknown() {
			rules[i].Services = emptyStrings
		}
		if rules[i].Ports.IsNull() || rules[i].Ports.IsUnknown() {
			rules[i].Ports = emptyInt64s
		}
	}
}

// portNotServiceValidator validates that ports don't include predefined service ports
type portNotServiceValidator struct{}

func (v *portNotServiceValidator) Description(ctx context.Context) string {
	return "Port should not be a predefined service port. Use 'services' attribute instead."
}

func (v *portNotServiceValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *portNotServiceValidator) ValidateInt64(ctx context.Context, req validator.Int64Request, resp *validator.Int64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	port := int(req.ConfigValue.ValueInt64())
	servicePorts := []struct {
		Port    int
		Service string
	}{
		{5672, "AMQP"},
		{5671, "AMQPS"},
		{443, "HTTPS"},
		{1883, "MQTT"},
		{8883, "MQTTS"},
		{61613, "STOMP"},
		{61614, "STOMPS"},
		{5552, "STREAM"},
		{5551, "STREAM_SSL"},
	}

	for _, sp := range servicePorts {
		if sp.Port == port {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Port should use 'services' attribute",
				fmt.Sprintf("Port %d is a predefined service port. Please add '%s' to the 'services' list instead of using the 'ports' list.", port, sp.Service),
			)
			return
		}
	}
}
