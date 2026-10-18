# `nubes_rabbitmq` — RabbitMQ

Брокер сообщений с открытым исходным кодом для надёжной доставки сообщений между
приложениями. Сервис создаёт кластер и отдаёт внутренние и внешние точки подключения
плюс креды администратора.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/rabbitmq.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_rabbitmq" "mq" {
  resource_name     = "my-rabbitmq"
  operation_timeout = "20m"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  cluster_configuration = {
    cpu      = 1000
    memory   = 1024
    disk     = 5
    replicas = 1 # 1 | 3 | 5 | 7
  }
  access_configuration = {
    master_ip_space    = "no-needed"
    master_access_list = jsonencode([])
  }
  mq_configuration = {
    max_message_size = 4
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
| `cpu` | integer | нет | `1000` | Квота CPU пода (milicores) |
| `memory` | integer | нет | `1024` | Квота памяти пода (MB) |
| `disk` | integer | нет | `5` | Квота диска (GB) |
| `replicas` | integer | нет | `1` | Число реплик: `1`, `3`, `5`, `7` |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `master_ip_space` | string | **да** | — | Из какого IP-space резервировать внешний IP |
| `master_access_list` | string (json) | нет | `[]` | Список разрешённых адресов |

### `mq_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `max_message_size` | integer | нет | `4` | Максимальный размер сообщения (MB) |

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
| `state_out.internalConnect` | string | Внутренний адрес: `rabbitmqk8s.<uid>.svc.cluster.local` (плоская строка) |
| `state_out.externalConnect` | string (json) | Внешний адрес: `{ admin = "admin.<uid>.dev.nubes.ru" }` — присутствует даже без внешнего IP |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_out_flat["monitoring.logs"]` | string | Ссылка на логи (Grafana) |
| `state_out_flat["monitoring.resourceMetrics"]` | string | Ссылка на метрики (Grafana) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets.adminUser` | string (sensitive) | Логин администратора |
| `vault_secrets.adminPass` | string (sensitive) | Пароль администратора |
| `vault_fields` | list | `["adminUser", "adminPass"]` |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |

```hcl
output "rabbit_admin" { value = nubes_rabbitmq.mq.vault_secrets["adminUser"] }
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
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание кластера |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `access_configuration`, `mq_configuration` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/rabbitmq.md).
