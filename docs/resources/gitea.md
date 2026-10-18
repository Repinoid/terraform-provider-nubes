---
page_title: "nubes_gitea Resource"
description: "Service: Gitea."
subcategory: "Приложения"
---

# `nubes_gitea` Resource

Service: Gitea.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/gitea/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/gitea/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `domain` | `string` | Optional, Computed | No | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `ingress_extra_annotations` | `string` | Optional, Computed | No | — |
| `ini_server_params` | `string` | Optional, Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `psql_uid` | `string` | Optional | No | UUID Услуги k8sPgsql. Пользователя techGitea, базу данных gitea создаёт автоматически Не изменяется после создания. |
| `resource_cpu` | `number` | Optional, Computed | No | Квота ядра пода (Указывается в mili). 1000Mili = 1 Ядро |
| `resource_disk` | `number` | Optional, Computed | No | Квота диска (Указывается в GIGAbytes) |
| `resource_instances` | `number` | Optional, Computed | No | Количество узлов |
| `resource_memory` | `number` | Optional, Computed | No | Квота памяти пода (Указывается в Megabytes) |
| `resource_name` | `string` | Required | No | — |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
