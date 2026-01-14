package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &instanceResource{}
	_ resource.ResourceWithConfigure   = &instanceResource{}
	_ resource.ResourceWithImportState = &instanceResource{}
	_ resource.ResourceWithModifyPlan  = &instanceResource{}
	// _ resource.ResourceWithConfigValidators = &instanceResource{}
	_ resource.Resource
)

type instanceResource struct {
	client *api.API
}

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

type instanceResourceModel struct {
	// Resource identifier
	ID types.String `tfsdk:"id"`
	// Required fields
	Name   types.String `tfsdk:"name"`
	Plan   types.String `tfsdk:"plan"`
	Region types.String `tfsdk:"region"`
	// Optional fields
	Tags            types.List   `tfsdk:"tags"`
	RmqVersion      types.String `tfsdk:"rmq_version"`
	VpcID           types.Int64  `tfsdk:"vpc_id"`
	VpcSubnet       types.String `tfsdk:"vpc_subnet"`
	Nodes           types.Int64  `tfsdk:"nodes"`
	NoDefaultAlarms types.Bool   `tfsdk:"no_default_alarms"`
	PreferredAz     types.List   `tfsdk:"preferred_az"`
	// Computed fields
	Url          types.String `tfsdk:"url"`
	ApiKey       types.String `tfsdk:"apikey"`
	Host         types.String `tfsdk:"host"`
	HostInternal types.String `tfsdk:"host_internal"`
	Vhost        types.String `tfsdk:"vhost"`
	Ready        types.Bool   `tfsdk:"ready"`
	Dedicated    types.Bool   `tfsdk:"dedicated"`
	Backend      types.String `tfsdk:"backend"`
	// Reasource fields
	KeepAssociatedVpc types.Bool  `tfsdk:"keep_associated_vpc"`
	Sleep             types.Int64 `tfsdk:"sleep"`
	Timeout           types.Int64 `tfsdk:"timeout"`
	// Nested Blocks
	CopySettings types.Set `tfsdk:"copy_settings"`
}

type copySettingsModel struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	Settings       types.List   `tfsdk:"settings"`
}

func (r *instanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_instance"
}

func (r *instanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to create and manage a CloudAMQP instance running either RabbitMQ or LavinMQ.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this resource instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the instance",
			},
			"plan": schema.StringAttribute{
				Required:    true,
				Description: "Name of the subscription plan",
				// Validators: []validator.String{
				// 	validators.PlanValidator{Client: r.client},
				// },
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Name of the region you want to create your instance in",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				// Validators: []validator.String{
				// 	validators.RegionValidator{Client: r.client},
				// },
			},
			"tags": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Tag the instances with optional tags",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"rmq_version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "RabbitMQ/LavinMQ version",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The VPC ID. Use this to create your instance in an existing VPC",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"vpc_subnet": schema.StringAttribute{
				Optional: true,
				Description: "Creates a dedicated VPC subnet, shouldn't overlap with other VPC subnet, " +
					"default subnet used 10.56.72.0/24",
				DeprecationMessage: "Support for this attribute will be removed in next major version (2.0)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nodes": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of nodes in cluster (plan must support it)",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"no_default_alarms": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Set to true to not create default alarms",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"preferred_az": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Preferred availability zone for the instance",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "URL of the CloudAMQP instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"apikey": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "API key for the CloudAMQP instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Computed:    true,
				Description: "External hostname for the CloudAMQP instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host_internal": schema.StringAttribute{
				Computed:    true,
				Description: "Internal hostname for the CloudAMQP instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vhost": schema.StringAttribute{
				Computed:    true,
				Description: "The virtual host",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ready": schema.BoolAttribute{
				Computed:    true,
				Description: "Indication when the instance is ready",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dedicated": schema.BoolAttribute{
				Computed:    true,
				Description: "Is the instance hosted on a dedicated server",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"backend": schema.StringAttribute{
				Computed:    true,
				Description: "Information about the backend used, either 'rabbitmq' or 'lavinmq'",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"keep_associated_vpc": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Keep associated VPC when deleting instance, default set to false",
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable sleep time in seconds between retries for instance operations",
				Default:     int64default.StaticInt64(30),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Configurable timeout time in seconds for instance operations",
				Default:     int64default.StaticInt64(3600),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"copy_settings": schema.SetNestedBlock{
				Description: "Copy settings from one instance to a new instance being created",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"subscription_id": schema.StringAttribute{
							Required:    true,
							Description: "Instance identifier of the CloudAMQP instance to copy the settings from",
						},
						"settings": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Array of one or more settings to be copied",
							Validators: []validator.List{
								listvalidator.ValueStringsAre(stringvalidator.OneOf(
									"alarms",
									"config",
									"definitions",
									"firewall",
									"logs",
									"metrics",
									"plugins",
								)),
							},
						},
					},
				},
			},
		},
	}
}

