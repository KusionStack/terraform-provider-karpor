package main

import (
	"context"
	"flag"
	"log"

	"github.com/KusionStack/terraform-provider-karpor/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	err := providerserver.Serve(
		context.Background(),
		provider.New(),
		providerserver.ServeOpts{
			Address: "registry.terraform.io/KusionStack/karpor",
			Debug:   debug,
		},
	)

	if err != nil {
		log.Fatal(err.Error())
	}
}
