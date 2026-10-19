---
page_title: "nubes_clickhouse"
description: "nubes_clickhouse TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_clickhouse` Resource

Service: ClickHouse.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/clickhouse/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_configuration` | `object` | Required | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `autoscale_configuration` | `object` | Required | — |
| `backup_configuration` | `object` | Required | — |
| `clickhouse_configuration` | `object` | Required | — |
| `clickhouse_files_configuration` | `string` | Optional, Computed | Формат:<br/> - path: somethings3.xml<br/>   xmlConfig: '&lt;clickhouse&gt;......&lt;/clickhouse&gt;' |
| `clickhouse_keeper_configuration` | `object` | Required | — |
| `cluster_configuration` | `object` | Required | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
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

### Nested attributes: `autoscale_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `enabled` | `bool` | Optional, Computed | Включает автоскейлинг PV. При включении параметра необходимо настроить autoScalePercentage и autoScaleTechWindow. Алерт при срабатывании вызывает modify у текущей услуги |
| `percent` | `number` | Optional, Computed | Требует включения autoScale. Позволяет указать в процентах расширение от текущего макс объема<br/>Минимальный размер на который расширяется - **1Gb** |
| `quota` | `number` | Optional, Computed | До какого размера возможно увеличивать диск, когда включен Autoscale |
| `schedule` | `number` | Optional, Computed | Требует включения autoScale. Окно обновления (пока не реализовано). Указывается час в виде Integer (0 / 5 / 12 / 23) |

### Nested attributes: `backup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `retain` | `number` | Optional, Computed | — |
| `s3_uid` | `string` | Required | Экземпляр S3 обязателен |
| `schedule` | `string` | Optional, Computed | — |

### Nested attributes: `clickhouse_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `version` | `string` | Optional, Computed | — |

### Nested attributes: `clickhouse_keeper_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `memory` | `number` | Optional, Computed | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | Если указан 0, то инстансы ClickHouse Keeper развёрнуты не будут |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | Указывается в GIGAbytes |
| `memory` | `number` | Optional, Computed | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | — |
| `shards` | `number` | Optional, Computed | **Шард** - это часть данных, разделенная на разные узлы (серверы) для горизонтального масштабирования и обработки больших объемов информации. <br/>**Реплика** - это точная копия данных внутри одного шарда |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `resource_realm` | `string` | Required | Кластер Kubernetes, на котором будет развернут экземпляр |