func (r *instanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ModifyPlan, similar to customDiff in SDK to validate plan and region against API.
// Forces replacement when changing between shared and dedicated plan types.
// | Operation | State | Plan | Config | Behavior |
// | --------- | ----- | ---- | ------ | -------- |
// | Create | Null | Has value | Has value | Validate plan/region, skip type change check |
// | Update | Has value | Has value | Has value | Validate plan/region, check type change |
// | Delete | Has value | Null | Has value | Skip immediatly |
// | Refresh | Has value | Has value | Has value | Validate plan/region if change detected |
func (r *instanceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Info(ctx, "Entering ModifyPlan for instanceResource")
	// Skip if resource is being deleted
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan, state instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.State.Raw.IsNull() {
		// New resource being created
		if err := r.client.ValidatePlan(plan.Plan.ValueString()); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("plan"), "Invalid Plan", err.Error())
		}
		if err := r.client.ValidateRegion(plan.Region.ValueString()); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("region"), "Invalid Region", err.Error())
		}
	} else {
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		// Check only if plan or region changed
		if plan.Plan.ValueString() != state.Plan.ValueString() {
			if err := r.client.ValidatePlan(plan.Plan.ValueString()); err != nil {
				resp.Diagnostics.AddAttributeError(path.Root("plan"), "Invalid Plan", err.Error())
			}
		}
		if plan.Region.ValueString() != state.Region.ValueString() {
			if err := r.client.ValidateRegion(plan.Region.ValueString()); err != nil {
				resp.Diagnostics.AddAttributeError(path.Root("region"), "Invalid Region", err.Error())
			}
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Skip plan type change check if creating new resource
	if req.State.Raw.IsNull() {
		return
	}

	// Replace when changing between shared and dedicated plan types
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldPlan := state.Plan.ValueString()
	newPlan := plan.Plan.ValueString()

	if oldPlan != newPlan {
		oldPlanType := r.isSharedPlan(oldPlan)
		newPlanType := r.isSharedPlan(newPlan)

		if oldPlanType != newPlanType {
			resp.RequiresReplace = append(resp.RequiresReplace, path.Root("plan"))
		}
	}
}

func (r *instanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Build create request
	createReq := model.InstanceCreateRequest{
		Name:   plan.Name.ValueString(),
		Plan:   plan.Plan.ValueString(),
		Region: plan.Region.ValueString(),
	}

	// Optional fields
	if !plan.Tags.IsUnknown() {
		var tags []string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Tags = tags
	}

	if !plan.VpcID.IsUnknown() && plan.VpcID.ValueInt64() > 0 {
		createReq.VpcID = plan.VpcID.ValueInt64Pointer()
	}

	if !plan.VpcSubnet.IsUnknown() {
		createReq.VpcSubnet = plan.VpcSubnet.ValueString()
	}

	if !plan.RmqVersion.IsUnknown() {
		createReq.RmqVersion = plan.RmqVersion.ValueString()
	}

	if !plan.NoDefaultAlarms.IsUnknown() {
		noDefaultAlarms := plan.NoDefaultAlarms.ValueBool()
		createReq.NoDefaultAlarms = &noDefaultAlarms
	} else {
		plan.NoDefaultAlarms = types.BoolValue(false)
	}

	if !plan.Nodes.IsUnknown() && plan.Nodes.ValueInt64() > 0 {
		planName := plan.Plan.ValueString()
		// Don't send if plan is shared or not legacy
		if !r.isSharedPlan(planName) || r.isLegacyPlan(planName) {
			createReq.Nodes = plan.Nodes.ValueInt64Pointer()
		}
	}

	if !plan.PreferredAz.IsUnknown() {
		var preferredAz *[]string
		resp.Diagnostics.Append(plan.PreferredAz.ElementsAs(ctx, &preferredAz, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.PreferredAz = preferredAz
	}

	// Handle copy_settings
	if !plan.CopySettings.IsUnknown() && len(plan.CopySettings.Elements()) > 0 {
		var copySettings []copySettingsModel
		resp.Diagnostics.Append(plan.CopySettings.ElementsAs(ctx, &copySettings, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// In next major version make it a single nested block, keep for backward compatibility
		if len(copySettings) > 0 {
			cs := copySettings[0]
			subscriptionID, err := strconv.ParseInt(cs.SubscriptionID.ValueString(), 10, 64)
			if err != nil {
				resp.Diagnostics.AddError(
					"Invalid subscription_id",
					fmt.Sprintf("Could not convert subscription_id to int64: %s", err),
				)
				return
			}

			var settings []string
			resp.Diagnostics.Append(cs.Settings.ElementsAs(ctx, &settings, false)...)
			if resp.Diagnostics.HasError() {
				return
			}

			createReq.CopySettings = &model.CopySettingsRequest{
				InstanceID: subscriptionID,
				Settings:   settings,
			}
		}
	}

	data, err := r.client.CreateInstance(timeoutCtx, createReq, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Instance",
			fmt.Sprintf("Could not create instance: %s", err),
		)
		return
	}

	r.populateResourceModel(ctx, data, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values
	if state.Sleep.IsNull() {
		state.Sleep = types.Int64Value(30)
	}
	if state.Timeout.IsNull() {
		state.Timeout = types.Int64Value(3600)
	}
	if state.KeepAssociatedVpc.IsNull() {
		state.KeepAssociatedVpc = types.BoolValue(false)
	}

	sleep := time.Duration(state.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(state.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	instanceID := state.ID.ValueString()

	data, err := r.client.ReadInstance(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Instance",
			fmt.Sprintf("Could not read instance with ID %s: %s", instanceID, err),
		)
		return
	}

	if data == nil {
		tflog.Warn(ctx, fmt.Sprintf("Resource drift: Instance: %s not found, removing from state", instanceID))
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateResourceModel(ctx, data, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state instanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	instanceID := plan.ID.ValueString()
	changed := false
	updateReq := model.InstanceUpdateRequest{}

	if !plan.Name.IsNull() && plan.Name.ValueString() != state.Name.ValueString() {
		updateReq.Name = plan.Name.ValueString()
		changed = true
	}
	if !plan.Plan.IsNull() && plan.Plan.ValueString() != state.Plan.ValueString() {
		updateReq.Plan = plan.Plan.ValueString()
		changed = true
	}
	if !plan.Nodes.IsNull() && plan.Nodes.ValueInt64() > 0 && plan.Nodes.ValueInt64() != state.Nodes.ValueInt64() {
		updateReq.Nodes = plan.Nodes.ValueInt64Pointer()
		changed = true
	}
	if !plan.Tags.IsNull() && !plan.Tags.Equal(state.Tags) {
		var tags []string
		resp.Diagnostics.Append(plan.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.Tags = tags
		changed = true
	}

	if !changed {
		tflog.Info(ctx, "No changes detected for instance update")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	err := r.client.UpdateInstance(timeoutCtx, instanceID, updateReq, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Instance",
			fmt.Sprintf("Could not update instance with ID %s: %s", instanceID, err),
		)
		return
	}

	data, err := r.client.ReadInstance(timeoutCtx, instanceID, sleep)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Instance",
			fmt.Sprintf("Could not create instance: %s", err),
		)
		return
	}

	r.populateResourceModel(ctx, data, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state instanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.ID.ValueString()
	keepVpc := state.KeepAssociatedVpc.ValueBool()

	err := r.client.DeleteInstance(ctx, instanceID, keepVpc)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Instance",
			fmt.Sprintf("Could not delete instance with ID %s: %s", instanceID, err),
		)
		return
	}
}

func (r *instanceResource) populateResourceModel(ctx context.Context, data *model.InstanceResponse, resourceModel *instanceResourceModel) {
	resourceModel.ID = types.StringValue(strconv.FormatInt(data.ID, 10))
	resourceModel.Name = types.StringValue(data.Name)
	resourceModel.Plan = types.StringValue(data.Plan)
	resourceModel.Region = types.StringValue(data.Region)
	resourceModel.Ready = types.BoolValue(data.Ready)
	resourceModel.ApiKey = types.StringValue(data.ApiKey)
	resourceModel.Url = types.StringValue(data.Url)
	resourceModel.Backend = types.StringValue(data.Backend)
	resourceModel.Nodes = types.Int64Value(data.Nodes)
	resourceModel.RmqVersion = types.StringValue(data.RmqVersion)
	resourceModel.Dedicated = types.BoolValue(data.Nodes > 0)
	resourceModel.Host = types.StringValue(data.HostnameExternal)
	resourceModel.HostInternal = types.StringValue(data.HostnameInternal)

	if data.Url != "" {
		urlInfo := r.client.UrlInformation(data.Url)
		if vhost, ok := urlInfo["vhost"].(string); ok {
			resourceModel.Vhost = types.StringValue(vhost)
		}
	}

	if data.VPC != nil {
		resourceModel.VpcID = types.Int64Value(data.VPC.ID)
		if data.VPC.Subnet != "" {
			resourceModel.VpcSubnet = types.StringValue(data.VPC.Subnet)
		}
	} else {
		resourceModel.VpcID = types.Int64Null()
		resourceModel.VpcSubnet = types.StringNull()
	}

	if len(data.Tags) > 0 {
		tags, _ := types.ListValueFrom(ctx, types.StringType, data.Tags)
		resourceModel.Tags = tags
	} else {
		emptyTags := []string{}
		tags, _ := types.ListValueFrom(ctx, types.StringType, emptyTags)
		resourceModel.Tags = tags
	}
}

// isSharedPlan returns true if the plan is a shared plan type
func (r *instanceResource) isSharedPlan(plan string) bool {
	sharedPlans := []string{"lemur", "tiger", "lemming", "ermine"}
	for _, p := range sharedPlans {
		if plan == p {
			return true
		}
	}
	return false
}

// isLegacyPlan returns true if the plan is a legacy plan type
func (r *instanceResource) isLegacyPlan(plan string) bool {
	legacyPlans := []string{"bunny", "rabbit", "panda", "ape", "hippo", "lion"}
	for _, p := range legacyPlans {
		if plan == p {
			return true
		}
	}
	return false
}
