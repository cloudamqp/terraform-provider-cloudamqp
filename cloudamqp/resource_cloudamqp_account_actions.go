package cloudamqp

import (
	"errors"
	"fmt"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAccountAction() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountActionRequest,
		Update: resourceAccountActionRequest,
		Read:   resourceAccountActionRead,
		Delete: resourceAccountActionRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The action to perform on the node",
				ValidateFunc: validateAccountAction(),
			},
		},
	}
}

func resourceAccountActionRequest(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		action     = d.Get("action")
		err        = errors.New("")
	)

	switch action {
	case "rotate-password":
		err = api.RotatePassword(instanceID)
	case "rotate-apikey":
		err = api.RotateApiKey(instanceID)
	}
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	return nil
}

func resourceAccountActionRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAccountActionRemove(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func validateAccountAction() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"rotate-password",
		"rotate-apikey",
	}, true)
}
