terraform {
  required_providers {
    karpor = {
      source = "KusionStack/karpor"
      version = "0.1.0"
    }
  }
}

provider "karpor" {
  api_endpoint = "https://api.karpor.example.com"
  api_key      = "your-api-key-here"
}

resource "karpor_cluster_registration" "example" {
  cluster_name    = "production-cluster"
  api_server_url  = "https://kubernetes.example.com"
  credentials     = file("~/.kube/config")
  description     = "Production Kubernetes cluster"
}
