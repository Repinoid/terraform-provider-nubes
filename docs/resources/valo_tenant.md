---
page_title: "nubes_valo_tenant"
description: "nubes_valo_tenant TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_valo_tenant` Resource

Service: VALO Cloud.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/valo_tenant/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `domain` | `string` | Required | No | Будет добавлено в зону valo.ngcloud.ru Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Не изменяется после создания. |
| `s3_storage` | `number` | Optional, Computed | No | Параметр указывается в `Gb` |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `valo_options` | `object` | Optional | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |

### Nested attributes: `valo_options`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `ff_logging` | `bool` | Optional, Computed | No | Тарифная опция |
| `ff_online_editor` | `bool` | Optional, Computed | No | Тарифная опция |
| `ff_public_link` | `bool` | Optional, Computed | No | Тарифная опция |
| `ff_themes` | `bool` | Optional, Computed | No | Тарифная опция |
| `ff_user_import` | `bool` | Optional, Computed | No | Тарифная опция |
| `ff_watermark` | `bool` | Optional, Computed | No | Тарифная опция |
| `request_approve_max_period` | `number` | Optional, Computed | No | Максимальное значение периода согласования (в днях) |
| `user_count` | `number` | Optional, Computed | No | — |
