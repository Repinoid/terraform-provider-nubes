# `nubes_kafka` — ApacheKafka

Отказоустойчивый брокер сообщений Apache Kafka для потоковой обработки данных и обмена
событиями между сервисами. Сервис создаёт кластер и отдаёт адреса подключения
(bootstrap-серверы брокеров).

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/kafka.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_kafka" "cluster" {
  resource_name     = "my-kafka"
  operation_timeout = "20m"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  cluster_configuration = {
    cpu      = 2000
    memory   = 2048
    disk     = 1
    replicas = 1 # ⚠️ обязателен: 1 | 3 | 5 | 7
  }
  access_configuration = {
    allocation_external = false
    master_access_list  = jsonencode([])
  }
}
```

Типовая связка с топиком и пользователем:

```hcl
resource "nubes_kafka_topic" "events" {
  kafka_id   = nubes_kafka.cluster.id
  name_topic = "events"
  partitions = 3
  replicas   = 1
}

resource "nubes_kafka_user" "app" {
  kafka_id   = nubes_kafka.cluster.id
  name_topic = "events"
  username   = "app-user"
  operations = "Create,Describe,Read,Write" # без пробелов
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
| `cpu` | integer | нет | `2000` | Квота CPU пода (milicores) |
| `memory` | integer | нет | `2048` | Квота памяти пода (MB) |
| `disk` | integer | нет | `1` | Квота диска (GB) |
| `replicas` | integer | **да** | — | Число реплик: `1`, `3`, `5`, `7` |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `allocation_external` | bool | нет | `false` | Выделить внешний адрес |
| `master_access_list` | string (json) | нет | `[]` | Список разрешённых адресов (JSON-массив) |

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
| `state_out.internalConnect` | string (json) | Точки подключения: `bootstrap.fqdn`, `bootstrap.port` (=9093), `brokers.<name>.fqdn`, `brokers.<name>.port` |
| `state_out.externalConnect` | string | Внешний адрес (`""` без `allocation_external = true`) |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_out_flat["monitoring.logs"]` | string | Ссылка на логи (Grafana) |
| `state_out_flat["monitoring.resourceMetrics"]` | string | Ссылка на метрики (Grafana) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

```hcl
output "kafka_bootstrap" {
  value = jsondecode(nubes_kafka.cluster.state_out_flat["internalConnect"])
}
```

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

---

## Подресурсы

### `nubes_kafka_topic`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `kafka_id` | string (uuid) | **да** | — | UUID кластера Kafka | — |
| `name_topic` | string | **да** | — | Имя топика | — |
| `partitions` | integer | **да** | — | Число партиций | можно увеличивать, уменьшать нельзя |
| `replicas` | integer | **да** | — | Число реплик | не больше числа реплик брокера |

Идентификатор: `<kafka_uid>:topic:<hash>`.

### `nubes_kafka_user`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `kafka_id` | string (uuid) | **да** | — | UUID кластера Kafka | — |
| `username` | string | **да** | — | Имя пользователя | — |
| `name_topic` | string | **да** | — | Полное имя топика | маска запрещена |
| `operations` | list | **да** | — | Разрешённые операции | `Create`, `Describe`, `Read`, `Write` — **без пробелов** |
| `group` | string | нет | `*` | Consumer-группа | — |
| `access_hosts` | string | нет | `*` | Разрешённые хосты | — |
| `password` | string (sensitive) | — | из Vault | Пароль пользователя | только чтение |

Идентификатор: `<kafka_uid>:user:<hash>`.

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание кластера |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `access_configuration` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `create_topic` / `modify_topic` / `delete_topic` | ✅ через `nubes_kafka_topic` | Управление топиками |
| `create_user` / `modify_user` / `delete_user` | ✅ через `nubes_kafka_user` | Управление пользователями |


See the [Terraform Registry attribute reference](../reference/resources/kafka.md).
