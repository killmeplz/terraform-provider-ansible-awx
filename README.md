# Terraform Provider for Ansible AWX

This Terraform provider allows you to manage AWX instances using Terraform.

## Configuration

Example provider configuration:

```hcl
provider "ansible_awx" {
  host  = "https://awx.example.com"
  token = "your-api-token"
}
```