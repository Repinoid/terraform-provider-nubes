# `nubes_s3bucket` — S3 бакет

Бакет внутри существующего S3-хранилища (`nubes_s3`) с настройками публичного доступа
и класса хранения. Управляет только бакетом — данные в нём не создаются.

> **Специфика Nubes: `destroy` для этого ресурса **удаляет** бакет**
> `terraform destroy` удаляет бакет из облака и из Terraform-state. Усыпление
> (`suspend`) у этого сервиса не поддерживается. Подробности — в разделе
> [«Жизненный цикл»](#zhiznennyy-cikl) и в
> общем руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_s3" "storage" {
  resource_name = "my-s3"
  display_name  = "my_s3_storage"
}

resource "nubes_s3bucket" "backups" {
  resource_name = "my-backups-bucket"

  s3_user_uid = nubes_s3.storage.id # UUID корневого S3
  bucket_name = "backups"           # ^[a-z0-9-]+$

  max_size  = "-1"   # -1 = безлимит, иначе размер в байтах
  placement = "HOT"  # HOT | COLD

  read_all = false   # публичное чтение объектов
  list_all = false   # публичный листинг
  cors_all = false   # публичное чтение настроек CORS
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `bucket_name` | string | **да** | — | Имя бакета S3 | `^[a-z0-9-]+$` |
| `s3_user_uid` | string (uuid) | нет | — | UID корневой услуги S3 (сервис 12), на основе которого создаётся бакет | ссылка на `nubes_s3` |
| `max_size` | string | нет | `-1` | Максимальный размер бакета в байтах; `-1` — безлимит | — |
| `read_all` | bool | нет | `false` | Публичный доступ на чтение объектов | `false`, `true` |
| `list_all` | bool | нет | `false` | Публичный доступ к листингу объектов | `false`, `true` |
| `cors_all` | bool | нет | `false` | Публичный доступ к чтению настроек CORS | `false`, `true` |
| `placement` | string | нет | `HOT` | Класс хранения по умолчанию (холоднее — дешевле место, дороже обращения) | `HOT`, `COLD` |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса (используется как `s3_bucket_uid` в `nubes_k8s_velero`) |
| `state_out` / `state_out_flat` | map | Ответ API |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

`terraform destroy` удаляет бакет из облака и из Terraform-state. Усыпление (`suspend`)
у этого сервиса не поддерживается — атрибута `suspend_on_destroy` нет.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — бакет остаётся в облаке (`state_only`); удаление в облаке не выполняется |

- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание бакета |
| `delete` | ✅ `destroy` | Удаление бакета |

> ℹ️ Доступны только `create` и `delete` — изменение параметров бакета не поддерживается.


See the [Terraform Registry attribute reference](../reference/resources/s3bucket.md).
