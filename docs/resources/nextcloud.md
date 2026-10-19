---
page_title: "nubes_nextcloud"
description: "nubes_nextcloud TEST resource schema for Nubes Cloud."
subcategory: "Хранилище"
---

# `nubes_nextcloud` Resource

Service: Nextcloud.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/nextcloud/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_configuration` | `object` | Required | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `app_configuration` | `object` | Required | — |
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
| `domain` | `string` | Optional, Computed | Если указан FQDN в формате: `myapp.example.com`, предварительно будет проверяться DNS-запись (А-запись **должна смотреть на IP Ingress кластера ресурсной платформы**)<br/><br/>Если указан в формате `myapp`, DNS-запись будет автоматически создана в служебной DNS-зоне |

### Nested attributes: `app_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `document_server` | `string` | Optional, Computed | Выбор сервера документов (onlyoffice)/(collabora). onlyoffice пока не доступен<br/>Размер пода collabora:<br/>1.8 ядра. 2Gi ОЗУ |
| `log_level` | `string` | Optional, Computed | Включает отображение логов самого nextcloud<br/>Логи уходят в /dev/stdout контейнера<br/>`NONE` отключает логирование |
| `version` | `string` | Optional, Computed | — |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | Указывается в GIGAbytes<br/>На данном диске хранятся данные Nextcloud |
| `memory` | `number` | Optional, Computed | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `psql_enabled` | `bool` | Optional, Computed | Инстанс обязан быть развёрнут на той же ресурсной платформе, где и будет располагаться Nextcloud |
| `psql_uid` | `string` | Required | Необходимо выбирать, когда выбран параметр `Подключить PostgreSQL из Личного кабинета` |
| `redis_enabled` | `bool` | Optional, Computed | Инстанс обязан быть развёрнут на той же ресурсной платформе, где и будет располагаться Nextcloud |
| `redis_uid` | `string` | Required | Необходимо выбирать, когда выбран параметр `Подключить Redis из Личного кабинета` |
| `resource_realm` | `string` | Required | Кластер Kubernetes, на котором будет развернут экземпляр |
