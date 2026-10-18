# `nubes_postgres` — PostgreSQL

Объектно-реляционная СУБД PostgreSQL (Zalando Postgres Operator). Поддерживает
репликацию, пулеры соединений, TLS, бэкапы в S3 и автоскейлинг.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/postgres.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
> ℹ️ Для Nextcloud нужен PostgreSQL с `ssl_required = false`.

---

## Пример использования

```hcl
resource "nubes_postgres" "db" {
  resource_name     = "my-postgres"
  operation_timeout = "20m"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  cluster_configuration = {
    cpu      = 500
    memory   = 512
    replicas = 1 # 1 | 3 | 5 | 7
    disk     = 10
  }
  access_configuration = {
    master_ip_space    = "no-needed"
    master_access_list = jsonencode([])
    slave_ip_space     = "no-needed"
    slave_access_list  = jsonencode([])
  }
  postgres_configuration = {
    version       = "17" # 17 | 16
    pooler_master = false
    pooler_slave  = false
    ssl_required  = true
  }
  backup_configuration = {
    s3_uid   = nubes_s3.storage.id
    retain   = 7
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

> `access_configuration.master_ip_space` и `slave_ip_space` **обязательны**. Если
> внешний адрес не нужен, указывайте `no-needed`.

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
| `replicas` | integer | нет | `1` | Число реплик: `1`, `3`, `5`, `7` |
| `disk` | integer | нет | `10` | Квота диска (GB) |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `master_ip_space` | string | **да** | — | IP-space для мастера (`no-needed` — без внешнего адреса) |
| `master_access_list` | string (json) | нет | `[]` | Разрешённые адреса мастера |
| `slave_ip_space` | string | **да** | — | IP-space для слейва (`no-needed` — без внешнего адреса) |
| `slave_access_list` | string (json) | нет | `[]` | Разрешённые адреса слейва |

### `postgres_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `version` | string | нет | `17` | Версия PostgreSQL: `17`, `16` |
| `pooler_master` | bool | нет | `false` | Пулер соединений для мастера |
| `pooler_slave` | bool | нет | `false` | Пулер соединений для слейва |
| `ssl_required` | bool | нет | `true` | Требовать TLS-соединение |

### `postgres_conf` *(опционально)*

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `log_connections` | string | нет | `off` | `off`, `on` |
| `log_disconnections` | string | нет | `off` | `off`, `on` |
| `max_connections` | integer | нет | `100` | Максимум одновременных подключений |

```hcl
postgres_conf = {
  log_connections    = "off"
  log_disconnections = "off"
  max_connections    = 100
}
```

### `backup_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `s3_uid` | string | нет | — | UUID или имя экземпляра S3 (сервис 12); имя резолвится в UUID |
| `retain` | integer | нет | `7` | Срок хранения бэкапов |
| `schedule` | string | нет | `0 0 * * *` | Cron-расписание бэкапов |
| `increment_count` | integer | нет | `0` | Число инкрементальных бэкапов |

### `autoscale_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `enabled` | bool | нет | `false` | Включить автоскейлинг |
| `percent` | integer | нет | `10` | Порог: `10`, `15`, `20` |
| `quota` | integer | нет | `100` | Квота автоскейлинга |
| `schedule` | integer | нет | `0` | Час запуска: `0`, `5`, `12`, `23` |

### `mtls_configuration` *(опционально)*

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `type` | string | нет | `off` | `off`, `manual`, `auto` |
| `duration_server` | integer | нет | `87600` | Срок жизни серверного сертификата (часы) |
| `duration_ca` | integer | нет | `175200` | Срок жизни CA (часы) |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `20m`) |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить вместо удаления |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.internalMaster` | string | Внутренний адрес мастера (`postgresqlk8s-master.<uid>.svc.cluster.local`) |
| `state_out.internalSlave` | string | Внутренний адрес слейва |
| `state_out_flat` | map | `instanceUid`, `internalMaster`, `externalMaster`, `monitoring.*` |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

> **Точки подключения**
> - `state_out.internalMaster` — внутренний адрес мастера
> - `state_out.externalMaster` — внешний адрес (если задан `master_ip_space`)
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
- `reconcile` запускается ресурсом `nubes_service_operation`; `recovery`, `restart` **недоступны** (только через ЛК).

---

## Подресурсы

### `nubes_postgres_user`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `postgres_id` | string (uuid) | **да** | — | UUID кластера PostgreSQL | — |
| `username` | string | **да** | — | Имя пользователя | длина 2…62 |
| `role` | string | **да** | — | Роль пользователя | `app_user` — базовые права (SELECT/INSERT/UPDATE/DELETE); `ddl_user` — расширенные (createdb, createrole) |
| `mtls_access` | bool | нет | `false` | Доступ по mTLS | `false`, `true` |
| `password` | string (sensitive) | — | генерируется | Пароль пользователя | только чтение |


See the [Terraform Registry attribute reference](../reference/resources/postgres.md).
