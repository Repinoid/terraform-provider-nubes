---
page_title: "nubes_mongodb"
description: "nubes_mongodb TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_mongodb` Resource

Service: MongoDB.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/mongodb/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_configuration` | `object` | Required | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `backup_configuration` | `object` | Required | — |
| `cluster_configuration` | `object` | Required | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mongo_configuration` | `object` | Required | — |
| `operation_timeout` | `string` | Optional | — |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | — |
| `startup_configuration` | `object` | Required | — |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `master_access_list` | `string` | Optional, Computed | Необходимо настраивать, когда зарезервирован внешний адрес<br/>Если передан пустой массив, доступ выделяется всем |
| `master_ip_space` | `string` | Required | Из какого Ip-Space резервировать внешний IP<br/>Список предоставляется из Организации, в котором развёрнута ресурсная платформа |

### Nested attributes: `backup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `retain` | `number` | Optional, Computed | — |
| `s3_uid` | `string` | Required | — |
| `schedule` | `string` | Optional, Computed | ``` * * * * * ```  &lt;-- Минута(0-59), Час(0-23), День(1-31), Месяц(1-12), День_недели(0-7, где 0-Вс) |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | — |
| `disk` | `number` | Optional, Computed | Указывается в Gb |
| `memory` | `number` | Optional, Computed | — |
| `replicas` | `number` | Optional, Computed | — |

### Nested attributes: `mongo_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `tls` | `bool` | Required | Шифрует данные между внешними подключениями и базой данных. |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `resource_realm` | `string` | Required | Кластер Kubernetes, на котором будет развернут экземпляр |
