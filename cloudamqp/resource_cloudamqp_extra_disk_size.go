package cloudamqp

import (
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceExtraDiskSize() *schema.Resource {
	return &schema.Resource{
		Create: resourceExtraDiskSizeUpdate,
		Read:   resourceExtraDiskSizeRead,
		Delete: resourceExtraDiskSizeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
		},
	}
}

func resourceExtraDiskSizeUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["extra_disk_size"] = d.Get("extra_disk_size")
	_, err := api.ResizeDisk(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}
	id := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(id)
	return nil
}

func resourceExtraDiskSizeRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceExtraDiskSizeDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
