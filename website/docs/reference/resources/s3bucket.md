# `nubes_s3bucket` Resource

Service: S3 бакет.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/s3bucket/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `bucket_name` | `string` | Required | No | Имя бакета s3 Не изменяется после создания. |
| `cors_all` | `bool` | Optional, Computed | No | Публичный доступ к чтению настроек позволяет любому пользователю интернета читать конфигурации CORS Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `list_all` | `bool` | Optional, Computed | No | Публичный доступ к листингу позволяет любому пользователю интернета получить список объектов в бакете. Не изменяется после создания. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `max_size` | `string` | Optional, Computed | No | Максимальный размер бакета s3. Если поставить -1, будет безлимит. Размер указывается в Bytes Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | No | — |
| `placement` | `string` | Optional, Computed | No | Класс хранилища, который используется по умолчанию при загрузке объектов. Чем холоднее, тем дешевле место, но дороже обращения Не изменяется после создания. |
| `read_all` | `bool` | Optional, Computed | No | Публичный доступ на чтение позволяет любому пользователю интернета скачивать объекты из бакета. Не изменяется после создания. |
| `resource_name` | `string` | Required | No | — |
| `s3_user_uid` | `string` | Optional | No | UID Корневой услуги S3. На основе него будет создан бакет. Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
