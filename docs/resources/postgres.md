---
page_title: "nubes_postgres"
description: "nubes_postgres TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_postgres` Resource

Service: PostgreSQL.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/postgres/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_configuration` | `object` | Required | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `autoscale_configuration` | `object` | Required | — |
| `backup_configuration` | `object` | Required | — |
| `cluster_configuration` | `object` | Required | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mtls_configuration` | `object` | Optional | — |
| `operation_timeout` | `string` | Optional | — |
| `postgres_conf` | `object` | Required | — |
| `postgres_configuration` | `object` | Required | — |
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
| `slave_access_list` | `string` | Optional, Computed | Необходимо настраивать, когда зарезервирован внешний адрес<br/>Если передан пустой массив, доступ выделяется всем |
| `slave_ip_space` | `string` | Required | Из какого Ip-Space резервировать внешний IP<br/>Список предоставляется из Организации, в котором развёрнута ресурсная платформа |

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
| `increment_count` | `number` | Optional, Computed | — |
| `retain` | `number` | Optional, Computed | Подразумевается количество **полных бэкапов**<br/><br/>Простой пример:<br/>- Бэкап создаётся каждый день<br/>- Кол-во инкрементов - 6<br/>- Кол-во бэкапов - 3<br/><br/> |
| `s3_uid` | `string` | Optional, Computed | Экземпляр S3 обязателен |
| `schedule` | `string` | Optional, Computed | — |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | Указывается в GIGAbytes |
| `memory` | `number` | Optional, Computed | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | — |

### Nested attributes: `mtls_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cert_ca` | `string` | Optional | — |
| `cert_server` | `string` | Optional | — |
| `duration_ca` | `number` | Optional, Computed | Указывается в **часах**<br/>Необходимо настраивать, когда параметр `Режим настройки` - **auto** |
| `duration_server` | `number` | Optional, Computed | Указывается в **часах**<br/>Не может быть больше чем время жизни СА<br/>Обновляется **раз в год**, если выставлено значение **более 2 лет**<br/>Обновляется **раз в пол года**, если выставлено значение **менее 2 лет**<br/>Необходимо настраивать, когда параметр `Режим настройки` - **auto** |
| `refresh_cert` | `bool` | Optional, Computed | На случай если его нужно обновить/заменить |
| `type` | `string` | Optional, Computed | `off` - Настройка MTLS производиться не будет<br/>`manual` - После создания кластера в полях данных будет запрос на создание сертификата. **Вам необходимо будет выписать сертификат через вашу CA** и уже после включить через операцию **modify** режим MTLS<br/>`auto` - Сертификаты будут созданы через Cert-Manager<br/><br/>Режимы **auto** и **manual** требует включить настройку `Включить sslMode` - **true**<br/>MTLS в режиме **auto** будет работать сразу, в режиме **manual** только после того как будет предоставлен **CA** и **Сертификат сервера** |

### Nested attributes: `postgres_conf`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `log_connections` | `string` | Optional, Computed | — |
| `log_disconnections` | `string` | Optional, Computed | — |
| `max_connections` | `number` | Optional, Computed | — |

### Nested attributes: `postgres_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `pooler_master` | `bool` | Optional, Computed | Включает pgPooler (на запись). Будет создан новый под<br/>Размер пода **200m/400Mi** |
| `pooler_slave` | `bool` | Optional, Computed | Включает pgPooler (на чтение). Будет создан новый под<br/>Размер пода **200m/400Mi** |
| `ssl_required` | `bool` | Optional, Computed | По умолчанию подключение обязательное (sslmode=require) |
| `version` | `string` | Optional, Computed | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `resource_realm` | `string` | Required | Кластер Kubernetes, на котором будет развернут экземпляр |
