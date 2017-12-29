package main

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/waveaccounting/go-cloudamqp/cloudamqp"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance",
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the plan, valid options are: lemur, tiger, bunny, rabbit, panda, ape, hippo, lion",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the region you want to create your instance in",
			},
			"vpc_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dedicated VPC subnet, shouldn't overlap with your current VPC's subnet",
			},
			"nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Number of nodes in cluster (plan must support it)",
			},
			"rmq_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "RabbitMQ version",
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudamqp.Client)
	fmt.Print(client)

	params := &cloudamqp.CreateInstanceParams{
		Name:   d.Get("name").(string),
		Plan:   d.Get("plan").(string),
		Region: d.Get("region").(string),
	}

	instance, _, err := client.Instances.Create(params)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instance.ID))
	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudamqp.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	instance, _, err := client.Instances.Get(id)
	if err != nil {
		// If the resource does not exist, inform Terraform. We want to immediately
		// return here to prevent further processing.
		d.SetId("")
		return nil
	}

	d.SetId(strconv.Itoa(instance.ID))
	d.Set("name", instance.Name)
	d.Set("region", instance.Region)
	d.Set("plan", instance.Plan)
	d.Set("url", instance.URL)
	d.Set("apikey", instance.APIKey)

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudamqp.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Instances.Delete(id)
	return err
}
