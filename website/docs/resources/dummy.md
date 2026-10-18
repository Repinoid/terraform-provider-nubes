# `nubes_dummy` — Болванка

Тестовая «болванка» платформы: создаёт инстанс, который можно заставить ждать,
падать на старте или в процессе, и записывать/не записывать данные в CMDB и Vault.
Служит для проверки поведения провайдера на разных сценариях (успех, отказ,
откат), а не для рабочей нагрузки.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/dummy.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_dummy" "test" {
  resource_name  = "my-dummy"
  resource_realm = "dummy"

  duration_ms = 0     # задержка перед завершением
  fail_at_start  = false
  fail_in_progress = false

  bodymessage = "hello"
  json_example = jsonencode({ foo = "bar" })

  map_example = {
    subparam1 = "a"
    subparam2 = "b"
    secret    = "s3cr3t"
  }

  map_fixed = {
    param1   = 1
    param2   = "1"
    param3   = 1
    bol1     = false
    min_str1 = "abracadabra"
  }

  arr = ["test"]
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `duration_ms` | integer | нет | `0` | Пауза (sleep) перед завершением, мс | — |
| `fail_at_start` | bool | нет | `false` | Аварийно завершить джоб, ничего не записав в CMDB | `false`, `true` |
| `fail_in_progress` | bool | нет | `false` | Начать джоб и упасть в процессе, не записав в CMDB | `false`, `true` |
| `where_fail` | integer | нет | `1` | Где упасть (с `fail_in_progress`): 1 — prepare, 2 — блок Data fill, 3 — после записи в Vault | `1`, `2`, `3` |
| `resource_realm` | string | нет | `dummy` | «Платформа» развёртывания | `dummy` |
| `bodymessage` | string | нет | — | Дата, уходящая в Vault | — |
| `json_example` | string | нет | — | Для тестирования интерфейса | — |
| `nested_ref_example` | string | нет | — | — | — |
| `arr` | list | **да** | — | Тестовый список | — |

### `map_example`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `subparam1` | string | нет | — | Для администратора |
| `subparam2` | string | нет | — | — |
| `secret` | string | нет | — | Секрет |

### `map_fixed`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `param1` | integer | нет | `1` | Integer, секрет |
| `param2` | string | нет | `1` | `1`, `2`, `3` |
| `param3` | integer | нет | `1` | Целое положительное |
| `bol1` | bool | нет | `false` | Булевский |
| `min_str1` | string | нет | `abracadabra` | Не меньше 10 символов |

### `array_map_fixed_example`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `param1` | integer | **да** | — | — |
| `param2` | string | **да** | — | `1`, `2`, `3` |
| `param3` | integer | нет | `1` | — |
| `bol1` | bool | **да** | — | `false`, `true` |
| `min_str1` | string | **да** | — | Не меньше 10 символов |

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
| `state_out` / `state_out_flat` | map | Ответ API |
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
| `keep_on_destroy` | `false` | Если `true` — облако не трогается, ресурс только убирается из состояния (`state_only`); приоритетнее `suspend_on_destroy` |

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

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ✅ `apply` | Меняются `duration_ms`, `fail_at_start`, `fail_in_progress`, `where_fail`, `bodymessage`, `map_example`, `json_example` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |
| `redeploy` | ✅ `apply` (при смене `git_revision`) | Пересборка из git (поле-триггер) |


See the [Terraform Registry attribute reference](../reference/resources/dummy.md).
