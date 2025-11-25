package cloudamqp

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/network"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &customCertificateResource{}
	_ resource.ResourceWithConfigure = &customCertificateResource{}
)

type customCertificateResource struct {
	client *api.API
}

func NewCustomCertificateResource() resource.Resource {
	return &customCertificateResource{}
}

type customCertificateResourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	InstanceID types.Int64  `tfsdk:"instance_id"`
	CA         types.String `tfsdk:"ca"`
	Cert       types.String `tfsdk:"cert"`
	PrivateKey types.String `tfsdk:"private_key"`
	SNIHosts   types.String `tfsdk:"sni_hosts"`
	Version    types.Int64  `tfsdk:"version"`
}

func (r *customCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "cloudamqp_custom_certificate"
}

func (r *customCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The identifier for this resource",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.Int64Attribute{
				Required:    true,
				Description: "The CloudAMQP instance identifier.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"ca": schema.StringAttribute{
				Required:    true,
				WriteOnly:   true,
				Description: "The PEM-encoded Certificate Authority (CA).",
			},
			"cert": schema.StringAttribute{
				Required:    true,
				WriteOnly:   true,
				Description: "The PEM-encoded certificate.",
			},
			"private_key": schema.StringAttribute{
				Required:    true,
				WriteOnly:   true,
				Description: "The PEM-encoded private key.",
			},
			"sni_hosts": schema.StringAttribute{
				Required:    true,
				Description: "A hostname (Server Name Indication) that this certificate applies to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
				Description: " An argument to trigger force new (default: 1)",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *customCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *customCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan customCertificateResourceModel
	var config customCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID := plan.InstanceID.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	request := model.CustomCertificateRequest{
		CA:         config.CA.ValueString(),
		Cert:       config.Cert.ValueString(),
		PrivateKey: config.PrivateKey.ValueString(),
		SNIHosts:   plan.SNIHosts.ValueString(),
	}

	jobResponse, err := r.client.CreateCustomCertificate(timeoutCtx, instanceID, request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating custom certificate", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *jobResponse.ID, 10*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for custom certificate job", err.Error())
		return
	}

	plan.ID = types.Int64Value(instanceID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *customCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Custom certificate resource doesn't support read operations
}

func (r *customCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Update is not supported - all changes require replacement
}

func (r *customCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state customCertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {

	}

	instanceID := state.InstanceID.ValueInt64()
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	jobResponse, err := r.client.DeleteCustomCertificate(timeoutCtx, instanceID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting custom certificate", err.Error())
		return
	}

	_, err = r.client.PollForJobCompleted(timeoutCtx, instanceID, *jobResponse.ID, 10*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Error polling for custom certificate deletion", err.Error())
		return
	}
}
