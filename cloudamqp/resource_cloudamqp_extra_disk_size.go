package cloudamqp

import (
	"fmt"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExtraDiskSize() *schema.Resource {
	return &schema.Resource{
		Create: resourceExtraDiskSizeUpdate,
		Read:   resourceExtraDiskSizeRead,
		Update: resourceExtraDiskSizeUpdate,
		Delete: resourceExtraDiskSizeDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "Instance identifier",
			},
			"extra_disk_size": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "Extra disk size in GB",
			},
			"allow_downtime": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When resizing disk, allow cluster downtime to do so",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Configurable sleep time in seconds between retries for resizing the disk",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Configurable timeout time in seconds for resizing the disk",
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subscription plan disk size",
						},
						"additional_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Additional added disk size",
						},
					},
				},
			},
		},
	}
}

func resourceExtraDiskSizeUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api     = meta.(*api.API)
		params  = make(map[string]interface{})
		sleep   = d.Get("sleep").(int)
		timeout = d.Get("timeout").(int)
	)

	params["extra_disk_size"] = d.Get("extra_disk_size")
	params["allow_downtime"] = d.Get("allow_downtime")

	_, err := api.ResizeDisk(d.Get("instance_id").(int), params, sleep, timeout)
	if err != nil {
		return err
	}
	id := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(id)

	return resourceExtraDiskSizeRead(d, meta)
}

func resourceExtraDiskSizeRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)

	data, err := api.ReadNodes(instanceID)
	if err != nil {
		return err
	}

	nodes := make([]map[string]interface{}, len(data))
	for k, v := range data {
		nodes[k] = readDiskNode(v)
	}

	if err = d.Set("nodes", nodes); err != nil {
		return fmt.Errorf("error setting nodes for resource %s, %s", d.Id(), err)
	}

	return nil
}

func resourceExtraDiskSizeDelete(d *schema.ResourceData, meta interface{}) error {
	// Just remove this resource from the state file, as the delete route does not exist in the backend
	// but we need to allow delete to happen, e.g. when you destroy your instance
	return nil
}

func readDiskNode(data map[string]interface{}) map[string]interface{} {
	node := make(map[string]interface{})
	for k, v := range data {
		if validateDiskSchemaAttribute(k) {
			node[k] = v
		}
	}
	return node
}

func validateDiskSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"disk_size",
		"additional_disk_size":
		return true
	}
	return false
}
