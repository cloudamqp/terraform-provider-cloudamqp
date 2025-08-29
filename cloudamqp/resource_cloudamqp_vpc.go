package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils/validators"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &vpcResource{}
	_ resource.ResourceWithConfigure   = &vpcResource{}
	_ resource.ResourceWithImportState = &vpcResource{}
)

type vpcResource struct {
	client *api.API
}

func NewVpcResource() resource.Resource {
	return &vpcResource{}
}

type vpcResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Region  types.String `tfsdk:"region"`
	Subnet  types.String `tfsdk:"subnet"`
	Tags    types.List   `tfsdk:"tags"`
	VpcName types.String `tfsdk:"vpc_name"`
}

func (r *vpcResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_vpc"
}

func (r *vpcResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID of the VPC instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the VPC instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "Region where the VPC is located",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnet": schema.StringAttribute{
				Required:    true,
				Description: "The VPC subnet in CIDR notation (e.g., 10.56.72.0/24)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validators.CidrValidator{},
				},
			},
			"tags": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Optional tags to associate with the VPC instance",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"vpc_name": schema.StringAttribute{
				Computed:    true,
				Description: "VPC name given when hosted at the cloud provider",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *vpcResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vpcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *vpcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tags := make([]string, 0)
	diag := plan.Tags.ElementsAs(ctx, &tags, false)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpc := model.VpcRequest{
		Name:   plan.Name.ValueString(),
		Region: plan.Region.ValueString(),
		Subnet: plan.Subnet.ValueString(),
		Tags:   tags,
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	data, err := r.client.CreateVPC(timeoutCtx, vpc)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create VPC Instance",
			fmt.Sprintf("Could not create VPC instance: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(data.ID))
	plan.VpcName = types.StringValue(data.VpcName)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not convert ID to integer: %s", err))
		return
	}

	data, err := r.client.ReadVPC(timeoutCtx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read VPC Instance",
			fmt.Sprintf("Could not read VPC instance with ID %d: %s", id, err),
		)
		return
	}

	state.Name = types.StringValue(data.Name)
	state.Region = types.StringValue(data.Region)
	state.Subnet = types.StringValue(data.Subnet)
	state.VpcName = types.StringValue(data.VpcName)

	if len(data.Tags) > 0 {
		tags, tagsDiag := types.ListValueFrom(ctx, types.StringType, data.Tags)
		resp.Diagnostics.Append(tagsDiag...)
		if resp.Diagnostics.HasError() {
			return
		}
		state.Tags = tags
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vpcResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var data model.VpcRequest
	data.Name = plan.Name.ValueString()
	data.Tags = make([]string, 0)
	resp.Diagnostics.Append(plan.Tags.ElementsAs(timeoutCtx, &data.Tags, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not convert ID to integer: %s", err))
		return
	}

	err = r.client.UpdateVPC(timeoutCtx, id, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update VPC instance",
			fmt.Sprintf("Could not update VPC instance with ID %d: %s", id, err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vpcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vpcResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	id, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not convert ID to integer: %s", err))
		return
	}

	err = r.client.DeleteVPC(timeoutCtx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete VPC Instance",
			fmt.Sprintf("Could not delete VPC instance with ID %d: %s", id, err),
		)
		return
	}
	resp.State.RemoveResource(ctx)
}
