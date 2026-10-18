---
page_title: "nubes_harbor Resource"
description: "Service: Container Registry."
subcategory: "Приложения"
---

# `nubes_harbor` Resource

Service: Container Registry.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/harbor/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/harbor/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `auth_type` | `string` | Optional, Computed | No | Определяет, какой тип авторизации будет использован после создания инстанса:<br/>- local - локальные пользователи<br/>- keycloak - авторизация через УЗ от личного кабинета Не изменяется после создания. |
| `domain` | `string` | Optional, Computed | No | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `pv_type` | `string` | Optional, Computed | No | Определяет тип хранилища, где будет храниться образ:<br/><br/>local - на локальных дисках. Необходимо также будет настроить размер хранилища<br/>s3 - будет создан бакет в s3. Необходимо также будет выбрать instanceUid s3 Не изменяется после создания. |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_cpu` | `number` | Optional, Computed | No | Квота ядра пода (Указывается в mili). 1000Mili = 1 Ядро |
| `resource_disk` | `number` | Optional, Computed | No | Необходимо настраивать, когда выбран тип хранилища - **local** |
| `resource_instances` | `number` | Optional, Computed | No | Количество узлов |
| `resource_memory` | `number` | Optional, Computed | No | Квота памяти пода (Указывается в Megabytes) |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Платформа для развертывания Не изменяется после создания. |
| `s3_uid` | `string` | Optional | No | UUID экземпляра услуги S3 Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
