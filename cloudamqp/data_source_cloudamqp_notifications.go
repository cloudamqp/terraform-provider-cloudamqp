package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &notificationsDataSource{}
	_ datasource.DataSourceWithConfigure = &notificationsDataSource{}
)

type notificationsDataSource struct {
	client *api.API
}

func NewNotificationsDataSource() datasource.DataSource {
	return &notificationsDataSource{}
}

type notificationsDataSourceModel struct {
	ID         types.String                            `tfsdk:"id"`
	InstanceID types.Int64                             `tfsdk:"instance_id"`
	Recipients []recipientNotificationsDataSourceModel `tfsdk:"recipients"`
}

type recipientNotificationsDataSourceModel struct {
	RecipientID types.Int64                              `tfsdk:"recipient_id"`
	Type        types.String                             `tfsdk:"type"`
	Value       types.String                             `tfsdk:"value"`
	Name        types.String                             `tfsdk:"name"`
	Options     types.Map                                `tfsdk:"options"`
	Responders  *[]notificationsDataSourceResponderModel `tfsdk:"responders"`
}

type notificationsDataSourceResponderModel struct {
	Type     types.String `tfsdk:"type"`
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Username types.String `tfsdk:"username"`
}

func (d *notificationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "cloudamqp_notifications"
}

func (d *notificationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about default or created notifications. Either use" +
			" recipient_id or name to retrieve the notification.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The identifier for this data source",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
			},
		},
		Blocks: map[string]schema.Block{
			"recipients": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"recipient_id": schema.Int64Attribute{
							Optional:    true,
							Computed:    true,
							Description: "Recipient identifier",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the notification.",
						},
						"value": schema.StringAttribute{
							Computed:    true,
							Description: "Notification endpoint, where to send the notifcation",
						},
						"name": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Optional display name of the recipient",
						},
						"options": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
							Description: "Optional key-value pair options parameters (e.g. dedupkey, rk)",
						},
					},
					Blocks: map[string]schema.Block{
						"responders": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Computed:    true,
										Description: "Responder type, valid options are: team, user, escalation, schedule",
									},
									"id": schema.StringAttribute{
										Computed:    true,
										Description: "Responder ID",
									},
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "Responder name",
									},
									"username": schema.StringAttribute{
										Computed:    true,
										Description: "Responder username",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *notificationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *notificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config notificationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := config.InstanceID.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	data, err := d.client.ListNotifications(timeoutCtx, instanceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Notification",
			fmt.Sprintf("Could not read notification: %s", err),
		)
		return
	}

	if data == nil {
		tflog.Warn(ctx, fmt.Sprintf("Resource drift detected for instance ID %d", instanceID))
		resp.State.RemoveResource(ctx)
		return
	}

	config.ID = types.StringValue(fmt.Sprintf("%d.notifications", instanceID))
	config.InstanceID = types.Int64Value(instanceID)

	for _, notification := range data {
		dataSourceModel := d.populateDataSourceModel(notification)
		config.Recipients = append(config.Recipients, dataSourceModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func (d *notificationsDataSource) populateDataSourceModel(data model.RecipientResponse) recipientNotificationsDataSourceModel {
	config := recipientNotificationsDataSourceModel{
		RecipientID: types.Int64Value(data.ID),
		Type:        types.StringValue(data.Type),
		Value:       types.StringValue(data.Value),
		Name:        types.StringValue(data.Name),
		Options:     types.MapNull(types.StringType),
		Responders:  nil,
	}

	switch data.Type {
	case "opsgenie", "opsgenie-eu":
		if data.Options != nil && data.Options.Responders != nil && len(*data.Options.Responders) > 0 {
			responderModels := make([]notificationsDataSourceResponderModel, len(*data.Options.Responders))
			for i, responder := range *data.Options.Responders {
				responderModel := notificationsDataSourceResponderModel{
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
			config.Responders = &responderModels
			config.Options = types.MapNull(types.StringType)
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
			config.Options = types.MapValueMust(types.StringType, opts)
		}
		config.Responders = nil
	}
	return config
}
