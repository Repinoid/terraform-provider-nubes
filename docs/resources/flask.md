---
page_title: "nubes_flask"
description: "nubes_flask TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_flask` Resource

Service: Web-сервер с фреймворком Flask.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/flask/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `app_configuration` | `object` | Required | No | — |
| `cluster_configuration` | `object` | Required | No | — |
| `git_revision` | `string` | Optional | No | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | No | — |
| `json_env` | `string` | Optional, Computed | No | — |
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
| `git_path` | `string` | Optional, Computed | No | Путь до репозитория приложения<br/>Обязательно должен быть передан как URL: `https://....`<br/><br/>Должен быть сетевой доступ до URL для скачивания исходников |
| `health_path` | `string` | Optional | No | Путь до healtcheck для запуска пода в Kubernetes, который отдаст 200 код<br/>Пустое значение означает, что будет проверяться исключительно доступность приложения по порту |
| `version` | `string` | Optional, Computed | No | — |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `memory` | `number` | Optional, Computed | No | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | No | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `resource_realm` | `string` | Required | No | Кластер Kubernetes, на котором будет развернут экземпляр |
