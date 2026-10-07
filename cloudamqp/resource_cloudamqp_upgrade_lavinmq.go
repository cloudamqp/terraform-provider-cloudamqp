package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &upgradeLavinMQResource{}
	_ resource.ResourceWithConfigure = &upgradeLavinMQResource{}
)

func NewUpgradeLavinMQResource() resource.Resource {
	return &upgradeLavinMQResource{}
}

type upgradeLavinMQResource struct {
	client *api.API
}

type upgradeLavinMQResourceModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	NewVersion types.String `tfsdk:"new_version"`
}

func (r *upgradeLavinMQResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_upgrade_lavinmq"
}

func (r *upgradeLavinMQResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Upgrade LavinMQ to the latest possible or a specific available version.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"new_version": schema.StringAttribute{
				Optional:    true,
				Description: "The new version to upgrade to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *upgradeLavinMQResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *upgradeLavinMQResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan upgradeLavinMQResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.client.UpgradeLavinMQ(ctx, plan.InstanceID.ValueInt64(), plan.NewVersion.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to upgrade LavinMQ", err.Error())
		return
	}

	if len(response) > 0 {
		tflog.Info(ctx, fmt.Sprintf("LavinMQ update result: %s", response))
	}

	plan.ID = types.StringValue(plan.InstanceID.String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *upgradeLavinMQResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

func (r *upgradeLavinMQResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *upgradeLavinMQResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(context.Background())
}
