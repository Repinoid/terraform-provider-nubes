---
page_title: "nubes_s3bucket"
description: "nubes_s3bucket TEST resource schema for Nubes Cloud."
subcategory: "Хранилище"
---

# `nubes_s3bucket` Resource

Service: S3 бакет.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/s3bucket/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `bucket_name` | `string` | Required | Имя бакета s3 Не изменяется после создания. |
| `cors_all` | `bool` | Optional, Computed | Публичный доступ к чтению настроек позволяет любому пользователю интернета читать конфигурации CORS Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `list_all` | `bool` | Optional, Computed | Публичный доступ к листингу позволяет любому пользователю интернета получить список объектов в бакете. Не изменяется после создания. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `max_size` | `string` | Optional, Computed | Максимальный размер бакета s3. Если поставить -1, будет безлимит. Размер указывается в Bytes Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | — |
| `placement` | `string` | Optional, Computed | Класс хранилища, который используется по умолчанию при загрузке объектов. Чем холоднее, тем дешевле место, но дороже обращения Не изменяется после создания. |
| `read_all` | `bool` | Optional, Computed | Публичный доступ на чтение позволяет любому пользователю интернета скачивать объекты из бакета. Не изменяется после создания. |
| `resource_name` | `string` | Required | — |
| `s3_user_uid` | `string` | Optional | UID Корневой услуги S3. На основе него будет создан бакет. Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
