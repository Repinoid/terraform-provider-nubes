---
page_title: "nubes_dummy"
description: "nubes_dummy TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_dummy` Resource

Service: Болванка.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/dummy/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `arr` | `string` | Required | Не изменяется после создания. |
| `array_map_fixed_example` | `string` | Optional, Computed | Array of map-fixed Не изменяется после создания. |
| `bodymessage` | `string` | Optional, Computed | Дата, уходящая в Vault |
| `duration_ms` | `number` | Optional, Computed | перед завершением подождать (sleep) столько миллисекунд |
| `fail_at_start` | `bool` | Optional, Computed | если true, завершить джоб, ничего не записывая в cmdb |
| `fail_in_progress` | `bool` | Optional, Computed | если true, корректно начать выполнение джоба, но в процессе завершить аварийно, ничего не записывая в cmdb |
| `git_revision` | `string` | Optional | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | — |
| `json_example` | `string` | Optional, Computed | Для тестирования интерфейса |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `map_example` | `object` | Required | — |
| `map_fixed` | `object` | Required | — |
| `nested_ref_example` | `string` | Optional, Computed | Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | — |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Optional, Computed | Платформа для развертывания Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
| `where_fail` | `number` | Optional, Computed | принимает 1,2,3. Работает с failInProgress. 1 - упасть в prepare. 2 - упасть внутри блока Data fill. 3 - упасть после vault write |

### Nested attributes: `map_example`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `secret` | `string` | Optional | это секрет |
| `subparam1` | `string` | Optional | это для администратора |
| `subparam2` | `string` | Optional | — |

### Nested attributes: `map_fixed`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `bol1` | `bool` | Optional, Computed | Булевский, обязательный |
| `min_str1` | `string` | Optional, Computed | не меньше 10 символов |
| `param1` | `number` | Optional, Computed | Integer, secret |
| `param2` | `string` | Optional, Computed | Дескрипшен |
| `param3` | `number` | Optional, Computed | целое положительное |
