---
page_title: "nubes_lucee"
description: "nubes_lucee TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_lucee` Resource

Service: Lucee.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/lucee/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_configuration` | `object` | Required | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `app_configuration` | `object` | Required | — |
| `cluster_configuration` | `object` | Required | — |
| `git_revision` | `string` | Optional | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | — |
| `json_env` | `string` | Optional, Computed | Передаётся в формате КЛЮЧ:ЗНАЧЕНИЕ |
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
| `git_path` | `string` | Optional, Computed | Путь до репозитория приложения<br/>Обязательно должен быть передан как URL: `https://....`<br/><br/>Должен быть сетевой доступ до URL для скачивания исходников |
| `health_path` | `string` | Optional | Путь до healtcheck для запуска пода в Kubernetes, который отдаст 200 код<br/>Пустое значение означает, что будет проверяться исключительно доступность приложения по порту |
| `version` | `string` | Optional, Computed | — |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `memory` | `number` | Optional, Computed | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `resource_realm` | `string` | Required | Кластер Kubernetes, на котором будет развернут экземпляр |
