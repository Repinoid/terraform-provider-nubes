# `nubes_mariadb` — Mariadb

Реляционная СУБД MariaDB (совместима с MySQL), разворачивается через `mariadb-operator`
в Kubernetes. Поддерживает репликацию, бэкапы в S3 и автоскейлинг.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/mariadb.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_mariadb" "db" {
  resource_name = "my-mariadb"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  cluster_configuration = {
    cpu      = 500
    memory   = 512
    replicas = 1 # только 1
    disk     = 10
  }
  access_configuration = {
    master_ip_space    = "no-needed"
    master_access_list = jsonencode([])
  }
  mariadb_configuration = {
    version = "9.4.0" # 8.4.0 | 9.4.0
  }
  backup_configuration = {
    s3_uid   = nubes_s3.storage.id # ⚠️ обязателен
    retain   = 720
    schedule = "0 0 * * *"
  }
  autoscale_configuration = {
    enabled  = false
    percent  = 10
    quota    = 100
    schedule = 0
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
| `resource_realm` | string | **да** | — | Платформа: `k8s-4-ext-nubes-ru` |

### `cluster_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `cpu` | integer | нет | `500` | Квота CPU пода (milicores) |
| `memory` | integer | нет | `512` | Квота памяти пода (MB) |
| `replicas` | integer | нет | `1` | Число реплик (только `1`) |
| `disk` | integer | нет | `10` | Квота диска (GB) |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `master_ip_space` | string | **да** | — | IP-space для доступа (`no-needed` — без внешнего адреса) |
| `master_access_list` | string (json) | нет | `[]` | Список разрешённых адресов |

### `mariadb_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `version` | string | нет | `9.4.0` | Версия: `8.4.0`, `9.4.0` |

### `backup_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `s3_uid` | string | **да** | — | UUID или имя экземпляра S3 (сервис 12); имя резолвится в UUID |
| `retain` | integer | нет | `720` | Срок хранения бэкапов |
| `schedule` | string | нет | `0 0 * * *` | Cron-расписание бэкапов |

### `autoscale_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `enabled` | bool | нет | `false` | Включить автоскейлинг |
| `percent` | integer | нет | `10` | Порог: `10`, `15`, `20` |
| `quota` | integer | нет | `100` | Квота автоскейлинга |
| `schedule` | integer | нет | `0` | Час запуска: `0`, `5`, `12`, `23` |

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
| `state_out` / `state_out_flat` | map | Ответ API (`instanceUid`, `internalConnect.master`, `monitoring`) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

> **Точки подключения**
> - `state_out_flat["internalConnect.master"]` = `mariak8s-mariadb-cluster.<uid>.svc.cluster.local`
> - `vault_secrets["adminUser"]`, `vault_secrets["adminPass"]` — креды администратора
>
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
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Подресурсы

### `nubes_mariadb_database`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `mariadb_id` | string (uuid) | **да** | — | UUID кластера MariaDB |
| `db_name` | string | нет | `db-<случайное>` | Имя базы данных |

### `nubes_mariadb_user`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `mariadb_id` | string (uuid) | **да** | — | UUID кластера MariaDB | — |
| `username` | string | нет | `user<случайное>` | Имя пользователя | длина 2…62; только латиница и цифры |
| `password` | string (sensitive) | нет | генерируется | Пароль пользователя | длина 8…62 |
| `db_name` | map (список) | нет | `[{}]` | Базы, к которым даётся полный доступ | `dbname` — **да**; база должна существовать |
| `access_hosts` | string | нет | `%` | MySQL host-паттерн, с которого разрешено подключение пользователю (напр. `%` — без ограничений, или конкретная подсеть). Определяет `user@host` в самой MariaDB | — |

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание кластера |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `access_configuration`, `mariadb_configuration`, `autoscale_configuration` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Остановка кластера (поды реально останавливаются) |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Запуск кластера |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `create_database` / `delete_database` | ✅ через `nubes_mariadb_database` | Управление БД |
| `create_user` / `delete_user` | ✅ через `nubes_mariadb_user` | Управление пользователями |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/mariadb.md).
