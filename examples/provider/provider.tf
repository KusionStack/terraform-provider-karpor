terraform {
  required_providers {
    karpor = {
      source  = "registry.terraform.io/KusionStack/karpor"
      version = "0.1.0"
    }
  }
}

provider "karpor" {
  api_endpoint = "https://api.karpor.example.com"
  api_key      = "your-api-key-here"
}
