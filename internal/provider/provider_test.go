package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "karpor" {
  api_endpoint = "https://127.0.0.1:7443"
  api_key      = "your-api-key-here"
  skip_tls_verify = true
}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"karpor": providerserver.NewProtocol6WithError(New("test")()),
	}
)
