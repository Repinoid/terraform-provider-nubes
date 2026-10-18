---
page_title: "nubes_pgadmin"
description: "nubes_pgadmin TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_pgadmin` Resource

Service: pgAdmin.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/pgadmin/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `domain` | `string` | Optional, Computed | No | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `git_revision` | `string` | Optional | No | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `login` | `string` | Required | No | Логин для подключения к pgAdmin, в формате email (login@email.net) Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | No | — |
| `password` | `string` | Required | Yes | Пароль для подключения к pgAdmin Не изменяется после создания. |
| `resource_cpu` | `number` | Optional, Computed | No | Квота ядра пода (milicores). 1000Mili = 1 Ядро |
| `resource_disk` | `number` | Optional, Computed | No | Квота диска (Указывается в GIGAbytes) |
| `resource_memory` | `number` | Optional, Computed | No | Квота памяти пода (MBytes) |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Required | No | Платформа для развертывания Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
