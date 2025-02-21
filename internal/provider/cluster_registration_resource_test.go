package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccClusterRegistration(t *testing.T) {
	// get root directory
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			// Make sure you have a valid kubeconfig file (only one cluster) in your home directory
			{
				Config: providerConfig + `
				resource "karpor_cluster_registration" "test" {
					cluster_name = "test-cluster"
					display_name = "test-display-name"
					credentials  =  file("~/config")
					description  = "test-description"
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("cluster_name"),
						knownvalue.StringExact("test-cluster"),
					),
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("display_name"),
						knownvalue.StringExact("test-display-name"),
					),
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test-description"),
					),
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("last_updated"),
						knownvalue.NotNull(),
					),
				},
			},
			// Update
			{
				Config: providerConfig + `
				resource "karpor_cluster_registration" "test" {
					cluster_name = "test-cluster"
					display_name = "test-display-name-updated"
					credentials  =  file("~/config")
					description  = "test-description-updated"
				}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("display_name"),
						knownvalue.StringExact("test-display-name-updated"),
					),
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test-description-updated"),
					),
					statecheck.ExpectKnownValue(
						"karpor_cluster_registration.test",
						tfjsonpath.New("last_updated"),
						knownvalue.NotNull(),
					),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
