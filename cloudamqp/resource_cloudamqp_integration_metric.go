package cloudamqp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &integrationMetricResource{}
	_ resource.ResourceWithConfigure   = &integrationMetricResource{}
	_ resource.ResourceWithImportState = &integrationMetricResource{}
)

type integrationMetricResource struct {
	client *api.API
}

func NewIntegrationMetricResource() resource.Resource {
	return &integrationMetricResource{}
}

type integrationMetricResourceModel struct {
	ID              types.String `tfsdk:"id"`
	InstanceID      types.Int64  `tfsdk:"instance_id"`
	Name            types.String `tfsdk:"name"`
	AccessKeyID     types.String `tfsdk:"access_key_id"`
	ApiKey          types.String `tfsdk:"api_key"`
	ClientEmail     types.String `tfsdk:"client_email"`
	Credentials     types.String `tfsdk:"credentials"`
	Email           types.String `tfsdk:"email"`
	IAMExternalID   types.String `tfsdk:"iam_external_id"`
	IAMRole         types.String `tfsdk:"iam_role"`
	IncludeAdQueues types.Bool   `tfsdk:"include_ad_queues"`
	PrivateKey      types.String `tfsdk:"private_key"`
	PrivateKeyID    types.String `tfsdk:"private_key_id"`
	ProjectID       types.String `tfsdk:"project_id"`
	QueueAllowlist  types.String `tfsdk:"queue_allowlist"`
	Region          types.String `tfsdk:"region"`
	Tags            types.String `tfsdk:"tags"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	VhostAllowlist  types.String `tfsdk:"vhost_allowlist"`
}

func (r *integrationMetricResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_integration_metric"
}

func (r *integrationMetricResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationMetricResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID of the integration metric",
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
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of log integration",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cloudwatch",
						"cloudwatch_v2",
						"datadog",
						"datadog_v2",
						"librato",
						"newrelic_v2",
						"stackdriver",
					),
				},
			},
			"access_key_id": schema.StringAttribute{
				Optional:    true,
				Description: "AWS access key identifier. (Cloudwatch)",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The API key for the integration service. (Librato, Data Dog, New Relic)",
			},
			"client_email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The client email. (Stackdriver)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"credentials": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Base64Encoded credentials. (Stackdriver)",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Description: "The email address registred for the integration service. (Librato)",
			},
			"iam_external_id": schema.StringAttribute{
				Optional:    true,
				Description: "External identifier that match the role you created. (Cloudwatch)",
			},
			"iam_role": schema.StringAttribute{
				Optional:    true,
				Description: "The ARN of the role to be assumed when publishing metrics. (Cloudwatch)",
			},
			"include_ad_queues": schema.BoolAttribute{
				Optional:    true,
				Description: "(optional) Include Auto-Delete queues",
			},
			"private_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The private key. (Stackdriver)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "Private key identifier. (Stackdriver)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Project ID. (Stackdriver)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"queue_allowlist": schema.StringAttribute{
				Optional:    true,
				Description: "(optional) allowlist using regular expression",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "AWS region for Cloudwatch and [US/EU] for Data dog/New relic. (Cloudwatch, Data Dog, New Relic)",
			},
			"secret_access_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "AWS secret key. (Cloudwatch)",
			},
			"tags": schema.StringAttribute{
				Optional:    true,
				Description: "(optional) tags. E.g. env=prod,region=europe",
			},
			"vhost_allowlist": schema.StringAttribute{
				Optional:    true,
				Description: "(optional) allowlist using regular expression",
			},
		},
	}
}

func (r *integrationMetricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("ImportState: ID=%s", req.ID))
	if !strings.Contains(req.ID, ",") {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {resource_id},{instance_id}")
		return
	}

	idSplit := strings.Split(req.ID, ",")
	if len(idSplit) != 2 {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {resource_id},{instance_id}")
		return
	}
	instanceID, err := strconv.Atoi(idSplit[1])
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance_id in import ID", fmt.Sprintf("Could not convert instance_id to int: %s", err))
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), idSplit[0])
	resp.State.SetAttribute(ctx, path.Root("instance_id"), int64(instanceID))
}

func (r *integrationMetricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan integrationMetricResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	instanceID := plan.InstanceID.ValueInt64()
	type_ := plan.Name.ValueString()
	request := r.populateRequest(&plan)

	id, err := r.client.CreateIntegrationMetric(timeoutCtx, instanceID, type_, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Metric Integration",
			fmt.Sprintf("Could not create metric integration: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(id)
	r.computedValues(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *integrationMetricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state integrationMetricResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	instanceID := state.InstanceID.ValueInt64()
	metricID := state.ID.ValueString()

	data, err := r.client.ReadIntegrationMetric(timeoutCtx, instanceID, metricID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Metric Integration",
			fmt.Sprintf("Could not read metric integration: %s", err),
		)
		return
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("metric integration not found, resource will be recreated: %s", state.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateResourceModel(&state, data)
	r.computedValues(&state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *integrationMetricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan integrationMetricResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	instanceID := plan.InstanceID.ValueInt64()
	metricID := plan.ID.ValueString()
	request := r.populateRequest(&plan)

	err := r.client.UpdateIntegrationMetric(timeoutCtx, instanceID, metricID, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Metric Integration",
			fmt.Sprintf("Could not update metric integration: %s", err),
		)
		return
	}

	r.computedValues(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *integrationMetricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationMetricResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	instanceID := state.InstanceID.ValueInt64()
	metricID := state.ID.ValueString()

	err := r.client.DeleteIntegrationMetric(timeoutCtx, instanceID, metricID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Metric Integration",
			fmt.Sprintf("Could not delete metric integration: %s", err),
		)
		return
	}
}

// Compute values that are not user set to keep backwards compatibility
func (r *integrationMetricResource) computedValues(resourceModel *integrationMetricResourceModel) {
	if resourceModel.Name.ValueString() != "stackdriver" {
		resourceModel.ProjectID = types.StringNull()
		resourceModel.PrivateKey = types.StringNull()
		resourceModel.ClientEmail = types.StringNull()
		resourceModel.PrivateKeyID = types.StringNull()
	} else if resourceModel.Name.ValueString() == "stackdriver" {
		if resourceModel.Credentials.ValueString() != "" {
			resourceModel.ProjectID = types.StringNull()
			resourceModel.PrivateKey = types.StringNull()
			resourceModel.ClientEmail = types.StringNull()
			resourceModel.PrivateKeyID = types.StringNull()
		} else {
			resourceModel.PrivateKeyID = types.StringValue("")
		}
	}
}

// Handle data conversion from API response to resource model
func (r *integrationMetricResource) populateResourceModel(resourceModel *integrationMetricResourceModel, data *model.MetricResponse) {
	resourceModel.Name = types.StringValue(data.Type)

	switch data.Type {
	case "cloudwatch", "cloudwatch_v2":
		resourceModel.Region = types.StringValue(*data.Config.Region)
		if data.Config.AccessKeyID != nil {
			resourceModel.AccessKeyID = types.StringValue(*data.Config.AccessKeyID)
		}
		if data.Config.SecretAccessKey != nil {
			resourceModel.SecretAccessKey = types.StringValue(*data.Config.SecretAccessKey)
		}
		if data.Config.IAMExternalID != nil {
			resourceModel.IAMExternalID = types.StringValue(*data.Config.IAMExternalID)
		}
		if data.Config.IAMRole != nil {
			resourceModel.IAMRole = types.StringValue(*data.Config.IAMRole)
		}
	case "datadog", "datadog_v2":
		resourceModel.Region = types.StringValue(*data.Config.Region)
		resourceModel.ApiKey = types.StringValue(*data.Config.APIKey)
	case "librato":
		resourceModel.Email = types.StringValue(*data.Config.Email)
		resourceModel.ApiKey = types.StringValue(*data.Config.APIKey)
	case "newrelic_v2":
		resourceModel.ApiKey = types.StringValue(*data.Config.APIKey)
		resourceModel.Region = types.StringValue(*data.Config.Region)
	case "stackdriver":
		if resourceModel.Credentials.ValueString() == "" {
			resourceModel.ClientEmail = types.StringValue(*data.Config.ClientEmail)
			resourceModel.PrivateKey = types.StringValue(*data.Config.PrivateKey)
			resourceModel.ProjectID = types.StringValue(*data.Config.ProjectID)
		}
	}

	// Set commons values
	if data.Config.IncludeAdQueues != nil {
		resourceModel.IncludeAdQueues = types.BoolValue(*data.Config.IncludeAdQueues)
	} else {
		resourceModel.IncludeAdQueues = types.BoolNull()
	}
	if data.Config.QueueRegex != nil {
		resourceModel.QueueAllowlist = types.StringValue(*data.Config.QueueRegex)
	} else {
		resourceModel.QueueAllowlist = types.StringNull()
	}
	if data.Config.VhostRegex != nil {
		resourceModel.VhostAllowlist = types.StringValue(*data.Config.VhostRegex)
	} else {
		resourceModel.VhostAllowlist = types.StringNull()
	}
	if data.Config.Tags != nil {
		resourceModel.Tags = types.StringValue(*data.Config.Tags)
	} else {
		resourceModel.Tags = types.StringNull()
	}
}

// Handle data conversion from resource model to API request
func (r *integrationMetricResource) populateRequest(plan *integrationMetricResourceModel) model.MetricRequest {
	var request model.MetricRequest

	switch plan.Name.ValueString() {
	case "cloudwatch", "cloudwatch_v2":
		request.Region = plan.Region.ValueString()
		if !plan.AccessKeyID.IsUnknown() {
			request.AccessKeyID = plan.AccessKeyID.ValueString()
		}
		if !plan.SecretAccessKey.IsUnknown() {
			request.SecretAccessKey = plan.SecretAccessKey.ValueString()
		}
		if !plan.IAMExternalID.IsUnknown() {
			request.IAMExternalID = plan.IAMExternalID.ValueString()
		}
		if !plan.IAMRole.IsUnknown() {
			request.IAMRole = plan.IAMRole.ValueString()
		}
	case "datadog", "datadog_v2":
		request.APIKey = plan.ApiKey.ValueString()
		request.Region = plan.Region.ValueString()
	case "librato":
		request.APIKey = plan.ApiKey.ValueString()
		request.Email = plan.Email.ValueString()
	case "newrelic_v2":
		request.APIKey = plan.ApiKey.ValueString()
		request.Region = plan.Region.ValueString()
	case "stackdriver":
		if plan.Credentials.ValueString() != "" {
			uDec, _ := base64.URLEncoding.DecodeString(plan.Credentials.ValueString())
			var jsonMap map[string]any
			json.Unmarshal([]byte(uDec), &jsonMap)
			request.ClientEmail = jsonMap["client_email"].(string)
			request.PrivateKeyID = jsonMap["private_key_id"].(string)
			request.PrivateKey = jsonMap["private_key"].(string)
			request.ProjectID = jsonMap["project_id"].(string)
		} else {
			request.ClientEmail = plan.ClientEmail.ValueString()
			request.PrivateKeyID = plan.PrivateKeyID.ValueString()
			request.PrivateKey = plan.PrivateKey.ValueString()
			request.ProjectID = plan.ProjectID.ValueString()
		}
	}

	// Set common values
	if !plan.IncludeAdQueues.IsUnknown() {
		request.IncludeAdQueues = plan.IncludeAdQueues.ValueBool()
	}
	if !plan.QueueAllowlist.IsUnknown() {
		request.QueueRegex = plan.QueueAllowlist.ValueString()
	}
	if !plan.Tags.IsUnknown() {
		request.Tags = plan.Tags.ValueString()
	}
	if !plan.VhostAllowlist.IsUnknown() {
		request.VhostRegex = plan.VhostAllowlist.ValueString()
	}

	return request
}
