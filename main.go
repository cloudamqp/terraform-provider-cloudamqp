package main

import (
	"net/http"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

var version string

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return cloudamqp.Provider(version, http.DefaultClient)
		},
	})
}
