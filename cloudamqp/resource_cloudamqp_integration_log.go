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
	_ resource.Resource                = &integrationLogResource{}
	_ resource.ResourceWithConfigure   = &integrationLogResource{}
	_ resource.ResourceWithImportState = &integrationLogResource{}
)

type integrationLogResource struct {
	client *api.API
}

func NewIntegrationLogResource() resource.Resource {
	return &integrationLogResource{}
}

type integrationLogResourceModel struct {
	ID                types.String `tfsdk:"id"`
	InstanceID        types.Int64  `tfsdk:"instance_id"`
	Name              types.String `tfsdk:"name"`
	Url               types.String `tfsdk:"url"`
	HostPort          types.String `tfsdk:"host_port"`
	Token             types.String `tfsdk:"token"`
	Region            types.String `tfsdk:"region"`
	AccessKeyID       types.String `tfsdk:"access_key_id"`
	SecretAccessKey   types.String `tfsdk:"secret_access_key"`
	ApiKey            types.String `tfsdk:"api_key"`
	Tags              types.String `tfsdk:"tags"`
	ProjectID         types.String `tfsdk:"project_id"`
	PrivateKey        types.String `tfsdk:"private_key"`
	ClientEmail       types.String `tfsdk:"client_email"`
	Host              types.String `tfsdk:"host"`
	SourceType        types.String `tfsdk:"sourcetype"`
	PrivateKeyID      types.String `tfsdk:"private_key_id"`
	Credentials       types.String `tfsdk:"credentials"`
	Endpoint          types.String `tfsdk:"endpoint"`
	Application       types.String `tfsdk:"application"`
	Subsystem         types.String `tfsdk:"subsystem"`
	TenantID          types.String `tfsdk:"tenant_id"`
	ApplicationID     types.String `tfsdk:"application_id"`
	ApplicationSecret types.String `tfsdk:"application_secret"`
	DceURI            types.String `tfsdk:"dce_uri"`
	Table             types.String `tfsdk:"table"`
	DcrID             types.String `tfsdk:"dcr_id"`
	Retention         types.Int64  `tfsdk:"retention"`
}

func (r *integrationLogResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_integration_log"
}

func (r *integrationLogResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *integrationLogResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID of the webhook",
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
						"azure_monitor",
						"cloudwatchlog",
						"coralogix",
						"datadog",
						"logentries",
						"loggly",
						"papertrail",
						"scalyr",
						"splunk",
						"stackdriver",
					),
				},
			},
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "The URL to push the logs to. (Papertrail)",
			},
			"host_port": schema.StringAttribute{
				Optional:    true,
				Description: "Destination to send the logs. (Splunk)",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The token used for authentication. (Loggly, Logentries, Splunk, Scalyr)",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "The region hosting integration service. (Cloudwatch, Datadog)",
			},
			"access_key_id": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "AWS access key identifier. (Cloudwatch)",
			},
			"secret_access_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "AWS secret access key. (Cloudwatch)",
			},
			"api_key": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The API key for the integration service. (Datadog)",
			},
			"tags": schema.StringAttribute{
				Optional:    true,
				Description: "Optional tags. E.g. env=prod,region=europe. (Cloudwatch, Datadog)",
			},
			"project_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The project ID for the integration service. (Stackdriver)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The private API key used for authentication. (Stackdriver, Coralogix)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The client email. (Stackdriver)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The host information. (Scalyr)",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"app.scalyr.com",
						"app.eu.scalyr.com",
					),
				},
			},
			"sourcetype": schema.StringAttribute{
				Optional:    true,
				Description: "Assign source type to the data exported, eg. generic_single_line. (Splunk)",
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
			"credentials": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Base64Encoded credentials. (Stackdriver)",
			},
			"endpoint": schema.StringAttribute{
				Optional:    true,
				Description: "The syslog destination to send the logs to. (Papertrail)",
			},
			"application": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the application. (Azure Monitor)",
			},
			"subsystem": schema.StringAttribute{
				Optional:    true,
				Description: "The name of the subsystem. (Azure Monitor)",
			},
			"tenant_id": schema.StringAttribute{
				Optional:    true,
				Description: "The tenant ID. (Azure Monitor)",
			},
			"application_id": schema.StringAttribute{
				Optional:    true,
				Description: "The application ID. (Azure Monitor)",
			},
			"application_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The application secret. (Azure Monitor)",
			},
			"dce_uri": schema.StringAttribute{
				Optional:    true,
				Description: "The DCE URI. (Coralogix)",
			},
			"table": schema.StringAttribute{
				Optional:    true,
				Description: "The table name to send the logs to. (Azure Monitor)",
			},
			"dcr_id": schema.StringAttribute{
				Optional:    true,
				Description: "The DCR ID. (Coralogix)",
			},
			"retention": schema.Int64Attribute{
				Optional:    true,
				Description: "The number of days to retain logs. (Cloudwatch)",
			},
		},
	}
}

