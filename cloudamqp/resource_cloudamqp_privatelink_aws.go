package cloudamqp

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePrivateLinkAws() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateLinkAwsCreate,
		ReadContext:   resourcePrivateLinkAwsRead,
		UpdateContext: resourcePrivateLinkAwsUpdate,
		DeleteContext: resourcePrivateLinkAwsDelete,
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
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the PrivateLink [enabled, pending, disabled]",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service name of the PrivateLink, needed when creating the endpoint",
			},
			"allowed_principals": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Allowed principals to access the endpoint service",
			},
			"active_zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Covering availability zones used when creating an Endpoint from other VPC",
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
				Default:     1800,
				Description: "Configurable timeout in seconds when enable PrivateLink",
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateValue("allowed_principals", func(ctx context.Context, value, meta any) error {
				for _, v := range value.([]any) {
					re := regexp.MustCompile(`^arn:aws:iam::\d{12}:(root|user/.+|role/.+)$`)
					if !re.MatchString(v.(string)) {
						return fmt.Errorf("invalid ARN : %v", v)
					}
				}
				return nil
			}),
		),
	}
}

func resourcePrivateLinkAwsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		params     = make(map[string][]any)
	)

	params["allowed_principals"] = d.Get("allowed_principals").([]any)
	err := api.EnablePrivatelink(ctx, instanceID, params, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	return resourcePrivateLinkAwsRead(ctx, d, meta)
}

func resourcePrivateLinkAwsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id()) // Uses d.Id() to allow import
		sleep         = d.Get("sleep").(int)
		timeout       = d.Get("timeout").(int)
	)

	// Set arguments during import
	if d.Get("instance_id").(int) == 0 {
		d.Set("instance_id", instanceID)
	}
	if sleep == 0 && timeout == 0 {
		sleep = 10
		d.Set("sleep", 10)
		timeout = 1800
		d.Set("timeout", 1800)
	}

	data, err := api.ReadPrivatelink(ctx, instanceID, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("privatelink not found, resource will be recreated: %s", d.Id()))
		d.SetId("")
		return nil
	}

	for k, v := range data {
		if validatePrivateLinkAwsSchemaAttribute(k) {
			d.Set(k, v)
		}
	}
	return diag.Diagnostics{}
}

func resourcePrivateLinkAwsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		params     = make(map[string][]any)
	)

	params["allowed_principals"] = d.Get("allowed_principals").([]any)
	err := api.UpdatePrivatelink(ctx, instanceID, params)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func resourcePrivateLinkAwsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)

	err := api.DisablePrivatelink(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validatePrivateLinkAwsSchemaAttribute(key string) bool {
	switch key {
	case "status",
		"service_name",
		"allowed_principals",
		"active_zones":
		return true
	}
	return false
}
