package cloudamqp

import (
	"log"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUpgradeLavinMQ() *schema.Resource {
	return &schema.Resource{
		Create: resourceUpgradeLavinMQInvoke,
		Read:   resourceUpgradeLavinMQRead,
		Update: resourceUpgradeLavinMQUpdate,
		Delete: resourceUpgradeLavinMQRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"new_version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The new version to upgrade to",
			},
		},
	}
}

func resourceUpgradeLavinMQInvoke(d *schema.ResourceData, meta interface{}) error {
	var (
		api         = meta.(*api.API)
		instanceID  = d.Get("instance_id").(int)
		new_version = d.Get("new_version").(string)
	)

	log.Printf("[DEBUG] - Upgrading LavinMQ instance %d to version %s", instanceID, new_version)
	response, err := api.UpgradeLavinMQ(instanceID, new_version)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instanceID))

	if len(response) > 0 {
		log.Println("[INFO] - ", response)
	}

	return nil
}

func resourceUpgradeLavinMQRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUpgradeLavinMQUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUpgradeLavinMQRemove(d *schema.ResourceData, meta interface{}) error {
	return nil
}
