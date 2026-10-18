---
page_title: "nubes_nextcloud"
description: "nubes_nextcloud TEST resource schema for Nubes Cloud."
subcategory: "Хранилище"
---

# `nubes_nextcloud` Resource

Service: Nextcloud.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/nextcloud/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `app_configuration` | `object` | Required | No | — |
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
| `domain` | `string` | Optional, Computed | No | Если указан FQDN в формате: `myapp.example.com`, предварительно будет проверяться DNS-запись (А-запись **должна смотреть на IP Ingress кластера ресурсной платформы**)<br/><br/>Если указан в формате `myapp`, DNS-запись будет автоматически создана в служебной DNS-зоне |

### Nested attributes: `app_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `document_server` | `string` | Optional, Computed | No | Выбор сервера документов (onlyoffice)/(collabora). onlyoffice пока не доступен<br/>Размер пода collabora:<br/>1.8 ядра. 2Gi ОЗУ |
| `log_level` | `string` | Optional, Computed | No | Включает отображение логов самого nextcloud<br/>Логи уходят в /dev/stdout контейнера<br/>`NONE` отключает логирование |
| `version` | `string` | Optional, Computed | No | — |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | No | Указывается в GIGAbytes<br/>На данном диске хранятся данные Nextcloud |
| `memory` | `number` | Optional, Computed | No | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | No | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `psql_enabled` | `bool` | Optional, Computed | No | Инстанс обязан быть развёрнут на той же ресурсной платформе, где и будет располагаться Nextcloud |
| `psql_uid` | `string` | Required | No | Необходимо выбирать, когда выбран параметр `Подключить PostgreSQL из Личного кабинета` |
| `redis_enabled` | `bool` | Optional, Computed | No | Инстанс обязан быть развёрнут на той же ресурсной платформе, где и будет располагаться Nextcloud |
| `redis_uid` | `string` | Required | No | Необходимо выбирать, когда выбран параметр `Подключить Redis из Личного кабинета` |
| `resource_realm` | `string` | Required | No | Кластер Kubernetes, на котором будет развернут экземпляр |
