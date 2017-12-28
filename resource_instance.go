package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	d.SetId(address)
	return nil
}

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*MyClient)

	// Attempt to read from an upstream API
	obj, ok := client.Get(d.Id())

	// If the resource does not exist, inform Terraform. We want to immediately
	// return here to prevent further processing.
	if !ok {
		d.SetId("")
		return nil
	}

	d.Set("address", obj.Address)
	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")
	return nil
}
