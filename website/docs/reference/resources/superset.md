# `nubes_superset` Resource

Service: Apache Superset.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/superset/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/superset/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `domain` | `string` | Optional, Computed | No | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `emails` | `string` | Optional, Computed | No | Учётная запись, которой будут выданы админские права. В формате email [ну не совсем. Для доменных учетных записей - без доменного суффикса]. Обязана находиться в Keycloak. Можно передавать несколько через пробел Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `resource_cpu` | `number` | Optional, Computed | No | Квота ядра пода (Указывается в mili). 1000 Mili = 1 ядро Не изменяется после создания. |
| `resource_disk` | `number` | Optional, Computed | No | Квота диска (Указывается в GIGAbytes) Не изменяется после создания. |
| `resource_instances` | `number` | Optional, Computed | No | Количество узлов Не изменяется после создания. |
| `resource_memory` | `number` | Optional, Computed | No | Квота памяти пода (Указывается в Megabytes) Не изменяется после создания. |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Платформа для развертывания Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
