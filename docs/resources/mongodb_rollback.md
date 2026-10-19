---
page_title: "nubes_mongodb_rollback"
description: "Runs a MongoDB rollback operation for an existing MongoDB instance."
subcategory: "Databases"
---

# `nubes_mongodb_rollback` Resource

Runs a MongoDB rollback operation for an existing MongoDB instance.

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `backup_name` | `string` | Required | — |
| `id` | `string` | Computed | — |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mongodb_id` | `string` | Required | — |
| `operation_timeout` | `string` | Optional | — |
| `run_id` | `string` | Required | — |
