package cloudamqp

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMaintenanceWindowUpdate,
		ReadContext:   resourceMaintenanceWindowRead,
		UpdateContext: resourceMaintenanceWindowUpdate,
		DeleteContext: resourceMaintenanceWindowDelete,
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
			"preferred_day": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preferred day of the week when to run maintenance",
				ValidateFunc: func(value any, key string) (warns []string, errs []error) {
					if value.(string) == "" {
						return
					}
					days := []string{"Monday", "Tuesday", "Wednesday", "Thursday",
						"Friday", "Saturday", "Sunday"}
					if !slices.Contains(days, value.(string)) {
						errs = append(errs, fmt.Errorf("should be a day of the week: %v", days))
					}
					return
				},
			},
			"preferred_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preferred time (UTC) the day when to run maintenance",
				ValidateFunc: func(value any, key string) (warns []string, errs []error) {
					if value.(string) == "" {
						return
					}
					r, _ := regexp.Compile("^([0-1][0-9]|2[0-3]):([0-5][0-9])$")
					if !r.MatchString(value.(string)) {
						errs = append(errs, fmt.Errorf("the time: %s, is in the wrong format hh:mm",
							value.(string)))
					}
					return
				},
			},
			"automatic_updates": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enable automatic updates",
				ValidateFunc: func(value any, key string) (warns []string, errs []error) {
					val := strings.ToLower(value.(string))
					if val == "" || val == "on" || val == "off" {
						return
					}
					errs = append(errs, fmt.Errorf("valid values are [on, off]"))
					return
				},
			},
		},
	}
}

func resourceMaintenanceWindowUpdate(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		data       model.Maintenance
	)

	if v := d.Get("preferred_day"); v != nil {
		data.PreferredDay = v.(string)
	}

	if v := d.Get("preferred_time"); v != nil {
		data.PreferredTime = v.(string)
	}

	if v := d.Get("automatic_updates"); v != nil {
		if v.(string) == "on" {
			data.AutomaticUpdates = utils.Pointer(true)
		} else if v.(string) == "off" {
			data.AutomaticUpdates = utils.Pointer(false)
		}
	}

	if d.HasChanges("preferred_day", "preferred_time", "automatic_updates") {
		if err := api.SetMaintenance(ctx, instanceID, data); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(instanceID))
	return diag.Diagnostics{}
}

func resourceMaintenanceWindowRead(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id())
	)
	// Set argument during import
	if instanceIDValue, ok := d.Get("instance_id").(int); ok && instanceIDValue == 0 {
		d.Set("instance_id", instanceID)
	}

	data, err := api.ReadMaintenance(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("maintenance window not found, resource will be recreated: %s", d.Id()))
		d.SetId("")
		return nil
	}

	d.Set("preferred_day", data.PreferredDay)
	d.Set("preferred_time", data.PreferredTime)

	if data.AutomaticUpdates != nil {
		if *data.AutomaticUpdates {
			d.Set("automatic_updates", "on")
		} else {
			d.Set("automatic_updates", "off")
		}
	}

	return diag.Diagnostics{}
}

func resourceMaintenanceWindowDelete(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	// Only remove from state because the maintenance window is managed by the API
	return diag.Diagnostics{}
}
