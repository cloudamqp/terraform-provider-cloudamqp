package main

import (
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return cloudamqp.Provider()
		},
	})
}