func (r *integrationLogResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

func (r *integrationLogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan integrationLogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	instanceID := plan.InstanceID.ValueInt64()
	type_ := plan.Name.ValueString()
	request := r.populateRequest(&plan)

	id, err := r.client.CreateIntegrationLog(timeoutCtx, instanceID, type_, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Log Integration",
			fmt.Sprintf("Could not create log integration: %s", err),
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

func (r *integrationLogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state integrationLogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	id := state.ID.ValueString()
	instanceID := state.InstanceID.ValueInt64()

	data, err := r.client.ReadIntegrationLog(timeoutCtx, instanceID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Log Integration",
			fmt.Sprintf("Could not read log integration with ID %s: %s", id, err),
		)
		return
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("log integration not found, resource will be recreated: %s", state.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateResourceModel(&state, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *integrationLogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan integrationLogResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	id := plan.ID.ValueString()
	instanceID := plan.InstanceID.ValueInt64()
	request := r.populateRequest(&plan)

	err := r.client.UpdateIntegrationLog(timeoutCtx, instanceID, id, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Log Integration",
			fmt.Sprintf("Could not update log integration with ID %s: %s", id, err),
		)
		return
	}

	r.computedValues(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *integrationLogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state integrationLogResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	id := state.ID.ValueString()
	instanceID := state.InstanceID.ValueInt64()

	err := r.client.DeleteIntegrationLog(timeoutCtx, instanceID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Log Integration",
			fmt.Sprintf("Could not delete log integration with ID %s: %s", id, err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

// Compute values that are not user set to keep backwards compatibility
func (r *integrationLogResource) computedValues(resourceModel *integrationLogResourceModel) {
	if resourceModel.Name.ValueString() == "coralogix" {
		resourceModel.ProjectID = types.StringNull()
		resourceModel.ClientEmail = types.StringNull()
		resourceModel.PrivateKeyID = types.StringNull()
	} else if resourceModel.Name.ValueString() != "stackdriver" {
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
func (r *integrationLogResource) populateResourceModel(resourceModel *integrationLogResourceModel, data *model.LogResponse) {
	resourceModel.Name = types.StringValue(data.Type)

	switch data.Type {
	case "azure_monitor":
		resourceModel.TenantID = types.StringValue(*data.Config.TenantID)
		resourceModel.ApplicationID = types.StringValue(*data.Config.ApplicationID)
		resourceModel.ApplicationSecret = types.StringValue(*data.Config.ApplicationSecret)
		resourceModel.DceURI = types.StringValue(*data.Config.DCEURI)
		resourceModel.Table = types.StringValue(*data.Config.Table)
		resourceModel.DcrID = types.StringValue(*data.Config.DCRID)
	case "cloudwatchlog":
		resourceModel.Region = types.StringValue(*data.Config.Region)
		resourceModel.AccessKeyID = types.StringValue(*data.Config.AccessKeyID)
		resourceModel.SecretAccessKey = types.StringValue(*data.Config.SecretAccessKey)
		if data.Config.Retention != nil {
			resourceModel.Retention = types.Int64Value(*data.Config.Retention)
		} else {
			resourceModel.Retention = types.Int64Null()
		}
		if data.Config.Tags != nil {
			resourceModel.Tags = types.StringValue(*data.Config.Tags)
		} else {
			resourceModel.Tags = types.StringNull()
		}
	case "coralogix":
		resourceModel.PrivateKey = types.StringValue(*data.Config.PrivateKey)
		resourceModel.Endpoint = types.StringValue(*data.Config.Endpoint)
		resourceModel.Application = types.StringValue(*data.Config.Application)
		resourceModel.Subsystem = types.StringValue(*data.Config.Subsystem)
	case "datadog":
		resourceModel.Region = types.StringValue(*data.Config.Region)
		resourceModel.ApiKey = types.StringValue(*data.Config.APIKey)
		if data.Config.Tags != nil {
			resourceModel.Tags = types.StringValue(*data.Config.Tags)
		} else {
			resourceModel.Tags = types.StringNull()
		}
	case "logentries":
		resourceModel.Token = types.StringValue(*data.Config.Token)
	case "loggly":
		resourceModel.Token = types.StringValue(*data.Config.Token)
	case "papertrail":
		resourceModel.Url = types.StringValue(*data.Config.URL)
	case "scalyr":
		resourceModel.Token = types.StringValue(*data.Config.Token)
		resourceModel.Host = types.StringValue(*data.Config.Host)
	case "splunk":
		resourceModel.HostPort = types.StringValue(*data.Config.HostPort)
		resourceModel.Token = types.StringValue(*data.Config.Token)
		resourceModel.SourceType = types.StringValue(*data.Config.Sourcetype)
	case "stackdriver":
		if resourceModel.Credentials.ValueString() == "" {
			resourceModel.ClientEmail = types.StringValue(*data.Config.ClientEmail)
			resourceModel.PrivateKey = types.StringValue(*data.Config.PrivateKey)
			resourceModel.ProjectID = types.StringValue(*data.Config.ProjectID)
		}
	}
}

// Handle data conversion from resource model to API request
func (r *integrationLogResource) populateRequest(plan *integrationLogResourceModel) model.LogRequest {
	var request model.LogRequest

	switch plan.Name.ValueString() {
	case "azure_monitor":
		request = model.LogRequest{
			TenantID:          plan.TenantID.ValueString(),
			ApplicationID:     plan.ApplicationID.ValueString(),
			ApplicationSecret: plan.ApplicationSecret.ValueString(),
			DCEURI:            plan.DceURI.ValueString(),
			Table:             plan.Table.ValueString(),
			DCRID:             plan.DcrID.ValueString(),
		}
	case "cloudwatchlog":
		request = model.LogRequest{
			Region:          plan.Region.ValueString(),
			AccessKeyID:     plan.AccessKeyID.ValueString(),
			SecretAccessKey: plan.SecretAccessKey.ValueString(),
		}
		if !plan.Retention.IsNull() {
			request.Retention = plan.Retention.ValueInt64()
		}
		if !plan.Tags.IsNull() {
			request.Tags = plan.Tags.ValueString()
		}
	case "coralogix":
		request = model.LogRequest{
			PrivateKey:  plan.PrivateKey.ValueString(),
			Endpoint:    plan.Endpoint.ValueString(),
			Application: plan.Application.ValueString(),
			Subsystem:   plan.Subsystem.ValueString(),
		}
	case "datadog":
		request = model.LogRequest{
			Region: plan.Region.ValueString(),
			APIKey: plan.ApiKey.ValueString(),
			Tags:   plan.Tags.ValueString(),
		}
	case "logentries":
		request = model.LogRequest{
			Token: plan.Token.ValueString(),
		}
	case "loggly":
		request = model.LogRequest{
			Token: plan.Token.ValueString(),
		}
	case "papertrail":
		request = model.LogRequest{
			URL: plan.Url.ValueString(),
		}
	case "scalyr":
		request = model.LogRequest{
			Token: plan.Token.ValueString(),
			Host:  plan.Host.ValueString(),
		}
	case "splunk":
		request = model.LogRequest{
			HostPort:   plan.HostPort.ValueString(),
			Token:      plan.Token.ValueString(),
			Sourcetype: plan.SourceType.ValueString(),
		}
	case "stackdriver":
		if plan.Credentials.ValueString() != "" {
			uDec, _ := base64.URLEncoding.DecodeString(plan.Credentials.ValueString())
			var jsonMap map[string]any
			json.Unmarshal([]byte(uDec), &jsonMap)
			request = model.LogRequest{
				ClientEmail:  jsonMap["client_email"].(string),
				PrivateKeyID: jsonMap["private_key_id"].(string),
				PrivateKey:   jsonMap["private_key"].(string),
				ProjectID:    jsonMap["project_id"].(string),
			}
		} else {
			request = model.LogRequest{
				ClientEmail:  plan.ClientEmail.ValueString(),
				PrivateKeyID: plan.PrivateKeyID.ValueString(),
				PrivateKey:   plan.PrivateKey.ValueString(),
				ProjectID:    plan.ProjectID.ValueString(),
			}
		}
	}

	return request
}
