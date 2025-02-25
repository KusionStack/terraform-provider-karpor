terraform {
  required_providers {
    karpor = {
      source  = "registry.terraform.io/KusionStack/karpor"
      version = "0.1.0"
    }
  }
}

provider "karpor" {
  api_endpoint    = "https://127.0.0.1:7443"
  api_key         = "your-api-key-here"
  skip_tls_verify = true
}

resource "karpor_cluster_registration" "example" {
  cluster_name = "local-cluster"
  display_name = "local-cluster-display-name"
  credentials  = file("~/config")
  description  = "local-cluster-description"
}


# make sure you have a existing demo cluster in karpor
# id is the cluster name
import {
  to = karpor_cluster_registration.imported
  id = "demo"
}

resource "karpor_cluster_registration" "imported" {
  cluster_name = "demo"
  # Setting credentials in imported resource will destroy and then create replacement
  # credentials  = file("~/config")
  description  = "demo-description"
  display_name = "demo-display-name"
}

output "cluster_name" {
  value = karpor_cluster_registration.example.cluster_name
}

output "display_name" {
  value = karpor_cluster_registration.example.display_name
}
