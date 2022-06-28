package cloudamqp

import (
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
	keys := []string{"extra_disk_size"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	_, err := api.ResizeDisk(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}

	d.SetId("NA")
	return nil
}

func resourceExtraDiskSizeRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// Need to fetch the extra disk size? Otherwise present in data source for nodes.

// func resourceExtraDiskSizeRead(d *schema.ResourceData, meta interface{}) error {
// 	// if strings.Contains(d.Id(), ",") {
// 	// 	s := strings.Split(d.Id(), ",")
// 	// 	d.SetId(s[0])
// 	// 	d.Set("name", s[0])
// 	// 	instanceID, _ := strconv.Atoi(s[1])
// 	// 	d.Set("instance_id", instanceID)
// 	// }
// 	// if d.Get("instance_id").(int) == 0 {
// 	// 	return errors.New("Missing instance identifier: {resource_id},{instance_id}")
// 	// }

// 	// api := meta.(*api.API)
// 	// data, err := api.ReadWebhook(d.Get("instance_id").(int), d.Id())

// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// for k, v := range data {
// 	// 	if validateWebhookSchemaAttribute(k) {
// 	// 		err = d.Set(k, v)

// 	// 		if err != nil {
// 	// 			return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
// 	// 		}
// 	// 	}
// 	// }

// 	return nil
// }

func resourceExtraDiskSizeDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
