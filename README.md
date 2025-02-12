![OSP](https://socialify.git.ci/KusionStack/terraform-provider-karpor/image?font=Raleway&language=1&name=1&owner=1&pattern=Plus&theme=Light)

# Terraform Provider for Karpor

A Terraform provider for managing cluster registration in Karpor

## Features

- Cluster Registration Management (`karpor_cluster_registration`)

## Installation

### Local Build & Install
```bash
make install  # Install to ~/.terraform.d/plugins
```

### Terraform Configuration
```hcl
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
  api_key      = "<your-api-key>"  # Recommend using environment variables
}
```

## Usage Example
```hcl
resource "karpor_cluster_registration" "production" {
  cluster_name    = "production-cluster"
  api_server_url  = "https://k8s.example.com"
  credentials     = file("~/.kube/config")
  description     = "Primary production cluster"
}
```

## Development Guide

### Requirements
- Go 1.21+
- Terraform 1.5+

### Common Commands
```bash
make build    # Build provider
make test     # Run unit tests
make testacc  # Run acceptance tests (requires API credentials)
```

### Test Configuration
Set environment variables before testing:
```bash
export KARPOR_ENDPOINT="https://api.karpor.example.com"
export KARPOR_API_KEY="your-api-key"
```

## Contributing
1. Create an issue describing the problem or feature request
2. Develop on a feature branch (feature/xxx)
3. Submit a PR with test cases
4. Pass CI pipeline verification

## License
Mozilla Public License 2.0, see [LICENSE](LICENSE)
