---
page_title: "nubes_service_operation"
description: "Triggers an explicit one-time operation on an existing service instance."
subcategory: "Operations"
---

# `nubes_service_operation` Resource

Triggers an explicit one-time operation on an existing service instance.

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `id` | `string` | Computed | — |
| `instance_uid` | `string` | Required | — |
| `operation` | `string` | Required | — |
| `params` | `["map","string"]` | Optional | — |
| `run_id` | `string` | Required | — |
