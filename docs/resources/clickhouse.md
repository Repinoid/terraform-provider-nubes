---
page_title: "nubes_clickhouse"
description: "nubes_clickhouse TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_clickhouse` Resource

Service: ClickHouse.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/clickhouse/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `autoscale_configuration` | `object` | Required | No | — |
| `backup_configuration` | `object` | Required | No | — |
| `clickhouse_configuration` | `object` | Required | No | — |
| `clickhouse_files_configuration` | `string` | Optional, Computed | No | Формат:<br/> - path: somethings3.xml<br/>   xmlConfig: '&lt;clickhouse&gt;......&lt;/clickhouse&gt;' |
| `clickhouse_keeper_configuration` | `object` | Required | No | — |
| `cluster_configuration` | `object` | Required | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
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

### Nested attributes: `autoscale_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `enabled` | `bool` | Optional, Computed | No | Включает автоскейлинг PV. При включении параметра необходимо настроить autoScalePercentage и autoScaleTechWindow. Алерт при срабатывании вызывает modify у текущей услуги |
| `percent` | `number` | Optional, Computed | No | Требует включения autoScale. Позволяет указать в процентах расширение от текущего макс объема<br/>Минимальный размер на который расширяется - **1Gb** |
| `quota` | `number` | Optional, Computed | No | До какого размера возможно увеличивать диск, когда включен Autoscale |
| `schedule` | `number` | Optional, Computed | No | Требует включения autoScale. Окно обновления (пока не реализовано). Указывается час в виде Integer (0 / 5 / 12 / 23) |

### Nested attributes: `backup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `retain` | `number` | Optional, Computed | No | — |
| `s3_uid` | `string` | Required | No | Экземпляр S3 обязателен |
| `schedule` | `string` | Optional, Computed | No | — |

### Nested attributes: `clickhouse_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `version` | `string` | Optional, Computed | No | — |

### Nested attributes: `clickhouse_keeper_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `memory` | `number` | Optional, Computed | No | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | No | Если указан 0, то инстансы ClickHouse Keeper развёрнуты не будут |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | No | Указывается в GIGAbytes |
| `memory` | `number` | Optional, Computed | No | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | No | — |
| `shards` | `number` | Optional, Computed | No | **Шард** - это часть данных, разделенная на разные узлы (серверы) для горизонтального масштабирования и обработки больших объемов информации. <br/>**Реплика** - это точная копия данных внутри одного шарда |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `resource_realm` | `string` | Required | No | Кластер Kubernetes, на котором будет развернут экземпляр |
