package cloudamqp

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/instance"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMaintenanceWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceMaintenanceWindowUpdate,
		Read:   resourceMaintenanceWindowRead,
		Update: resourceMaintenanceWindowUpdate,
		Delete: resourceMaintenanceWindowDelete,
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
						errs = append(errs, fmt.Errorf("the time: %s, in wrong on format hh:mm",
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
					if value.(string) == "" || value.(string) == "on" || value.(string) == "off" {
						return
					}
					errs = append(errs, fmt.Errorf("valid values are [on, off]"))
					return
				},
			},
		},
	}
}

func resourceMaintenanceWindowUpdate(d *schema.ResourceData, meta any) error {
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

	if err := api.SetMaintenance(instanceID, data); err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instanceID))
	return nil
}

func resourceMaintenanceWindowRead(d *schema.ResourceData, meta any) error {
	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id())
	)

	// Set argument during import
	if d.Get("instance_id").(int) == 0 {
		d.Set("instance_id", instanceID)
	}

	data, err := api.ReadMaintenance(instanceID)
	if err != nil {
		return err
	}

	d.Set("preferred_day", data.PreferredDay)
	d.Set("preferred_time", data.PreferredTime)

	if data.AutomaticUpdates != nil {
		d.Set("automatic_updates", "off")
		if *data.AutomaticUpdates {
			d.Set("automatic_updates", "on")
		}
	}

	return nil
}

func resourceMaintenanceWindowDelete(d *schema.ResourceData, meta any) error {
	// Only remove from state
	return nil
}
