package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDAMQP_APIKEY", nil),
				Description: "The API key used to connect to CloudAMQP",
			},
			"base_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDAMQP_BASE_URL", "https://customer.cloudamqp.com/api/"),
				Description: "The CloudAMQP Base API URL",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudamqp_instance": resourceInstance(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey:  d.Get("apikey").(string),
		BaseURL: d.Get("base_url").(string),
	}

	return config.Client()
}
