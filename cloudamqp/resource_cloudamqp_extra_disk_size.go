package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	models "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &extraDiskSizeResource{}
	_ resource.ResourceWithConfigure = &extraDiskSizeResource{}
)

type extraDiskSizeResource struct {
	client *api.API
}

func NewExtraDiskSizeResource() resource.Resource {
	return &extraDiskSizeResource{}
}

type extraDiskSizeResourceModel struct {
	ID            types.String `tfsdk:"id"`
	InstanceID    types.Int64  `tfsdk:"instance_id"`
	ExtraDiskSize types.Int64  `tfsdk:"extra_disk_size"`
	AllowDowntime types.Bool   `tfsdk:"allow_downtime"`
	Sleep         types.Int64  `tfsdk:"sleep"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	Nodes         types.List   `tfsdk:"nodes"`
}

var extraDiskNodeObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":                 types.StringType,
		"disk_size":            types.Int64Type,
		"additional_disk_size": types.Int64Type,
	},
}

func (r *extraDiskSizeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_extra_disk_size"
}

func (r *extraDiskSizeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resize additional disk for an instance and wait until all nodes are configured.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier (same as instance_id)",
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
			"extra_disk_size": schema.Int64Attribute{
				Required:    true,
				Description: "Extra disk size in GB. Valid values: 0, 25, 50, 100, 250, 500, 1000, 2000",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.OneOf(0, 25, 50, 100, 250, 500, 1000, 2000),
				},
			},
			"allow_downtime": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "When resizing disk, allow cluster downtime to do so",
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(30),
				Description: "Sleep time in seconds between retries for disk resize polling",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:           true,
				Computed:           true,
				Default:            int64default.StaticInt64(1800),
				Description:        "Timeout in seconds for disk resize and node convergence",
				DeprecationMessage: "The timeout attribute is deprecated and will be removed in next 2.0.0 version. Use the timeouts block instead.",
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"nodes": schema.ListAttribute{
				Computed:    true,
				Description: "Nodes and resulting disk sizes",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":                 types.StringType,
						"disk_size":            types.Int64Type,
						"additional_disk_size": types.Int64Type,
					},
				},
			},
		},
	}
}

func (r *extraDiskSizeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *extraDiskSizeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan extraDiskSizeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Hour)
	defer cancel()

	if err := r.applyResize(timeoutCtx, &plan); err != nil {
		resp.Diagnostics.AddError("Error Resizing Extra Disk", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *extraDiskSizeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state extraDiskSizeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data, err := r.client.ListNodes(ctx, state.InstanceID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Extra Disk State", err.Error())
		return
	}

	nodeValues := make([]attr.Value, len(data))
	for i, node := range data {
		nodeValues[i] = types.ObjectValueMust(
			extraDiskNodeObjectType.AttrTypes,
			map[string]attr.Value{
				"name":                 types.StringValue(node.Name),
				"disk_size":            types.Int64Value(node.DiskSize),
				"additional_disk_size": types.Int64Value(node.AdditionalDiskSize),
			},
		)
	}

	state.Nodes = types.ListValueMust(
		extraDiskNodeObjectType,
		nodeValues,
	)
	state.ID = types.StringValue(state.InstanceID.String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *extraDiskSizeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan extraDiskSizeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Hour)
	defer cancel()

	if err := r.applyResize(timeoutCtx, &plan); err != nil {
		resp.Diagnostics.AddError("Error Updating Extra Disk", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *extraDiskSizeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No delete route exists in backend; this resource is state-only on destroy.
	resp.State.RemoveResource(ctx)
}

func (r *extraDiskSizeResource) applyResize(ctx context.Context, plan *extraDiskSizeResourceModel) error {
	params := models.ExtraDiskRequest{
		ExtraDisk:     plan.ExtraDiskSize.ValueInt64(),
		AllowDowntime: plan.AllowDowntime.ValueBool(),
	}

	_, err := r.client.ResizeDisk(ctx, plan.InstanceID.ValueInt64(), params, plan.Sleep.ValueInt64())
	if err != nil {
		return fmt.Errorf("could not resize extra disk for instance %d: %w", plan.InstanceID.ValueInt64(), err)
	}

	plan.ID = types.StringValue(plan.InstanceID.String())

	data, err := r.client.ListNodes(ctx, plan.InstanceID.ValueInt64())
	if err != nil {
		return fmt.Errorf("could not read nodes after resize for instance %d: %w", plan.InstanceID.ValueInt64(), err)
	}

	nodeValues := make([]attr.Value, len(data))
	for i, node := range data {
		nodeValues[i] = types.ObjectValueMust(
			extraDiskNodeObjectType.AttrTypes,
			map[string]attr.Value{
				"name":                 types.StringValue(node.Name),
				"disk_size":            types.Int64Value(node.DiskSize),
				"additional_disk_size": types.Int64Value(node.AdditionalDiskSize),
			},
		)
	}
	plan.Nodes = types.ListValueMust(
		extraDiskNodeObjectType,
		nodeValues,
	)

	return nil
}
