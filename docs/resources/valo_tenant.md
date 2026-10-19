---
page_title: "nubes_valo_tenant"
description: "nubes_valo_tenant TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_valo_tenant` Resource

Service: VALO Cloud.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/valo_tenant/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `domain` | `string` | Required | Будет добавлено в зону valo.ngcloud.ru Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Optional, Computed | Не изменяется после создания. |
| `s3_storage` | `number` | Optional, Computed | Параметр указывается в `Gb` |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `valo_options` | `object` | Optional | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |

### Nested attributes: `valo_options`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `ff_logging` | `bool` | Optional, Computed | Тарифная опция |
| `ff_online_editor` | `bool` | Optional, Computed | Тарифная опция |
| `ff_public_link` | `bool` | Optional, Computed | Тарифная опция |
| `ff_themes` | `bool` | Optional, Computed | Тарифная опция |
| `ff_user_import` | `bool` | Optional, Computed | Тарифная опция |
| `ff_watermark` | `bool` | Optional, Computed | Тарифная опция |
| `request_approve_max_period` | `number` | Optional, Computed | Максимальное значение периода согласования (в днях) |
| `user_count` | `number` | Optional, Computed | — |
