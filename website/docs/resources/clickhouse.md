# `nubes_clickhouse` — ClickHouse

Колоночная СУБД для аналитики в реальном времени. Сервис создаёт кластер ClickHouse
(с опциональным ClickHouse Keeper, автодоскейлом и бэкапами в S3) и отдаёт адреса
подключения по HTTP и TCP.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/clickhouse.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> У этого ресурса `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`),
> а не удаляет. Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

Минимальный (обязательные параметры):

```hcl
resource "nubes_clickhouse" "db" {
  resource_name = "my-clickhouse"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  access_configuration = {
    master_ip_space = "no-needed"
  }
  backup_configuration = {
    s3_uid = nubes_s3.storage.id # ⚠️ обязателен
  }
}
```

Типовой (с явными ресурсами и бэкапами):

```hcl
resource "nubes_clickhouse" "db" {
  resource_name     = "my-clickhouse"
  operation_timeout = "20m"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  cluster_configuration = {
    cpu      = 1000 # milicores
    memory   = 2048 # MB
    disk     = 10   # GB
    replicas = 1    # 1 | 3 | 5 | 7
    shards   = 1    # 1 | 2 | 3
  }
  access_configuration = {
    master_ip_space    = "no-needed"
    master_access_list = jsonencode([])
  }
  clickhouse_keeper_configuration = {
    cpu      = 500
    memory   = 512
    replicas = 0 # 0 | 1 | 3
  }
  clickhouse_configuration = {
    version = "25.10"
  }
  backup_configuration = {
    s3_uid   = nubes_s3.storage.id
    retain   = 7
    schedule = "0 0 * * *"
  }
  autoscale_configuration = {
    enabled  = false
    percent  = 10  # 10 | 15 | 20
    quota    = 100
    schedule = 0   # 0 | 5 | 12 | 23
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
| `memory` | integer | нет | `2048` | Квота памяти пода (MB) |
| `disk` | integer | нет | `10` | Квота диска (GB) |
| `replicas` | integer | нет | `1` | Число реплик: `1`, `3`, `5`, `7` |
| `shards` | integer | нет | `1` | Число шардов: `1`, `2`, `3` |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `master_ip_space` | string | **да** | — | IP-space для доступа |
| `master_access_list` | string (json) | нет | `[]` | Список разрешённых адресов (JSON-массив) |

### `clickhouse_keeper_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `cpu` | integer | нет | `500` | Квота CPU Keeper (milicores) |
| `memory` | integer | нет | `512` | Квота памяти Keeper (MB) |
| `replicas` | integer | нет | `0` | Число реплик Keeper: `0`, `1`, `3` |

### `clickhouse_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `version` | string | нет | `25.10` | Версия ClickHouse |

### `backup_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `s3_uid` | string (uuid) | **да** | — | UUID экземпляра S3 (сервис 12) |
| `retain` | integer | нет | `7` | Срок хранения бэкапов |
| `schedule` | string | нет | `0 0 * * *` | Cron-расписание бэкапов |

### `autoscale_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `enabled` | bool | нет | `false` | Включить автоскейлинг |
| `percent` | integer | нет | `10` | Порог: `10`, `15`, `20` |
| `quota` | integer | нет | `100` | Квота автоскейлинга |
| `schedule` | integer | нет | `0` | Час запуска: `0`, `5`, `12`, `23` |

### `clickhouse_files_configuration` *(опционально)*

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `path` | string | нет | `something.xml` | Путь к XML-конфигу |
| `xml_config` | string | нет | `<clickhouse></clickhouse>` | Содержимое XML-конфига |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `20m`) |
| `log_level` | string | нет | `none` | Логирование стадий операции: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить (`true`) вместо удаления (`false`) |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply (=import) |

---

## Attribute Reference

Вычисляемые атрибуты (только чтение):

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.localConnectionUrlHttp` | string | HTTP-endpoint внутри кластера: `http://clickhouse-clickhousek8s.<uid>.svc.cluster.local:8123` |
| `state_out.localConnectionUrlTcp` | string | TCP-endpoint внутри кластера: `...:9000` |
| `state_out.externalConnectionUrlHttp` | string | Внешний HTTP-endpoint (`""` без внешнего доступа) |
| `state_out.externalConnectionUrlTcp` | string | Внешний TCP-endpoint (`""` без внешнего доступа) |
| `state_out.databases` | string (json) | Базы данных (`{}` пока нет БД) |
| `state_out.users` | string (json) | Пользователи (`{}` пока нет юзеров) |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_out_flat["monitoring.logs"]` | string | Ссылка на логи (Grafana) |
| `state_out_flat["monitoring.resourceMetrics"]` | string | Ссылка на метрики (Grafana) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | `["users"]` |

```hcl
output "ch_http" { value = nubes_clickhouse.db.state_out_flat["localConnectionUrlHttp"] }
```

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

⚠️ `terraform destroy` **не удаляет** инстанс, а **усыпляет** его
(`suspend_on_destroy = true` по умолчанию). Данные сохраняются.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy` | `true` | Инстанс переводится в `suspend` (не удаляется) |
| `suspend_on_destroy` | `false` | Выполняется операция `delete` — инстанс удаляется в облаке |
| `keep_on_destroy` | `false` | Если `true` — облако не трогается, ресурс только убирается из состояния (`state_only`); приоритетнее `suspend_on_destroy` |

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
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Подресурсы

### `nubes_clickhouse_database`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `clickhouse_id` | string (uuid) | **да** | — | UUID кластера ClickHouse |
| `db_name` | string | нет | `db01` | Имя базы данных |

### `nubes_clickhouse_user`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `clickhouse_id` | string (uuid) | **да** | — | UUID кластера ClickHouse | — |
| `username` | string | нет | `myusername` | Имя пользователя | длина 2…62 |
| `db_name` | map (список) | нет | `[{}]` | Базы, к которым даётся доступ | `db_name` — **да** |
| `access_hosts` | string | нет | `::/0` | Разрешённые адреса подключения | — |

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`. Идентификатор: `<clickhouse_uid>:<kind>:<hash>`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание кластера |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `clickhouse_keeper_configuration`, `access_configuration`, `clickhouse_configuration`, `backup_configuration`, `autoscale_configuration`, `clickhouse_files_configuration` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `create_database` / `delete_database` | ✅ через `nubes_clickhouse_database` | Управление БД |
| `create_user` / `delete_user` | ✅ через `nubes_clickhouse_user` | Управление пользователями |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/clickhouse.md).
