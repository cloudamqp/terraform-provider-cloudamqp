package main

import (
	"context"
	"log"
	"net/http"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

var version string

func main() {
	ctx := context.Background()

	muxServer, err := tf5muxserver.NewMuxServer(
		ctx, cloudamqp.Provider(version, http.DefaultClient).GRPCProvider,
		providerserver.NewProtocol5(cloudamqp.New(version, http.DefaultClient)),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = tf5server.Serve(
		"registry.terraform.io/cloudamqp/cloudamqp",
		muxServer.ProviderServer,
	)

	if err != nil {
		log.Fatal(err)
	}
}
