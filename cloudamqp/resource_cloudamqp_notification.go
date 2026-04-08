package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &notificationResource{}
	_ resource.ResourceWithConfigure   = &notificationResource{}
	_ resource.ResourceWithImportState = &notificationResource{}
)

type notificationResource struct {
	client *api.API
}

func NewNotificationResource() resource.Resource {
	return &notificationResource{}
}

type notificationResourceModel struct {
	ID         types.String                          `tfsdk:"id"`
	InstanceID types.Int64                           `tfsdk:"instance_id"`
	Type       types.String                          `tfsdk:"type"`
	Value      types.String                          `tfsdk:"value"`
	Name       types.String                          `tfsdk:"name"`
	Options    types.Map                             `tfsdk:"options"`
	Responders *[]notificationResourceResponderModel `tfsdk:"responders"`
}

type notificationResourceResponderModel struct {
	Type     types.String `tfsdk:"type"`
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
}

func (r *notificationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_notification"
}

func (r *notificationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "This resource allows you to create and manage notification endpoints (recipients) " +
			"to be used together with alarms.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "The instance identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Description: "Type of the notification, valid options are: email, opsgenie, opsgenie-eu, " +
					"pagerduty, signl4, slack, teams, victorops, webhook",
				Validators: []validator.String{
					stringvalidator.OneOfCaseInsensitive(
						"email",
						"opsgenie",
						"opsgenie-eu",
						"pagerduty",
						"signl4",
						"slack",
						"teams",
						"victorops",
						"webhook",
					),
				},
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: "Notification endpoint, where to send the notification",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional display name of the recipient",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"options": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional key-value pair options parameters (e.g. dedupkey, rk)",
			},
		},
		Blocks: map[string]schema.Block{
			"responders": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Responder type, valid options are: team, user, escalation, schedule",
							Validators: []validator.String{
								stringvalidator.OneOfCaseInsensitive(
									"escalation",
									"schedule",
									"team",
									"user",
								),
							},
						},
						"id": schema.StringAttribute{
							Optional:    true,
							Description: "Responder ID",
						},
						"name": schema.StringAttribute{
							Optional:    true,
							Description: "Responder name",
						},
						"username": schema.StringAttribute{
							Optional:    true,
							Description: "Responder username",
						},
					},
				},
			},
		},
	}
}

func (r *notificationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *notificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("ImportState: ID=%s", req.ID))
	if !strings.Contains(req.ID, ",") {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {recipient_id},{instance_id}")
		return
	}

	idSplit := strings.Split(req.ID, ",")
	if len(idSplit) != 2 {
		resp.Diagnostics.AddError("Invalid import ID format", "Expected format: {recipient_id},{instance_id}")
		return
	}

	instanceID, err := strconv.ParseInt(idSplit[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance_id in import ID", fmt.Sprintf("Could not convert instance_id to int: %s", err))
		return
	}

	resp.State.SetAttribute(ctx, path.Root("id"), idSplit[0])
	resp.State.SetAttribute(ctx, path.Root("instance_id"), instanceID)
}

func (r *notificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan notificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := r.populateRequest(ctx, plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	data, err := r.client.CreateNotification(timeoutCtx, plan.InstanceID.ValueInt64(), &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Notification",
			fmt.Sprintf("Could not create notification: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *notificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state notificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	recipientID := state.ID.ValueString()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	data, err := r.client.ReadNotification(timeoutCtx, instanceID, recipientID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Notification",
			fmt.Sprintf("Could not read notification: %s", err),
		)
		return
	}

	if data == nil {
		tflog.Warn(ctx, fmt.Sprintf("Resource drift detected for notification ID %s or instance ID %d", recipientID, instanceID))
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateResourceModel(*data, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *notificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan notificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := r.populateRequest(ctx, plan)
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	err := r.client.UpdateNotification(timeoutCtx, plan.InstanceID.ValueInt64(), plan.ID.ValueString(), &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Notification",
			fmt.Sprintf("Could not update notification: %s", err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *notificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state notificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := state.InstanceID.ValueInt64()
	recipientID := state.ID.ValueString()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	err := r.client.DeleteNotification(timeoutCtx, instanceID, recipientID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Notification",
			fmt.Sprintf("Could not delete notification: %s", err),
		)
		return
	}
}

func (r *notificationResource) populateRequest(ctx context.Context, plan notificationResourceModel) model.RecipientRequest {
	params := model.RecipientRequest{
		Type:  plan.Type.ValueString(),
		Value: plan.Value.ValueString(),
		Name:  plan.Name.ValueString(),
	}

	switch plan.Type.ValueString() {
	case "opsgenie", "opsgenie-eu":
		if plan.Responders != nil && len(*plan.Responders) > 0 {
			list := make([]model.RecipientResponder, len(*plan.Responders))
			for i, responder := range *plan.Responders {
				list[i] = model.RecipientResponder{
					Type:     responder.Type.ValueString(),
					ID:       responder.ID.ValueStringPointer(),
					Name:     responder.Name.ValueStringPointer(),
					Username: responder.Username.ValueStringPointer(),
				}
			}
			params.Options = &model.RecipientOptions{
				Responders: &list,
			}
		}
	case "pagerduty", "victorops":
		if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
			var opts map[string]string
			plan.Options.ElementsAs(ctx, &opts, false)
			if len(opts) > 0 {
				options := &model.RecipientOptions{}
				if v, ok := opts["dedupkey"]; ok {
					options.DedupKey = &v
				}
				if v, ok := opts["rk"]; ok {
					options.RK = &v
				}
				params.Options = options
			}
		}
	}

	return params
}

func (r *notificationResource) populateResourceModel(data model.RecipientResponse, state *notificationResourceModel) {
	state.Type = types.StringValue(data.Type)
	state.Value = types.StringValue(data.Value)
	state.Name = types.StringValue(data.Name)
	state.Options = types.MapNull(types.StringType)
	state.Responders = nil

	switch data.Type {
	case "opsgenie", "opsgenie-eu":
		if data.Options != nil && data.Options.Responders != nil && len(*data.Options.Responders) > 0 {
			responderModels := make([]notificationResourceResponderModel, len(*data.Options.Responders))
			for i, responder := range *data.Options.Responders {
				responderModel := notificationResourceResponderModel{
					Type: types.StringValue(responder.Type),
				}
				if responder.ID != nil {
					responderModel.ID = types.StringValue(*responder.ID)
				} else {
					responderModel.ID = types.StringNull()
				}
				if responder.Name != nil {
					responderModel.Name = types.StringValue(*responder.Name)
				} else {
					responderModel.Name = types.StringNull()
				}
				if responder.Username != nil {
					responderModel.Username = types.StringValue(*responder.Username)
				} else {
					responderModel.Username = types.StringNull()
				}
				responderModels[i] = responderModel
			}
			state.Responders = &responderModels
		}
	case "pagerduty", "victorops":
		if data.Options != nil && (data.Options.DedupKey != nil || data.Options.RK != nil) {
			opts := map[string]attr.Value{}
			if data.Options.DedupKey != nil {
				opts["dedupkey"] = types.StringValue(*data.Options.DedupKey)
			}
			if data.Options.RK != nil {
				opts["rk"] = types.StringValue(*data.Options.RK)
			}
			state.Options = types.MapValueMust(types.StringType, opts)
		}
	}
}
