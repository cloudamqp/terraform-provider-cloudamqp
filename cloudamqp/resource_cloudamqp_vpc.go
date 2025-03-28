package cloudamqp

import (
	"context"
	"fmt"
	"net"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcCreate,
		ReadContext:   resourceVpcRead,
		UpdateContext: resourceVpcUpdate,
		DeleteContext: resourceVpcDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the VPC instance",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The hosted region for the standalone VPC instance",
			},
			"subnet": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val any, key string) (warns []string, errs []error) {
					v := val.(string)
					_, _, err := net.ParseCIDR(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("subnet: %v", err))
					}
					return
				},
				Description: "The VPC subnet",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Tag the VPC instance with optional tags",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC name given when hosted at the cloud provider",
			},
		},
	}
}

func resourceVpcCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	keys := []string{"name", "region", "subnet", "tags"}
	params := make(map[string]any)
	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			params[k] = v
		}
	}

	data, err := api.CreateVpcInstance(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(data["id"].(string))
	return resourceVpcRead(ctx, d, meta)
}

func resourceVpcRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	data, err := api.ReadVpcInstance(ctx, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		if validateVpcSchemaAttribute(k) {
			err = d.Set(k, v)
			if err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return diag.Diagnostics{}
}

func resourceVpcUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	keys := []string{"name", "tags"}
	params := make(map[string]any)
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = d.Get(k)
		}
	}

	if err := api.UpdateVpcInstance(ctx, d.Id(), params); err != nil {
		return diag.FromErr(err)
	}

	return resourceVpcRead(ctx, d, meta)
}

func resourceVpcDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	if err := api.DeleteVpcInstance(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validateVpcSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"region",
		"subnet",
		"tags",
		"vpc_name":
		return true
	}
	return false
}
