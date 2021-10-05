package cloudamqp

import (
	"fmt"
	"log"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomDomainCreate,
		Read:   resourceCustomDomainRead,
		Update: resourceCustomDomainUpdate,
		Delete: resourceCustomDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The custom hostname.",
			},
		},
	}
}

func resourceCustomDomainCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instanceID := d.Get("instance_id").(int)
	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::create instance id: %v", instanceID)
	hostname := d.Get("hostname").(string)
	data, err := api.CreateCustomDomain(instanceID, hostname)
	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::create data: %v", data)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instanceID))
	d.Set("instance_id", instanceID)

	for k, v := range data {
		if validateCustomDomainSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourceCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::read instance id: %v", instanceID)

	data, err := api.ReadCustomDomain(instanceID)

	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::read data: %v", data)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instanceID))
	d.Set("instance_id", instanceID)

	for k, v := range data {
		if validateCustomDomainSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourceCustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())
	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::update instance id: %v", instanceID)
	hostname := d.Get("hostname").(string)
	data, err := api.UpdateCustomDomain(instanceID, hostname)
	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::create data: %v", data)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instanceID))
	d.Set("instance_id", instanceID)

	for k, v := range data {
		if validateCustomDomainSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourceCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instanceID, _ := strconv.Atoi(d.Id())

	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::delete instance id: %v", instanceID)

	data, err := api.DeleteCustomDomain(instanceID)
	log.Printf("[DEBUG] cloudamqp::resource::custom_domain::delete data: %v", data)

	if err != nil {
		return err
	}

	return nil
}

func validateCustomDomainSchemaAttribute(key string) bool {
	switch key {
	case "hostname":
		return true
	}
	return false
}
