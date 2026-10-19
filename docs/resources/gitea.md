---
page_title: "nubes_gitea"
description: "nubes_gitea TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_gitea` Resource

Service: Gitea.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/gitea/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `domain` | `string` | Optional, Computed | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `ingress_extra_annotations` | `string` | Optional, Computed | — |
| `ini_server_params` | `string` | Optional, Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `psql_uid` | `string` | Optional | UUID Услуги k8sPgsql. Пользователя techGitea, базу данных gitea создаёт автоматически Не изменяется после создания. |
| `resource_cpu` | `number` | Optional, Computed | Квота ядра пода (Указывается в mili). 1000Mili = 1 Ядро |
| `resource_disk` | `number` | Optional, Computed | Квота диска (Указывается в GIGAbytes) |
| `resource_instances` | `number` | Optional, Computed | Количество узлов |
| `resource_memory` | `number` | Optional, Computed | Квота памяти пода (Указывается в Megabytes) |
| `resource_name` | `string` | Required | — |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
