---
page_title: "nubes_http Resource"
description: "Service: Простой HTTP контейнер."
subcategory: "Приложения"
---

# `nubes_http` Resource

Service: Простой HTTP контейнер.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/http/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `domain` | `string` | Optional, Computed | No | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `git_revision` | `string` | Optional | No | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | No | — |
| `ingress_annotations` | `string` | Optional, Computed | No | Формат ключ/значение<br/>Поддерживает пока простые аннотации |
| `json_env` | `string` | Optional, Computed | No | Массив переменных, будут переданы в env приложения |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `registry_path` | `string` | Optional, Computed | No | Путь до image, который будет установлен как deployment |
| `resource_cpu` | `number` | Optional, Computed | No | Квота ядра пода (Указывается в mili) |
| `resource_instances` | `number` | Optional, Computed | No | Количество узлов |
| `resource_memory` | `number` | Optional, Computed | No | Квота памяти пода (Указывается в Megabytes) |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Платформа для развертывания Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
