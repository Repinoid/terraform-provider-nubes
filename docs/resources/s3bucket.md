---
page_title: "nubes_s3bucket"
description: "S3 bucket resource for the Nubes TEST stand."
subcategory: "Хранилище"
---

# `nubes_s3bucket` Resource

Creates a bucket within an existing S3 storage service. This resource manages
the bucket itself; it does not create objects in the bucket.

See the [TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/s3bucket/)
for service details and usage examples.

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | Adopt an existing instance by name during create. |
| `bucket_name` | `string` | Required | No | Имя бакета s3. |
| `cors_all` | `bool` | Optional, Computed | No | Публичный доступ к настройкам CORS. |
| `id` | `string` | Computed | No | UUID инстанса. |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Не удалять бакет из облака при `destroy`. |
| `list_all` | `bool` | Optional, Computed | No | Публичный доступ к списку объектов. |
| `log_level` | `string` | Optional | No | Уровень журналирования ресурса. |
| `max_size` | `string` | Optional, Computed | No | Максимальный размер бакета в байтах; `-1` — безлимит. |
| `operation_timeout` | `string` | Optional | No | Максимальное время ожидания операции. |
| `placement` | `string` | Optional, Computed | No | Класс хранения: `HOT` или `COLD`. |
| `read_all` | `bool` | Optional, Computed | No | Публичное чтение объектов. |
| `resource_name` | `string` | Required | No | Имя инстанса. |
| `s3_user_uid` | `string` | Optional | No | UID корневой услуги S3. |
| `state_out` | `["map","string"]` | Computed | No | API output. |
| `state_out_flat` | `["map","string"]` | Computed | No | Flattened API output. |
| `state_params` | `["map","string"]` | Computed | No | API parameters. |
| `state_params_flat` | `["map","string"]` | Computed | No | Flattened API parameters. |
| `vault_fields` | `["list","string"]` | Computed | No | Доступные поля Vault. |
| `vault_secrets` | `["map","string"]` | Computed | Yes | Секреты Vault. |
| `vault_url` | `string` | Computed | No | URL Vault. |
| `vault_user_path` | `string` | Computed | No | Путь пользователя Vault. |

## Lifecycle

This resource supports `create` and `delete`; parameter modification and
suspension are not supported. `keep_on_destroy = true` removes the object from
Terraform state without deleting the bucket from the cloud.
