package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/integrations"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &integrationLogAgentResource{}
	_ resource.ResourceWithConfigure   = &integrationLogAgentResource{}
	_ resource.ResourceWithImportState = &integrationLogAgentResource{}
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
	Cloudwatch *cloudwatchModel `tfsdk:"cloudwatch"`
}

type cloudwatchModel struct {
	IAMRole       types.String `tfsdk:"iam_role"`
	IAMExternalID types.String `tfsdk:"iam_external_id"`
	Region        types.String `tfsdk:"region"`
	LogGroupName  types.String `tfsdk:"log_group_name"`
	LogStreamName types.String `tfsdk:"log_stream_name"`
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
						Required:    true,
						Description: "AWS IAM role ARN",
					},
					"iam_external_id": schema.StringAttribute{
						Required:    true,
						Description: "External identifier that matches the role you created",
					},
					"region": schema.StringAttribute{
						Required:    true,
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

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	instanceID := plan.InstanceID.ValueInt64()
	request := r.populateRequest(&plan, intType)

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

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	id := plan.ID.ValueString()
	instanceID := plan.InstanceID.ValueInt64()
	request := r.populateRequest(&plan, intType)

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

// getIntegrationType returns the API type string based on which block is populated
func (r *integrationLogAgentResource) getIntegrationType(m *integrationLogAgentResourceModel) (string, error) {
	if m.Cloudwatch != nil {
		return "cloudwatch_v2", nil
	}
	return "", fmt.Errorf("exactly one integration block must be set (e.g. cloudwatch)")
}

// populateRequest converts the resource model to an API request
func (r *integrationLogAgentResource) populateRequest(plan *integrationLogAgentResourceModel, intType string) model.LogAgentRequest {
	switch intType {
	case "cloudwatch_v2":
		cw := plan.Cloudwatch
		req := model.LogAgentRequest{
			Region:        cw.Region.ValueString(),
			IAMRole:       cw.IAMRole.ValueString(),
			IAMExternalID: cw.IAMExternalID.ValueString(),
		}
		if !cw.LogGroupName.IsNull() && !cw.LogGroupName.IsUnknown() {
			req.LogGroupName = cw.LogGroupName.ValueString()
		}
		if !cw.LogStreamName.IsNull() && !cw.LogStreamName.IsUnknown() {
			req.LogStreamName = cw.LogStreamName.ValueString()
		}
		return req
	}
	return model.LogAgentRequest{}
}

// populateResourceModel fills the resource model from the API response
func (r *integrationLogAgentResource) populateResourceModel(m *integrationLogAgentResourceModel, data *model.LogAgentResponse) {
	switch data.Type {
	case "cloudwatch_v2":
		if m.Cloudwatch == nil {
			m.Cloudwatch = &cloudwatchModel{}
		}
		m.Cloudwatch.IAMRole = types.StringPointerValue(data.Config.IAMRole)
		m.Cloudwatch.IAMExternalID = types.StringPointerValue(data.Config.IAMExternalID)
		m.Cloudwatch.Region = types.StringPointerValue(data.Config.Region)
		m.Cloudwatch.LogGroupName = types.StringPointerValue(data.Config.LogGroupName)
		m.Cloudwatch.LogStreamName = types.StringPointerValue(data.Config.LogStreamName)
	}
}
