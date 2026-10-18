# `nubes_redis` — Redis

Кластер Redis под кеши (k8s-operator, порт 6379). Отдаёт адрес подключения и
учётные данные администратора.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/redis.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_redis" "cache" {
  resource_name     = "my-redis"
  operation_timeout = "15m"

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
| `cpu` | integer | нет | `500` | Квота CPU пода (milicores) |
| `memory` | integer | нет | `512` | Квота памяти пода (MB) |
| `replicas` | integer | нет | `1` | Число реплик: `1`, `3`, `5`, `7` |
| `disk` | integer | нет | `10` | Квота диска (GB) |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `master_ip_space` | string | **да** | — | IP-space для мастера |
| `master_access_list` | string (json) | нет | `[]` | Разрешённые адреса мастера |
| `slave_ip_space` | string | **да** | — | IP-space для слейва |
| `slave_access_list` | string (json) | нет | `[]` | Разрешённые адреса слейва |

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
| `state_out` / `state_out_flat` | map | Ответ API (адрес подключения, `instanceUid`, `monitoring`) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

> **Точки подключения**
> - `state_out_flat["internalConnect.master"]` — хост Redis
> - `vault_secrets["adminPass"]` — пароль (строка, не JSON)
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

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание кластера |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `access_configuration` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/redis.md).
