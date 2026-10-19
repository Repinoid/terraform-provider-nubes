---
page_title: "nubes_harbor"
description: "nubes_harbor TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_harbor` Resource

Service: Container Registry.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/harbor/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `auth_type` | `string` | Optional, Computed | Определяет, какой тип авторизации будет использован после создания инстанса:<br/>- local - локальные пользователи<br/>- keycloak - авторизация через УЗ от личного кабинета Не изменяется после создания. |
| `domain` | `string` | Optional, Computed | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `pv_type` | `string` | Optional, Computed | Определяет тип хранилища, где будет храниться образ:<br/><br/>local - на локальных дисках. Необходимо также будет настроить размер хранилища<br/>s3 - будет создан бакет в s3. Необходимо также будет выбрать instanceUid s3 Не изменяется после создания. |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_cpu` | `number` | Optional, Computed | Квота ядра пода (Указывается в mili). 1000Mili = 1 Ядро |
| `resource_disk` | `number` | Optional, Computed | Необходимо настраивать, когда выбран тип хранилища - **local** |
| `resource_instances` | `number` | Optional, Computed | Количество узлов |
| `resource_memory` | `number` | Optional, Computed | Квота памяти пода (Указывается в Megabytes) |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Required | Платформа для развертывания Не изменяется после создания. |
| `s3_uid` | `string` | Optional | UUID экземпляра услуги S3 Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
