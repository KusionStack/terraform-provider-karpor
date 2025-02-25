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

data "karpor_cluster" "example" {
  cluster_name = "local-cluster"
}

output "cluster_name" {
  value = data.karpor_cluster.example.cluster_name
}

output "display_name" {
  value = data.karpor_cluster.example.display_name
}

output "description" {
  value = data.karpor_cluster.example.description
}

