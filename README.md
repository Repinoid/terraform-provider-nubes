# Terraform Provider for Nubes Cloud

Terraform provider for managing Nubes Cloud resources.

## Install

```hcl
terraform {
  required_providers {
    nubes = {
      source  = "repinoid/nubes"
      version = "~> 3.5"
    }
  }
}
```

Configure the provider with the Nubes TEST API endpoint and an API token:

```hcl
provider "nubes" {
  api_endpoint = "https://lk-api-gateway-test.ngcloud.ru/api/v1/svc"
  api_token    = var.api_token
}
```

Keep API tokens and other credentials outside source control. See the
[TEST S3 bucket documentation](https://registry.terraform.io/providers/Repinoid/nubes/3.5.0/docs/resources/s3bucket)
and [provider versions](https://registry.terraform.io/providers/Repinoid/nubes).

Version `3.x` targets TEST. Historical `1.x` versions target PROD.
