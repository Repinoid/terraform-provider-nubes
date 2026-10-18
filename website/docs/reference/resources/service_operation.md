# `nubes_service_operation` Resource

Use this resource to execute a supported one-time operation for an existing service instance.
Set a new `run_id` when a subsequent operation should execute.

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `id` | `string` | Computed | No | — |
| `instance_uid` | `string` | Required | No | — |
| `operation` | `string` | Required | No | — |
| `params` | `["map","string"]` | Optional | No | — |
| `run_id` | `string` | Required | No | — |
