# `nubes_mongodb` — MongoDB

Документная СУБД MongoDB (PSMDB) с репликацией и резервным копированием в S3. Сервис
создаёт кластер и отдаёт внутренний адрес подключения.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/mongodb.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_mongodb" "db" {
  resource_name     = "my-mongodb"
  operation_timeout = "15m"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  cluster_configuration = {
    cpu      = 1000
    memory   = 1024
    disk     = 5
    replicas = 1
  }
  access_configuration = {
    master_ip_space    = "no-needed"
    master_access_list = jsonencode([])
  }
  backup_configuration = {
    s3_uid   = nubes_s3.storage.id # ⚠️ обязателен
    retain   = 7
    schedule = "0 0 * * *"
  }
  mongo_configuration = {
    tls = false # ⚠️ обязателен
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

### `cluster_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `cpu` | integer | нет | `1000` | Квота CPU пода (milicores) |
| `memory` | integer | нет | `1024` | Квота памяти пода (MB) |
| `disk` | integer | нет | `5` | Квота диска (GB) |
| `replicas` | integer | нет | `1` | Число реплик |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `master_ip_space` | string | **да** | — | IP-space для доступа |
| `master_access_list` | string (json) | нет | `[]` | Список разрешённых адресов (JSON-массив) |

### `backup_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `s3_uid` | string (uuid) | **да** | — | UUID экземпляра S3 (сервис 12) |
| `retain` | integer | нет | `7` | Срок хранения бэкапов |
| `schedule` | string | нет | `0 0 * * *` | Cron-расписание бэкапов |

### `mongo_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `tls` | bool | **да** | — | Включить TLS |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `15m`) |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить вместо удаления |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.internalConnect` | string | Внутренний адрес: `mongok8s-psmdb-db-master.<uid>.svc.cluster.local` |
| `state_out.backups` | string (json) | Список бэкапов (`{}`, пока нет бэкапов) |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_out_flat["monitoring.logs"]` | string | Ссылка на логи (Grafana) |
| `state_out_flat["monitoring.resourceMetrics"]` | string | Ссылка на метрики (Grafana) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | **Пусто** — админ-пользователей по умолчанию нет |
| `vault_fields` | list | `[]` |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |

> ℹ️ `internalConnect` — плоская строка, а не объект.
> ℹ️ `vault_secrets` пуст: пользователей создавайте через `nubes_mongodb_user`.

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy` | `true` | Инстанс переводится в `suspend` (не удаляется) |
| `suspend_on_destroy` | `false` | Выполняется операция `delete` — инстанс удаляется в облаке |
| `keep_on_destroy` | `false` | Если `true` — ресурс вообще не трогается в облаке (state_only) |

Возврат из `suspend`: повторный `terraform apply` с `adopt_existing_on_create = true`
— провайдер найдёт инстанс по имени и выполнит `resume` (тот же `id`).

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
- `reconcile` запускается ресурсом `nubes_service_operation`; `rollback` **недоступен** (только через ЛК).

---

## Подресурсы

### `nubes_mongodb_user`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `mongodb_id` | string (uuid) | **да** | — | UUID кластера MongoDB | — |
| `username` | string | **да** | — | Имя пользователя | — |
| `db_name` | string | **да** | — | База данных | — |
| `role` | string | **да** | — | Роль пользователя | `read`, `readWrite`, `dbAdmin`, `dbOwner` |
| `password` | string (sensitive) | — | генерируется | Пароль пользователя | только чтение |

### `nubes_mongodb_backup`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `mongodb_id` | string (uuid) | **да** | — | UUID кластера MongoDB | — |
| `backup_name` | string | **да** | — | Имя резервной копии | строчная латиница, цифры, дефис |

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание кластера |
| `modify` | ✅ `apply` | Изменение конфигурации кластера |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |
| `rollback` | ⛔ нет | Только вручную в ЛК — откат кластера к резервной копии |
| управление пользователями/бэкапами | ✅ через `nubes_mongodb_user` / `nubes_mongodb_backup` | |


See the [Terraform Registry attribute reference](../reference/resources/mongodb.md).
