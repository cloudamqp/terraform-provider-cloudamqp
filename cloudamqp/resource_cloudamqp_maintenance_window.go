package cloudamqp

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"
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
	_ resource.Resource                = &maintenanceWindowResource{}
	_ resource.ResourceWithConfigure   = &maintenanceWindowResource{}
	_ resource.ResourceWithImportState = &maintenanceWindowResource{}
)

type maintenanceWindowResource struct {
	client *api.API
}

func NewMaintenanceWindowResource() resource.Resource {
	return &maintenanceWindowResource{}
}

type maintenanceWindowResourceModel struct {
	ID               types.String `tfsdk:"id"`
	InstanceID       types.Int64  `tfsdk:"instance_id"`
	PreferredDay     types.String `tfsdk:"preferred_day"`
	PreferredTime    types.String `tfsdk:"preferred_time"`
	AutomaticUpdates types.String `tfsdk:"automatic_updates"`
}

func (r *maintenanceWindowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_maintenance_window"
}

func (r *maintenanceWindowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource ID of the maintenance window",
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
			"preferred_day": schema.StringAttribute{
				Optional:    true,
				Description: "Preferred day of the week when to run maintenance",
				Validators: []validator.String{
					stringvalidator.OneOf("Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"),
				},
			},
			"preferred_time": schema.StringAttribute{
				Optional:    true,
				Description: "Preferred time (UTC) the day when to run maintenance",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([0-1][0-9]|2[0-3]):([0-5][0-9])$`),
						"must be in format hh:mm (e.g., 14:30)",
					),
				},
			},
			"automatic_updates": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable automatic updates (only available for LavinMQ)",
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
			},
		},
	}
}

func (r *maintenanceWindowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *maintenanceWindowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *maintenanceWindowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan maintenanceWindowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	instanceID := int(plan.InstanceID.ValueInt64())
	data := model.Maintenance{}

	if !plan.PreferredDay.IsNull() && plan.PreferredDay.ValueString() != "" {
		data.PreferredDay = plan.PreferredDay.ValueString()
	}

	if !plan.PreferredTime.IsNull() && plan.PreferredTime.ValueString() != "" {
		data.PreferredTime = plan.PreferredTime.ValueString()
	}

	if !plan.AutomaticUpdates.IsNull() && plan.AutomaticUpdates.ValueString() != "" {
		if plan.AutomaticUpdates.ValueString() == "on" {
			data.AutomaticUpdates = utils.Pointer(true)
		} else if plan.AutomaticUpdates.ValueString() == "off" {
			data.AutomaticUpdates = utils.Pointer(false)
		}
	}

	err := r.client.SetMaintenance(timeoutCtx, instanceID, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Maintenance Window",
			fmt.Sprintf("Could not create maintenance window: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(strconv.Itoa(instanceID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *maintenanceWindowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state maintenanceWindowResourceModel
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

	// Set instance_id during import
	if state.InstanceID.IsNull() {
		state.InstanceID = types.Int64Value(int64(id))
	}

	data, err := r.client.ReadMaintenance(timeoutCtx, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Maintenance Window",
			fmt.Sprintf("Could not read maintenance window with ID %d: %s", id, err),
		)
		return
	}

	// Resource drift: resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("maintenance window not found, resource will be recreated: %s", state.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	state.PreferredDay = types.StringValue(data.PreferredDay)
	state.PreferredTime = types.StringValue(data.PreferredTime)

	if data.AutomaticUpdates != nil {
		if *data.AutomaticUpdates {
			state.AutomaticUpdates = types.StringValue("on")
		} else {
			state.AutomaticUpdates = types.StringValue("off")
		}
	} else {
		state.AutomaticUpdates = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *maintenanceWindowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan maintenanceWindowResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	instanceID := int(plan.InstanceID.ValueInt64())
	data := model.Maintenance{}

	if !plan.PreferredDay.IsNull() && plan.PreferredDay.ValueString() != "" {
		data.PreferredDay = plan.PreferredDay.ValueString()
	}

	if !plan.PreferredTime.IsNull() && plan.PreferredTime.ValueString() != "" {
		data.PreferredTime = plan.PreferredTime.ValueString()
	}

	if !plan.AutomaticUpdates.IsNull() && plan.AutomaticUpdates.ValueString() != "" {
		if plan.AutomaticUpdates.ValueString() == "on" {
			data.AutomaticUpdates = utils.Pointer(true)
		} else if plan.AutomaticUpdates.ValueString() == "off" {
			data.AutomaticUpdates = utils.Pointer(false)
		}
	}

	err := r.client.SetMaintenance(timeoutCtx, instanceID, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Maintenance Window",
			fmt.Sprintf("Could not update maintenance window: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *maintenanceWindowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Only remove from state because the maintenance window is managed by the API
	resp.State.RemoveResource(ctx)
}
