# karpor_cluster_registration

Register and manage Kubernetes clusters in Karpor

## Example Usage

```hcl
resource "karpor_cluster_registration" "example" {
  cluster_name    = "production-cluster"
  api_server_url  = "https://kubernetes.example.com"
  credentials     = file("~/.kube/config")
  description     = "Production Kubernetes cluster"
}
```

## Argument Reference

- `cluster_name` (Required) - Unique name for the cluster
- `api_server_url` (Required) - Kubernetes API server URL
- `credentials` (Required) - Path to kubeconfig file
- `description` (Optional) - Human-readable description

## Attributes Reference

- `id` - Unique identifier for the registration
- `registration_time` - Timestamp of registration
- `health_status` - Current cluster health status
