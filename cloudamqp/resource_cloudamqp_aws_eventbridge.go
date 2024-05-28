package cloudamqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type awsEventBridgeResource struct {
	client *api.API
}

type awsEventBridgeResourceModel struct {
	Id           types.String `tfsdk:"id"`
	InstanceID   types.Int64  `tfsdk:"instance_id"`
	AwsAccountId types.String `tfsdk:"aws_account_id"`
	AwsRegion    types.String `tfsdk:"aws_region"`
	Vhost        types.String `tfsdk:"vhost"`
	QueueName    types.String `tfsdk:"queue"`
	WithHeaders  types.Bool   `tfsdk:"with_headers"`
	Status       types.String `tfsdk:"status"`
}

type awsEventBridgeResourceApiModel struct {
	AwsAccountId string `json:"aws_account_id"`
	AwsRegion    string `json:"aws_region"`
	Vhost        string `json:"vhost"`
	QueueName    string `json:"queue"`
	WithHeaders  bool   `json:"with_headers"`
}

func (r *awsEventBridgeResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if request.ProviderData == nil {
		return
	}

	client, ok := request.ProviderData.(*api.API)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", request.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *awsEventBridgeResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "cloudamqp_integration_aws_eventbridge"
}

func (r *awsEventBridgeResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"aws_account_id": schema.StringAttribute{
				Required:    true,
				Description: "The 12 digit AWS Account ID where you want the events to be sent to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"aws_region": schema.StringAttribute{
				Required:    true,
				Description: "The AWS region where you the events to be sent to. (e.g. us-west-1, us-west-2, ..., etc.)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vhost": schema.StringAttribute{
				Required:    true,
				Description: "The VHost the queue resides in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"queue": schema.StringAttribute{
				Required:    true,
				Description: "A (durable) queue on your RabbitMQ instance.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"with_headers": schema.BoolAttribute{
				Required:    true,
				Description: "Include message headers in the event data.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Always set to null, unless there is an error starting the EventBridge",
			},
		},
	}
}

func (r *awsEventBridgeResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var data awsEventBridgeResourceModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.Plan.Get(ctx, &data)...)

	if response.Diagnostics.HasError() {
		return
	}

	apiModel := awsEventBridgeResourceApiModel{
		AwsAccountId: data.AwsAccountId.ValueString(),
		AwsRegion:    data.AwsRegion.ValueString(),
		Vhost:        data.Vhost.ValueString(),
		QueueName:    data.QueueName.ValueString(),
		WithHeaders:  data.WithHeaders.ValueBool(),
	}

	var params map[string]interface{}
	temp, err := json.Marshal(apiModel)
	if err != nil {
		response.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while creating the resource create request. "+
				"Please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error(),
		)
		return
	}
	// TODO: This is totally a hack to get the struct into a map[string]interface{}
	// It is very unlikely this will fail after the first one succeeds, so it should be fine to ignore the error
	// Maybe after the api is moved into the repo we can improve the interface
	_ = json.Unmarshal(temp, &params)

	apiResponse, err := r.client.CreateAwsEventBridge(int(data.InstanceID.ValueInt64()), params)
	if err != nil {
		response.Diagnostics.AddError(
			"Failed to Create Resource",
			"An error occurred while calling the api to create the surface, verify your permissions are correct.\n\n"+
				"JSON Error: "+err.Error(),
		)
		return
	}

	data.Id = types.StringValue(apiResponse["id"].(string))
	data.Status = types.StringNull()

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &data)...)
}

func (r *awsEventBridgeResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var state awsEventBridgeResourceModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if strings.Contains(state.Id.ValueString(), ",") {
		log.Printf("[DEBUG] cloudamqp::resource::aws-eventbridge::read id contains : %v", state.Id.String())
		s := strings.Split(state.Id.ValueString(), ",")
		log.Printf("[DEBUG] cloudamqp::resource::aws-eventbridge::read split ids: %v, %v", s[0], s[1])
		state.Id = types.StringValue(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		state.InstanceID = types.Int64Value(int64(instanceID))
	}
	if state.InstanceID.ValueInt64() == 0 {
		response.Diagnostics.AddError("Missing instance identifier {resource_id},{instance_id}", "")
		return
	}

	var (
		id         = state.Id.ValueString()
		instanceID = int(state.InstanceID.ValueInt64())
	)

	log.Printf("[DEBUG] cloudamqp::resource::aws-eventbridge::read ID: %v, instanceID %v", id, instanceID)
	data, err := r.client.ReadAwsEventBridge(instanceID, id)
	if err != nil {
		response.Diagnostics.AddError("Something went wrong while reading the aws event bridge", fmt.Sprintf("%v", err))
		return
	}

	state.AwsAccountId = types.StringValue(data["aws_account_id"].(string))
	state.AwsRegion = types.StringValue(data["aws_region"].(string))
	state.Vhost = types.StringValue(data["vhost"].(string))
	state.QueueName = types.StringValue(data["queue"].(string))
	state.WithHeaders = types.BoolValue(data["with_headers"].(bool))

	// Save data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &state)...)

	return
}

func (r *awsEventBridgeResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// This resource does not implement the Update function
}

func (r *awsEventBridgeResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var data awsEventBridgeResourceModel

	// Read Terraform plan data into the model
	response.Diagnostics.Append(request.State.Get(ctx, &data)...)
	var id = data.Id.ValueString()
	err := r.client.DeleteAwsEventBridge(int(data.InstanceID.ValueInt64()), id)

	if err != nil {
		response.Diagnostics.AddError("An error occurred while deleting cloudamqp_integration_aws_eventbridge",
			fmt.Sprintf("Error deleting Cloudamqp event bridge %s: %s", id, err),
		)
	}
}

func (r *awsEventBridgeResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}
