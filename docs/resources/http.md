---
page_title: "nubes_http"
description: "nubes_http TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_http` Resource

Service: Простой HTTP контейнер.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/http/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `domain` | `string` | Optional, Computed | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `git_revision` | `string` | Optional | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | — |
| `ingress_annotations` | `string` | Optional, Computed | Формат ключ/значение<br/>Поддерживает пока простые аннотации |
| `json_env` | `string` | Optional, Computed | Массив переменных, будут переданы в env приложения |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `registry_path` | `string` | Optional, Computed | Путь до image, который будет установлен как deployment |
| `resource_cpu` | `number` | Optional, Computed | Квота ядра пода (Указывается в mili) |
| `resource_instances` | `number` | Optional, Computed | Количество узлов |
| `resource_memory` | `number` | Optional, Computed | Квота памяти пода (Указывается в Megabytes) |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Required | Платформа для развертывания Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
