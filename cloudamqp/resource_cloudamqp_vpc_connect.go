package cloudamqp

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	AWS_ARN_VALIDATE_RE, _        = regexp.Compile(`\Aarn:aws:iam::\d{12}:(root|((user|role)(/[^/]+)?))\z`)
	AZURE_SUBS_VALIDATE_RE, _     = regexp.Compile(`\A[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\z`)
	GCP_PROJECT_ID_VALIDATE_RE, _ = regexp.Compile(`\A[a-z][0-9a-z-]{4,28}[0-9a-z]\z`)
)

func resourceVpcConnect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcConnectCreate,
		ReadContext:   resourceVpcConnectRead,
		UpdateContext: resourceVpcConnectUpdate,
		DeleteContext: resourceVpcConnectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region where the CloudAMQP instance is hosted",
			},
			"allowed_principals": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of allowed prinicpals used by AWS",
			},
			"approved_subscriptions": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of approved subscriptions used by Azure",
			},
			"allowed_projects": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of allowed projects used by GCP",
			},
			"active_zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Covering availability zones used when creating an endpoint from other VPC. [AWS]",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Private Service Connect [enabled, pending, disabled]",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service name (alias for Azure) of the PrivateLink.",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Configurable sleep in seconds between retries when enable PrivateLink",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Configurable timeout in seconds when enable PrivateLink",
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateValue("allowed_principals", func(ctx context.Context,
				value, meta any) error {

				for _, v := range value.([]any) {
					if AWS_ARN_VALIDATE_RE.MatchString(v.(string)) {
						continue
					} else {
						return fmt.Errorf("invalid format for 'allowed_principals': %s", v.(string))
					}
				}
				return nil
			}),
			customdiff.ValidateValue("approved_subscriptions", func(ctx context.Context, value,
				meta any) error {

				for _, v := range value.([]any) {
					if AZURE_SUBS_VALIDATE_RE.MatchString(v.(string)) {
						continue
					} else {
						return fmt.Errorf("invalid format for 'approved_subscriptions': %s", v.(string))
					}
				}
				return nil
			}),
			customdiff.ValidateValue("allowed_projects", func(ctx context.Context, value,
				meta any) error {

				for _, v := range value.([]any) {
					if GCP_PROJECT_ID_VALIDATE_RE.MatchString(v.(string)) {
						continue
					} else {
						return fmt.Errorf("invalid format for 'allowed_projets': %s", v.(string))
					}
				}
				return nil
			}),
		),
	}
}

func resourceVpcConnectCreate(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		region     = d.Get("region").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		params     = make(map[string][]any)
	)

	switch getPlatform(region) {
	case "amazon":
		params["allowed_principals"] = d.Get("allowed_principals").([]any)
	case "azure":
		params["approved_subscriptions"] = d.Get("approved_subscriptions").([]any)
	case "google":
		params["allowed_projects"] = d.Get("allowed_projects").([]any)
	default:
		return diag.Errorf("invalid region")
	}

	err := api.EnableVpcConnect(ctx, instanceID, params, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	return resourceVpcConnectRead(ctx, d, meta)
}

func resourceVpcConnectRead(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id()) // Uses d.Id() to allow import
	)

	// Set arguments during import
	if d.Get("instance_id").(int) == 0 {
		d.Set("instance_id", instanceID)
	}
	if d.Get("sleep").(int) == 0 && d.Get("timeout").(int) == 0 {
		d.Set("sleep", 10)
		d.Set("timeout", 3600)
	}

	data, err := api.ReadVpcConnect(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("vpc connect not found, resource will be recreated: %s", d.Id()))
		d.SetId("")
		return nil
	}

	d.Set("active_zones", []string{})
	for k, v := range data {
		if validateVpcConnectSchemaAttribute(k) {
			if k == "alias" {
				d.Set("service_name", v)
			} else {
				d.Set(k, v)
			}
		}
	}
	return diag.Diagnostics{}
}

func resourceVpcConnectUpdate(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		region     = d.Get("region").(string)
		params     = make(map[string][]any)
	)

	switch getPlatform(region) {
	case "amazon":
		params["allowed_principals"] = d.Get("allowed_principals").([]any)
	case "azure":
		params["approved_subscriptions"] = d.Get("approved_subscriptions").([]any)
	case "google":
		params["allowed_projects"] = d.Get("allowed_projects").([]any)
	default:
		return diag.Errorf("invalid region")
	}

	err := api.UpdateVpcConnect(ctx, instanceID, params)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func resourceVpcConnectDelete(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)

	err := api.DisableVpcConnect(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validateVpcConnectSchemaAttribute(key string) bool {
	switch key {
	case "active_zones",
		"alias",
		"allowed_principals",
		"allowed_projects",
		"approved_subscriptions",
		"service_name",
		"status":
		return true
	}
	return false
}

func getPlatform(region string) string {
	regionSplit := strings.Split(region, "::")
	switch regionSplit[0] {
	case "amazon-web-services":
		return "amazon"
	case "azure-arm":
		return "azure"
	case "google-compute-engine":
		return "google"
	default:
		return ""
	}
}
