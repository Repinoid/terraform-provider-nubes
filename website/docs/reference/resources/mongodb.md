# `nubes_mongodb` Resource

Service: MongoDB.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/mongodb/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/mongodb/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `backup_configuration` | `object` | Required | No | — |
| `cluster_configuration` | `object` | Required | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mongo_configuration` | `object` | Required | No | — |
| `operation_timeout` | `string` | Optional | No | — |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | No | — |
| `startup_configuration` | `object` | Required | No | — |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `master_access_list` | `string` | Optional, Computed | No | Необходимо настраивать, когда зарезервирован внешний адрес<br/>Если передан пустой массив, доступ выделяется всем |
| `master_ip_space` | `string` | Required | No | Из какого Ip-Space резервировать внешний IP<br/>Список предоставляется из Организации, в котором развёрнута ресурсная платформа |

### Nested attributes: `backup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `retain` | `number` | Optional, Computed | No | — |
| `s3_uid` | `string` | Required | No | — |
| `schedule` | `string` | Optional, Computed | No | ``` * * * * * ```  &lt;-- Минута(0-59), Час(0-23), День(1-31), Месяц(1-12), День_недели(0-7, где 0-Вс) |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | — |
| `disk` | `number` | Optional, Computed | No | Указывается в Gb |
| `memory` | `number` | Optional, Computed | No | — |
| `replicas` | `number` | Optional, Computed | No | — |

### Nested attributes: `mongo_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `tls` | `bool` | Required | No | Шифрует данные между внешними подключениями и базой данных. |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `resource_realm` | `string` | Optional, Computed | No | Кластер Kubernetes, на котором будет развернут экземпляр |
