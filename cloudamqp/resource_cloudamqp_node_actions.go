package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &nodeActionsResource{}
	_ resource.ResourceWithConfigure = &nodeActionsResource{}
)

func NewNodeActionsResource() resource.Resource {
	return &nodeActionsResource{}
}

type nodeActionsResource struct {
	client *api.API
}

type nodeActionsResourceModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	NodeName   types.String `tfsdk:"node_name"`
	NodeNames  types.List   `tfsdk:"node_names"`
	Action     types.String `tfsdk:"action"`
	Sleep      types.Int64  `tfsdk:"sleep"`
	Timeout    types.Int64  `tfsdk:"timeout"`
}

func (r *nodeActionsResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

func (r *nodeActionsResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "cloudamqp_node_actions"
}

func (r *nodeActionsResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Description: "Perform actions on CloudAMQP nodes, such as start, stop, restart, or reboot. " +
			"Actions can be performed on individual nodes or multiple nodes at once. " +
			"Cluster-level actions (cluster.start, cluster.stop, cluster.restart) affect all nodes in the cluster.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Resource identifier",
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "Instance identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"node_name": schema.StringAttribute{
				Optional:           true,
				Description:        "The name of the node (deprecated: use node_names instead)",
				DeprecationMessage: "Use node_names instead. This attribute will be removed in a future version.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("node_names")),
				},
			},
			"node_names": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of node names to perform the action on. For cluster-level actions (cluster.start, cluster.stop, cluster.restart), this can be omitted or should include all nodes.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.ConflictsWith(path.MatchRoot("node_name")),
				},
			},
			"action": schema.StringAttribute{
				Required:    true,
				Description: "The action to perform.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"start",
						"stop",
						"restart",
						"reboot",
						"mgmt.restart",
						"cluster.start",
						"cluster.stop",
						"cluster.restart",
					),
				},
			},
			"sleep": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Sleep interval in seconds between polling for node status (default: 10)",
				Default:     int64default.StaticInt64(10),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Timeout in seconds for the action to complete (default: 1800)",
				Default:     int64default.StaticInt64(1800),
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
		},
	}
}

func (r *nodeActionsResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan nodeActionsResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Determine which nodes to act on
	var nodeNames []string
	// Handle backward compatibility with node_name
	if !plan.NodeName.IsNull() && !plan.NodeName.IsUnknown() {
		nodeNames = []string{plan.NodeName.ValueString()}
	} else if !plan.NodeNames.IsNull() && !plan.NodeNames.IsUnknown() {
		diag := plan.NodeNames.ElementsAs(ctx, &nodeNames, false)
		response.Diagnostics.Append(diag...)
		if response.Diagnostics.HasError() {
			return
		}
	}

	// For cluster actions, if no nodes specified, get all nodes
	action := plan.Action.ValueString()
	isClusterAction := action == "cluster.start" || action == "cluster.stop" || action == "cluster.restart"

	if isClusterAction && len(nodeNames) == 0 {
		nodes, err := r.client.ListNodes(ctx, plan.InstanceID.ValueInt64())
		if err != nil {
			response.Diagnostics.AddError(
				"Failed to List Nodes",
				fmt.Sprintf("Could not list nodes for cluster action: %s", err),
			)
			return
		}
		for _, node := range nodes {
			nodeNames = append(nodeNames, node.Name)
		}
	}

	// Validate that we have at least one node
	if len(nodeNames) == 0 {
		response.Diagnostics.AddError(
			"No Nodes Specified",
			"Either node_name or node_names must be specified for non-cluster actions, or nodes must exist for cluster actions",
		)
		return
	}

	// Create timeout context
	sleep := time.Duration(plan.Sleep.ValueInt64()) * time.Second
	timeout := time.Duration(plan.Timeout.ValueInt64()) * time.Second
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Perform the action
	if err := r.client.PostAction(timeoutCtx, plan.InstanceID.ValueInt64(), nodeNames, action, sleep); err != nil {
		response.Diagnostics.AddError(
			"Node Action Failed",
			fmt.Sprintf("Failed to perform action '%s' on nodes: %s", action, err),
		)
		return
	}

	// Set ID to instance ID
	plan.ID = types.StringValue(plan.InstanceID.String())
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *nodeActionsResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	// This resource does not implement the Read function
	// Actions are one-time operations and don't have persistent state to read
}

func (r *nodeActionsResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// This resource does not implement the Update function
	// All attributes have RequiresReplace, so any change will trigger recreation
}

func (r *nodeActionsResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// This resource does not implement the Delete function
	// Actions are one-time operations with no cleanup needed
	response.State.RemoveResource(ctx)
}
