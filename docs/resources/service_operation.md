---
page_title: "nubes_service_operation"
description: "Triggers an explicit one-time operation on an existing service instance."
subcategory: "Operations"
---

# `nubes_service_operation` Resource

Triggers an explicit one-time operation on an existing service instance.

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `id` | `string` | Computed | No | — |
| `instance_uid` | `string` | Required | No | — |
| `operation` | `string` | Required | No | — |
| `params` | `["map","string"]` | Optional | No | — |
| `run_id` | `string` | Required | No | — |
