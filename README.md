# Terraform Provider for Nubes Cloud

Terraform provider for managing Nubes Cloud resources.

## Install

```hcl
terraform {
  required_providers {
    nubes = {
      source  = "repinoid/nubes"
      version = "~> 1.5"
    }
  }
}
```

Configure the provider with the Nubes PROD API endpoint and an API token:

```hcl
provider "nubes" {
  api_endpoint = "https://lk-api-gateway.ngcloud.ru/api/v1/svc"
  api_token    = var.api_token
}
```

Keep API tokens and other credentials outside source control. See the
[provider documentation](https://repinoid.github.io/terraform-provider-nubes/)
and the [Terraform Registry](https://registry.terraform.io/providers/repinoid/nubes/latest).

This public package is an independently prepared PROD release. Existing Nubes
private-registry and tf-docs distributions remain separate and unchanged.
