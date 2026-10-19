---
page_title: "nubes_superset"
description: "nubes_superset TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_superset` Resource

Service: Apache Superset.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/superset/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `domain` | `string` | Optional, Computed | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `emails` | `string` | Optional, Computed | Учётная запись, которой будут выданы админские права. В формате email [ну не совсем. Для доменных учетных записей - без доменного суффикса]. Обязана находиться в Keycloak. Можно передавать несколько через пробел Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `resource_cpu` | `number` | Optional, Computed | Квота ядра пода (Указывается в mili). 1000 Mili = 1 ядро Не изменяется после создания. |
| `resource_disk` | `number` | Optional, Computed | Квота диска (Указывается в GIGAbytes) Не изменяется после создания. |
| `resource_instances` | `number` | Optional, Computed | Количество узлов Не изменяется после создания. |
| `resource_memory` | `number` | Optional, Computed | Квота памяти пода (Указывается в Megabytes) Не изменяется после создания. |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Required | Платформа для развертывания Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
