---
page_title: "nubes_postgres Resource"
description: "Service: PostgreSQL."
subcategory: "Базы данных"
---

# `nubes_postgres` Resource

Service: PostgreSQL.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/postgres/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/postgres/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `autoscale_configuration` | `object` | Required | No | — |
| `backup_configuration` | `object` | Required | No | — |
| `cluster_configuration` | `object` | Required | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mtls_configuration` | `object` | Optional | No | — |
| `operation_timeout` | `string` | Optional | No | — |
| `postgres_conf` | `object` | Required | No | — |
| `postgres_configuration` | `object` | Required | No | — |
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
| `slave_access_list` | `string` | Optional, Computed | No | Необходимо настраивать, когда зарезервирован внешний адрес<br/>Если передан пустой массив, доступ выделяется всем |
| `slave_ip_space` | `string` | Required | No | Из какого Ip-Space резервировать внешний IP<br/>Список предоставляется из Организации, в котором развёрнута ресурсная платформа |

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
| `increment_count` | `number` | Optional, Computed | No | — |
| `retain` | `number` | Optional, Computed | No | Подразумевается количество **полных бэкапов**<br/><br/>Простой пример:<br/>- Бэкап создаётся каждый день<br/>- Кол-во инкрементов - 6<br/>- Кол-во бэкапов - 3<br/><br/> |
| `s3_uid` | `string` | Optional, Computed | No | Экземпляр S3 обязателен |
| `schedule` | `string` | Optional, Computed | No | — |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | No | Указывается в GIGAbytes |
| `memory` | `number` | Optional, Computed | No | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | No | — |

### Nested attributes: `mtls_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cert_ca` | `string` | Optional | No | — |
| `cert_server` | `string` | Optional | No | — |
| `duration_ca` | `number` | Optional, Computed | No | Указывается в **часах**<br/>Необходимо настраивать, когда параметр `Режим настройки` - **auto** |
| `duration_server` | `number` | Optional, Computed | No | Указывается в **часах**<br/>Не может быть больше чем время жизни СА<br/>Обновляется **раз в год**, если выставлено значение **более 2 лет**<br/>Обновляется **раз в пол года**, если выставлено значение **менее 2 лет**<br/>Необходимо настраивать, когда параметр `Режим настройки` - **auto** |
| `refresh_cert` | `bool` | Optional, Computed | No | На случай если его нужно обновить/заменить |
| `type` | `string` | Optional, Computed | No | `off` - Настройка MTLS производиться не будет<br/>`manual` - После создания кластера в полях данных будет запрос на создание сертификата. **Вам необходимо будет выписать сертификат через вашу CA** и уже после включить через операцию **modify** режим MTLS<br/>`auto` - Сертификаты будут созданы через Cert-Manager<br/><br/>Режимы **auto** и **manual** требует включить настройку `Включить sslMode` - **true**<br/>MTLS в режиме **auto** будет работать сразу, в режиме **manual** только после того как будет предоставлен **CA** и **Сертификат сервера** |

### Nested attributes: `postgres_conf`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `log_connections` | `string` | Optional, Computed | No | — |
| `log_disconnections` | `string` | Optional, Computed | No | — |
| `max_connections` | `number` | Optional, Computed | No | — |

### Nested attributes: `postgres_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `pooler_master` | `bool` | Optional, Computed | No | Включает pgPooler (на запись). Будет создан новый под<br/>Размер пода **200m/400Mi** |
| `pooler_slave` | `bool` | Optional, Computed | No | Включает pgPooler (на чтение). Будет создан новый под<br/>Размер пода **200m/400Mi** |
| `ssl_required` | `bool` | Optional, Computed | No | По умолчанию подключение обязательное (sslmode=require) |
| `version` | `string` | Optional, Computed | No | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `resource_realm` | `string` | Optional, Computed | No | Кластер Kubernetes, на котором будет развернут экземпляр |
