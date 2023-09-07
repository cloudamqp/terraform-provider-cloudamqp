package cloudamqp

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNodeAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceNodeActionRequest,
		Update: resourceNodeActionRequest,
		Read:   resourceNodeActionRead,
		Delete: resourceNodeActionRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"node_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the node",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The action to perform on the node",
				ValidateFunc: validateNodeAction(),
			},
			"running": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the node is running",
			},
		},
	}
}

func resourceNodeActionRequest(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data := make(map[string]interface{})
	nodeName := d.Get("node_name").(string)
	data, err := api.PostAction(d.Get("instance_id").(int), nodeName, d.Get("action").(string))
	if err != nil {
		return err
	}
	d.SetId(nodeName)
	d.Set("running", data["running"])
	return nil
}

func resourceNodeActionRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data := make(map[string]interface{})
	data, err := api.ReadNode(d.Get("instance_id").(int), d.Get("node_name").(string))
	if err != nil {
		return err
	}
	d.Set("running", data["running"])
	return nil
}

func resourceNodeActionRemove(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func validateNodeAction() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"start",
		"stop",
		"restart",
		"reboot",
		"mgmt.restart",
	}, true)
}
