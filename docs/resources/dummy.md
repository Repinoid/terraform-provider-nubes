---
page_title: "nubes_dummy Resource"
description: "Service: Болванка."
subcategory: "Приложения"
---

# `nubes_dummy` Resource

Service: Болванка.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/dummy/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/dummy/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `arr` | `string` | Required | No | Не изменяется после создания. |
| `array_map_fixed_example` | `string` | Optional, Computed | No | Array of map-fixed Не изменяется после создания. |
| `bodymessage` | `string` | Optional, Computed | No | Дата, уходящая в Vault |
| `duration_ms` | `number` | Optional, Computed | No | перед завершением подождать (sleep) столько миллисекунд |
| `fail_at_start` | `bool` | Optional, Computed | No | если true, завершить джоб, ничего не записывая в cmdb |
| `fail_in_progress` | `bool` | Optional, Computed | No | если true, корректно начать выполнение джоба, но в процессе завершить аварийно, ничего не записывая в cmdb |
| `git_revision` | `string` | Optional | No | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | No | — |
| `json_example` | `string` | Optional, Computed | No | Для тестирования интерфейса |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `map_example` | `object` | Required | No | — |
| `map_fixed` | `object` | Required | No | — |
| `nested_ref_example` | `string` | Optional, Computed | No | Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | No | — |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Required | No | Платформа для развертывания Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
| `where_fail` | `number` | Optional, Computed | No | принимает 1,2,3. Работает с failInProgress. 1 - упасть в prepare. 2 - упасть внутри блока Data fill. 3 - упасть после vault write |

### Nested attributes: `map_example`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `secret` | `string` | Optional | No | это секрет |
| `subparam1` | `string` | Optional | No | это для администратора |
| `subparam2` | `string` | Optional | No | — |

### Nested attributes: `map_fixed`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `bol1` | `bool` | Optional, Computed | No | Булевский, обязательный |
| `min_str1` | `string` | Optional, Computed | No | не меньше 10 символов |
| `param1` | `number` | Optional, Computed | No | Integer, secret |
| `param2` | `string` | Optional, Computed | No | Дескрипшен |
| `param3` | `number` | Optional, Computed | No | целое положительное |
