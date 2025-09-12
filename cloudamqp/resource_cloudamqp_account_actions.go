package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &accountActionsResource{}
	_ resource.ResourceWithConfigure = &accountActionsResource{}
)

func NewAccountActionsResource() resource.Resource {
	return &accountActionsResource{}
}

type accountActionsResource struct {
	client *api.API
}

type accountActionsResourceModel struct {
	Id         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	Action     types.String `tfsdk:"action"`
}

func (r *accountActionsResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *accountActionsResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "cloudamqp_account_actions"
}

func (r *accountActionsResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Preform actions on your CloudAMQP account, such as rotating passwords or API keys.",
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
			"action": schema.StringAttribute{
				Required:    true,
				Description: "The action to perform on the node",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("rotate-password", "rotate-apikey", "enable-vpc"),
				},
			},
		},
	}
}

func (r *accountActionsResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan accountActionsResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	switch plan.Action.ValueString() {
	case "rotate-password":
		if err := r.client.RotatePassword(ctx, int(plan.InstanceID.ValueInt64())); err != nil {
			response.Diagnostics.AddError(
				"Rotate Password Error",
				fmt.Sprintf("An error occurred while rotating the password for instance %d: %s", plan.InstanceID, err),
			)
			return
		}
	case "rotate-apikey":
		if err := r.client.RotateApiKey(ctx, int(plan.InstanceID.ValueInt64())); err != nil {
			response.Diagnostics.AddError(
				"Rotate API Key Error",
				fmt.Sprintf("An error occurred while rotating the API key for instance %d: %s", plan.InstanceID, err),
			)
			return
		}
	case "enable-vpc":
		if err := r.client.EnableVPC(ctx, int(plan.InstanceID.ValueInt64())); err != nil {
			response.Diagnostics.AddError(
				"Enable VPC Error",
				fmt.Sprintf("An error occurred while enabling the VPC for instance %d: %s", plan.InstanceID, err),
			)
			return
		}
	}

	plan.Id = types.StringValue(plan.InstanceID.String())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *accountActionsResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	// This resource does not implement the Read function
}

func (r *accountActionsResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// This resource does not implement the Update function
}

func (r *accountActionsResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// This resource does not implement the Delete function
	response.State.RemoveResource(ctx)
}
