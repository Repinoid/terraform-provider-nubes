---
page_title: "nubes_mongodb_rollback"
description: "Runs a MongoDB rollback operation for an existing MongoDB instance."
subcategory: "Databases"
---

# `nubes_mongodb_rollback` Resource

Runs a MongoDB rollback operation for an existing MongoDB instance.

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
