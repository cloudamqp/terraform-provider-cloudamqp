package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                     = &integrationLogAgentResource{}
	_ resource.ResourceWithConfigure        = &integrationLogAgentResource{}
	_ resource.ResourceWithImportState      = &integrationLogAgentResource{}
	_ resource.ResourceWithConfigValidators = &integrationLogAgentResource{}
)

type integrationLogAgentResource struct {
	client *api.API
}

func NewIntegrationLogAgentResource() resource.Resource {
	return &integrationLogAgentResource{}
}

type integrationLogAgentResourceModel struct {
	ID         types.String     `tfsdk:"id"`
	InstanceID types.Int64      `tfsdk:"instance_id"`
	Cloudwatch  *cloudwatchModel  `tfsdk:"cloudwatch"`
	Uptrace     *uptraceModel     `tfsdk:"uptrace"`
	Splunk      *splunkModel      `tfsdk:"splunk"`
	Coralogix   *coralogixModel   `tfsdk:"coralogix"`
	Datadog     *datadogModel     `tfsdk:"datadog"`
	CustomOTLP  *customOtlpModel  `tfsdk:"custom_otlp"`
	GoogleCloud *googleCloudModel `tfsdk:"google_cloud"`
	Grafana     *grafanaModel     `tfsdk:"grafana"`
}

type cloudwatchModel struct {
	IAMRole       types.String `tfsdk:"iam_role"`
	IAMExternalID types.String `tfsdk:"iam_external_id"`
	Region        types.String `tfsdk:"region"`
	LogGroupName  types.String `tfsdk:"log_group_name"`
	LogStreamName types.String `tfsdk:"log_stream_name"`
}

type uptraceModel struct {
	DSN types.String `tfsdk:"dsn"`
}

type splunkModel struct {
	Endpoint   types.String `tfsdk:"hec_endpoint"`
	Token      types.String `tfsdk:"token"`
	SourceType types.String `tfsdk:"source_type"`
}

type coralogixModel struct {
	PrivateKey  types.String `tfsdk:"private_key"`
	Application types.String `tfsdk:"application"`
	Subsystem   types.String `tfsdk:"subsystem"`
	Region      types.String `tfsdk:"region"`
}

type datadogModel struct {
	APIKey types.String `tfsdk:"api_key"`
	Region types.String `tfsdk:"region"`
	Tags   types.String `tfsdk:"tags"`
}

type customOtlpModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Headers  types.Map    `tfsdk:"headers"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type googleCloudModel struct {
	ServiceAccountFile        types.String `tfsdk:"service_account_file"`
	ServiceAccountFileVersion types.Int64  `tfsdk:"service_account_file_version"`
	ProjectID                 types.String `tfsdk:"project_id"`
	ClientEmail               types.String `tfsdk:"client_email"`
	PrivateKeyID              types.String `tfsdk:"private_key_id"`
	Tags                      types.String `tfsdk:"tags"`
}

type grafanaModel struct {
	Endpoint          types.String `tfsdk:"endpoint"`
	GrafanaInstanceID types.String `tfsdk:"grafana_instance_id"`
	APIToken          types.String `tfsdk:"api_token"`
}

func (r *integrationLogAgentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_integration_log_agent"
}

func (r *integrationLogAgentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Schema defines the schema for the resource
func (r *integrationLogAgentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID of the log integration",
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
		},
		Blocks: map[string]schema.Block{
			"cloudwatch": schema.SingleNestedBlock{
				Description: "CloudWatch OTLP log integration configuration",
				Attributes: map[string]schema.Attribute{
					"iam_role": schema.StringAttribute{
						Optional:    true,
						Description: "AWS IAM role ARN",
					},
					"iam_external_id": schema.StringAttribute{
						Optional:    true,
						Description: "External identifier that matches the role you created",
					},
					"region": schema.StringAttribute{
						Optional:    true,
						Description: "AWS region",
					},
					"log_group_name": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The name of the CloudWatch log group",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"log_stream_name": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The name of the CloudWatch log stream",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"uptrace": schema.SingleNestedBlock{
				Description: "Uptrace OTLP log integration configuration",
				Attributes: map[string]schema.Attribute{
					"dsn": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Uptrace DSN (Data Source Name) URL",
					},
				},
			},
			"splunk": schema.SingleNestedBlock{
				Description: "Splunk HEC log integration configuration",
				Attributes: map[string]schema.Attribute{
					"endpoint": schema.StringAttribute{
						Optional:    true,
						Description: "Splunk HEC endpoint URL (e.g. https://your-instance.splunkcloud.com:8088/services/collector)",
					},
					"token": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
						Description: "Splunk HEC token",
					},
					"source_type": schema.StringAttribute{
						Optional:    true,
						Description: "Splunk source type (leave empty to use the token's default)",
					},
				},
			},
			"coralogix": schema.SingleNestedBlock{
				Description: "Coralogix log integration configuration",
				Attributes: map[string]schema.Attribute{
					"private_key": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Coralogix private key (always starts with cxtp_...)",
					},
					"application": schema.StringAttribute{
						Optional:    true,
						Description: "Application name, used to group logs by environment",
					},
					"subsystem": schema.StringAttribute{
						Optional:    true,
						Description: "Subsystem name, used to group logs by service within an application",
					},
					"region": schema.StringAttribute{
						Optional:    true,
						Description: "Coralogix region (US1, US2, US3, EU1, EU2, AP1, AP2, AP3)",
						Validators: []validator.String{
							stringvalidator.OneOf("US1", "US2", "US3", "EU1", "EU2", "AP1", "AP2", "AP3"),
						},
					},
				},
			},
			"datadog": schema.SingleNestedBlock{
				Description: "Datadog log integration configuration",
				Attributes: map[string]schema.Attribute{
					"api_key": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
						Description: "Datadog API key",
					},
					"region": schema.StringAttribute{
						Optional:    true,
						Description: "Datadog region (US1, US3, US5, EU, AP2)",
						Validators: []validator.String{
							stringvalidator.OneOf("US1", "US3", "US5", "EU", "AP2"),
						},
					},
					"tags": schema.StringAttribute{
						Optional:    true,
						Description: "Comma-separated tags to attach to logs (e.g. env=prod,region=eu)",
					},
				},
			},
			"custom_otlp": schema.SingleNestedBlock{
				Description: "Custom OTLP log integration configuration",
				Attributes: map[string]schema.Attribute{
					"endpoint": schema.StringAttribute{
						Optional:    true,
						Description: "OTLP HTTP endpoint URL (e.g. http://otlp.uptrace.dev:4318)",
					},
					"headers": schema.MapAttribute{
						Optional:    true,
						Sensitive:   true,
						ElementType: types.StringType,
						Description: "Key-value HTTP headers for authentication (e.g. uptrace-dsn: https://token@api.uptrace.dev/project_id). Mutually exclusive with username/password.",
					},
					"username": schema.StringAttribute{
						Optional:    true,
						Description: "Username for HTTP basic auth. Must be set together with password. Mutually exclusive with headers.",
					},
					"password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
						Description: "Password for HTTP basic auth. Must be set together with username. Mutually exclusive with headers.",
					},
				},
			},
			"google_cloud": schema.SingleNestedBlock{
				Description: "Google Cloud log integration configuration",
				Attributes: map[string]schema.Attribute{
					"service_account_file": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
						Description: "Base64-encoded Google service account key JSON file",
					},
					"service_account_file_version": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(1),
						Description: "Version of the write-only service_account_file. Increment to trigger an update when the file contents change (default: 1).",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"project_id": schema.StringAttribute{
						Computed:    true,
						Description: "Google Cloud project ID (computed from service_account_file)",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"client_email": schema.StringAttribute{
						Computed:    true,
						Description: "Google service account client email (computed from service_account_file)",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"private_key_id": schema.StringAttribute{
						Computed:    true,
						Sensitive:   true,
						Description: "Google service account private key ID (computed from service_account_file)",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"tags": schema.StringAttribute{
						Optional:    true,
						Description: "Comma-separated tags to attach to logs (e.g. env=prod,region=eu)",
					},
				},
			},
			"grafana": schema.SingleNestedBlock{
				Description: "Grafana Cloud (Loki) log integration configuration",
				Attributes: map[string]schema.Attribute{
					"endpoint": schema.StringAttribute{
						Optional:    true,
						Description: "Grafana Loki endpoint URL",
					},
					"grafana_instance_id": schema.StringAttribute{
						Optional:    true,
						Description: "Grafana Cloud instance ID",
					},
					"api_token": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						WriteOnly:   true,
						Description: "Grafana Cloud API token",
					},
				},
			},
		},
	}
}

// ImportState imports the resource state from the API into the Terraform state
func (r *integrationLogAgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// Create creates the resource and sets the initial Terraform state
func (r *integrationLogAgentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	intType, err := r.getIntegrationType(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	// Read WriteOnly fields from config (not available in plan)
	var config integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.GoogleCloud != nil && !config.GoogleCloud.ServiceAccountFile.IsNull() {
		if plan.GoogleCloud == nil {
			plan.GoogleCloud = &googleCloudModel{}
		}
		plan.GoogleCloud.ServiceAccountFile = config.GoogleCloud.ServiceAccountFile
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	instanceID := plan.InstanceID.ValueInt64()
	request, err := r.populateRequest(&plan, intType)
	if err != nil {
		resp.Diagnostics.AddError("Invalid service_account_file", err.Error())
		return
	}

	id, err := r.client.CreateIntegrationLogAgent(timeoutCtx, instanceID, intType, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Log Integration",
			fmt.Sprintf("Could not create log integration: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read fetches the resource state from the API and updates the Terraform state accordingly
func (r *integrationLogAgentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	id := state.ID.ValueString()
	instanceID := state.InstanceID.ValueInt64()

	data, err := r.client.ReadIntegrationLogAgent(timeoutCtx, instanceID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Log Integration",
			fmt.Sprintf("Could not read log integration with ID %s: %s", id, err),
		)
		return
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("Log integration not found, resource will be recreated: %s", state.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateResourceModel(&state, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies the resource and updates the integration in the API
func (r *integrationLogAgentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	intType, err := r.getIntegrationType(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	// Read WriteOnly fields from config (not available in plan)
	var config integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.GoogleCloud != nil && !config.GoogleCloud.ServiceAccountFile.IsNull() {
		if plan.GoogleCloud == nil {
			plan.GoogleCloud = &googleCloudModel{}
		}
		plan.GoogleCloud.ServiceAccountFile = config.GoogleCloud.ServiceAccountFile
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	id := plan.ID.ValueString()
	instanceID := plan.InstanceID.ValueInt64()
	request, err := r.populateRequest(&plan, intType)
	if err != nil {
		resp.Diagnostics.AddError("Invalid service_account_file", err.Error())
		return
	}

	err = r.client.UpdateIntegrationLogAgent(timeoutCtx, instanceID, id, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Log Integration",
			fmt.Sprintf("Could not update log integration with ID %s: %s", id, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the resource from the state and deletes the integration from the API
func (r *integrationLogAgentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationLogAgentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	id := state.ID.ValueString()
	instanceID := state.InstanceID.ValueInt64()

	err := r.client.DeleteIntegrationLogAgent(timeoutCtx, instanceID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Log Integration",
			fmt.Sprintf("Could not delete log integration with ID %s: %s", id, err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}
