package cloudamqp

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomDomainCreate,
		ReadContext:   resourceCustomDomainRead,
		UpdateContext: resourceCustomDomainUpdate,
		DeleteContext: resourceCustomDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The custom hostname.",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Configurable sleep time in seconds between retries for custom domain configuration",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Configurable timeout time in seconds for custom domain configuration",
			},
		},
	}
}

func resourceCustomDomainCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID := d.Get("instance_id").(int)
	hostname := d.Get("hostname").(string)
	sleep := d.Get("sleep").(int)
	timeout := d.Get("timeout").(int)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	data, err := api.CreateCustomDomain(timeoutCtx, instanceID, hostname, time.Duration(sleep)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(instanceID))
	d.Set("instance_id", instanceID)

	for k, v := range data {
		if validateCustomDomainSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return diag.Diagnostics{}
}

func resourceCustomDomainRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	sleep := d.Get("sleep").(int)
	timeout := d.Get("timeout").(int)

	// Set defaults during import
	if d.Get("instance_id").(int) == 0 {
		d.Set("instance_id", instanceID)
	}
	if sleep == 0 && timeout == 0 {
		sleep = 10
		d.Set("sleep", 10)
		timeout = 1800
		d.Set("timeout", 1800)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	data, err := api.ReadCustomDomain(timeoutCtx, instanceID, time.Duration(sleep)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(instanceID))
	d.Set("instance_id", instanceID)

	for k, v := range data {
		if validateCustomDomainSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return diag.Diagnostics{}
}

func resourceCustomDomainUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	hostname := d.Get("hostname").(string)
	sleep := d.Get("sleep").(int)
	timeout := d.Get("timeout").(int)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	data, err := api.UpdateCustomDomain(timeoutCtx, instanceID, hostname, time.Duration(sleep)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(instanceID))
	d.Set("instance_id", instanceID)

	for k, v := range data {
		if validateCustomDomainSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return diag.Diagnostics{}
}

func resourceCustomDomainDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	sleep := d.Get("sleep").(int)
	timeout := d.Get("timeout").(int)

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	_, err := api.DeleteCustomDomain(timeoutCtx, instanceID, time.Duration(sleep)*time.Second)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func validateCustomDomainSchemaAttribute(key string) bool {
	switch key {
	case "hostname":
		return true
	}
	return false
}
