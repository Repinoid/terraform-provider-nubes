---
page_title: "nubes_mongodb_rollback Resource"
description: "Service: MongoDB."
subcategory: "Базы данных"
---

# `nubes_mongodb_rollback` Resource

Service: MongoDB.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/mongodb/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/mongodb/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `backup_name` | `string` | Required | No | — |
| `id` | `string` | Computed | No | — |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mongodb_id` | `string` | Required | No | — |
| `operation_timeout` | `string` | Optional | No | — |
| `run_id` | `string` | Required | No | — |
