# `nubes_nextcloud` — Nextcloud

Платформа для хранения и совместной работы с файлами (файлы, календари, документы
Collabora). Разворачивается в Kubernetes через Helm-чарт `nextcloud/nextcloud`.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/nextcloud.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
> **Nextcloud не самодостаточен**
> Для создания нужны заранее развёрнутые в той же ресурсной платформе **PostgreSQL**
> (с параметром `ssl_required = false`) и **Redis**.
>
---

## Пример использования

```hcl
resource "nubes_postgres" "db" {
  resource_name = "my-nextcloud-pg"

  startup_configuration = { resource_realm = "k8s-4-ext-nubes-ru" }
  access_configuration  = {}
  postgres_configuration = {
    version      = "17"
    ssl_required = false # ⚠️ обязательно для Nextcloud
  }
  backup_configuration = {}
  autoscale_configuration = {}
}

resource "nubes_redis" "cache" {
  resource_name = "my-nextcloud-redis"

  startup_configuration = { resource_realm = "k8s-4-ext-nubes-ru" }
}

resource "nubes_nextcloud" "cloud" {
  resource_name = "my-nextcloud"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
    psql_uid       = nubes_postgres.db.id  # ⚠️ обязателен
    redis_uid      = nubes_redis.cache.id  # ⚠️ обязателен
    psql_enabled   = true
    redis_enabled  = true
  }
  cluster_configuration = {
    cpu      = 1000 # 1000 | 2000 | 3000 | 4000
    memory   = 1024 # 1024…8192 с шагом 1024
    replicas = 1    # 1 | 2
    disk     = 10
  }
  access_configuration = {
    domain = "my-nextcloud"
  }
  app_configuration = {
    version         = "33.0.4"
    document_server = "collabora"
    log_level       = "INFO"
  }
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса |

### `startup_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_realm` | string | **да** | — | Платформа: `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` |
| `psql_uid` | string (uuid) | **да** | — | UUID экземпляра PostgreSQL (сервис 90) с `sslRequired=false` |
| `redis_uid` | string (uuid) | **да** | — | UUID экземпляра Redis (сервис 91) |
| `psql_enabled` | bool | нет | `true` | Использовать PostgreSQL |
| `redis_enabled` | bool | нет | `true` | Использовать Redis |

### `cluster_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `cpu` | integer | нет | `1000` | Квота CPU пода (milicores): `1000`, `2000`, `3000`, `4000` |
| `memory` | integer | нет | `1024` | Квота памяти пода (MB): `1024`…`8192` |
| `replicas` | integer | нет | `1` | Число реплик: `1`, `2` |
| `disk` | integer | нет | `10` | Квота диска (GB) |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `domain` | string | нет | `nextcloud-<случайное>` | FQDN или имя |

### `app_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `version` | string | нет | `33.0.4` | Версия Nextcloud |
| `document_server` | string | нет | `collabora` | Сервер документов |
| `log_level` | string | нет | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`, `NONE` |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить вместо удаления |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out` / `state_out_flat` | map | Ответ API (адрес, `instanceUid`, `monitoring`) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy` | `true` | Инстанс переводится в `suspend` (не удаляется) |
| `suspend_on_destroy` | `false` | Выполняется операция `delete` — инстанс удаляется в облаке |
| `keep_on_destroy` | `false` | Если `true` — ресурс вообще не трогается в облаке (state_only) |

Возврат из `suspend`: повторный `terraform apply` с `adopt_existing_on_create = true`
— провайдер найдёт инстанс по имени и выполнит `resume`.

> **Зачем нужно усыпление**
> `suspend` — временная остановка: ресурсы высвобождаются, а инстанс можно
> вернуть в работу в любой момент повторным `apply`.
>
> **Удаление усыплённого инстанса**
> Удалить усыплённый инстанс можно, когда он пробыл в `suspend` **14 дней**.
>
> ⚠️ Из Terraform это невозможно — операция `delete` для усыплённого инстанса
> недоступна. Удаляйте его вручную в Личном кабинете (ЛК).
>
- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.
- `reconcile` запускается ресурсом `nubes_service_operation`; `recovery`, `restart` **недоступны** (только через ЛК).

---

## Подресурсы

### `nubes_nextcloud_backup`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `nextcloud_id` | string (uuid) | **да** | — | UUID инстанса Nextcloud |
| `startup_configuration` | map (список) | **да** | — | Параметры бэкапа |

`startup_configuration` (элемент списка):

| Имя | Тип | Обязательный | Описание |
|---|---|---|---|
| `velero_uid` | string | **да** | UUID сервиса Velero, выполняющего бэкап |

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `app_configuration` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `create_backup` | ✅ через `nubes_nextcloud_backup` | Резервное копирование |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |
| `recovery` | ⛔ нет | Только вручную в ЛК |
| `restart` | ⛔ нет | Только вручную в ЛК |


See the [Terraform Registry attribute reference](../reference/resources/nextcloud.md).
