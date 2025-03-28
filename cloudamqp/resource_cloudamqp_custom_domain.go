package cloudamqp

import (
	"context"
	"strconv"

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
		},
	}
}

func resourceCustomDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID := d.Get("instance_id").(int)
	hostname := d.Get("hostname").(string)
	data, err := api.CreateCustomDomain(ctx, instanceID, hostname)
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

func resourceCustomDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())

	data, err := api.ReadCustomDomain(ctx, instanceID)
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

func resourceCustomDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	hostname := d.Get("hostname").(string)

	data, err := api.UpdateCustomDomain(ctx, instanceID, hostname)
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

func resourceCustomDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())

	_, err := api.DeleteCustomDomain(ctx, instanceID)
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
