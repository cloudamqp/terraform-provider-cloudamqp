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

func resourcePrivateLinkAzure() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateLinkAzureCreate,
		ReadContext:   resourcePrivateLinkAzureRead,
		UpdateContext: resourcePrivateLinkAzureUpdate,
		DeleteContext: resourcePrivateLinkAzureDelete,
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
				Description: "Service name (alias) of the PrivateLink, needed when creating the endpoint",
			},
			"server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the server having the PrivateLink enabled",
			},
			"approved_subscriptions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Approved subscriptions to access the endpoint service",
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
			customdiff.ValidateValue("approved_subscriptions", func(ctx context.Context, value, meta any) error {
				for _, v := range value.([]any) {
					re := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
					if !re.MatchString(v.(string)) {
						return fmt.Errorf("invalid Subscription ID : %v", v)
					}
				}
				return nil
			}),
		),
	}
}

func resourcePrivateLinkAzureCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		params     = make(map[string][]any)
	)

	params["approved_subscriptions"] = d.Get("approved_subscriptions").([]any)
	err := api.EnablePrivatelink(ctx, instanceID, params, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	return resourcePrivateLinkAzureRead(ctx, d, meta)
}

func resourcePrivateLinkAzureRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
		if validatePrivateLinkAzureSchemaAttribute(k) {
			if k == "alias" {
				d.Set("service_name", v)
			} else {
				d.Set(k, v)
			}
		}
	}
	return diag.Diagnostics{}
}

func resourcePrivateLinkAzureUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		params     = make(map[string][]any)
	)

	params["approved_subscriptions"] = d.Get("approved_subscriptions").([]any)
	err := api.UpdatePrivatelink(ctx, instanceID, params)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func resourcePrivateLinkAzureDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func validatePrivateLinkAzureSchemaAttribute(key string) bool {
	switch key {
	case "status",
		"server_name",
		"alias",
		"approved_subscriptions":
		return true
	}
	return false
}
