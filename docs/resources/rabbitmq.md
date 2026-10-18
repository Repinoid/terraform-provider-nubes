---
page_title: "nubes_rabbitmq"
description: "nubes_rabbitmq TEST resource schema for Nubes Cloud."
subcategory: "Очереди"
---

# `nubes_rabbitmq` Resource

Service: RabbitMQ.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/rabbitmq/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `cluster_configuration` | `object` | Required | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mq_configuration` | `object` | Required | No | — |
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

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | — |
| `disk` | `number` | Optional, Computed | No | — |
| `memory` | `number` | Optional, Computed | No | — |
| `replicas` | `number` | Optional, Computed | No | — |

### Nested attributes: `mq_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `max_message_size` | `number` | Optional, Computed | No | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `resource_realm` | `string` | Required | No | — |
