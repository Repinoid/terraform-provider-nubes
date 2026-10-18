# `nubes_valo_tenant` — VALO Cloud

Логически изолированный тенант сервиса VALO Cloud: защищённое облачное хранилище,
синхронизация данных и доступ к документам с разных устройств. Создаётся в тестовом
режиме на 2 месяца; тенант отдаёт адрес подключения.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/valo_tenant.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_valo_tenant" "tenant" {
  resource_name = "my-valo"
  domain        = "myvalotenant" # одиночная метка, без точек
  # resource_realm = "valo.ngcloud.ru" # дефолт (фиксирован)

  s3_storage = 10 # 1…999 (GB)

  valo_options = {
    ff_logging                 = true
    ff_online_editor           = true
    ff_public_link             = true
    ff_themes                  = true
    ff_user_import             = true
    ff_watermark               = true
    request_approve_max_period = 2
    user_count                 = 10
  }
}
```

> **Обратите внимание:** `domain` проверяется паттерном `^[a-z0-9](?:[a-z0-9-]{3,61}[a-z0-9])?$` —
> **одиночная метка без точек** (FQDN не поддерживается).

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `domain` | string | **да** | — | Имя тенанта (добавляется в зону `valo.ngcloud.ru`) | ⛔ без точек, длина 3…61 |
| `resource_realm` | string | нет | `valo.ngcloud.ru` | «Платформа» — сервис VALO Cloud | `valo.ngcloud.ru` |
| `s3_storage` | integer | нет | `10` | Объём хранилища (GB) | 1…999 |

### `valo_options`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `ff_logging` | bool | нет | `true` | Журналирование |
| `ff_online_editor` | bool | нет | `true` | Онлайн-редактор документов |
| `ff_public_link` | bool | нет | `true` | Публичные ссылки |
| `ff_themes` | bool | нет | `true` | Темы оформления |
| `ff_user_import` | bool | нет | `true` | Импорт пользователей |
| `ff_watermark` | bool | нет | `true` | Водяные знаки |
| `request_approve_max_period` | integer | нет | `2` | Максимальный период согласования (дней) |
| `user_count` | integer | нет | `10` | Количество пользователей |

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
| `state_out.isDemo` | string | Признак демо-тенанта (`"true"`) |
| `state_out.urlConnect` | string | Адрес тенанта: `<domain>.valo.ngcloud.ru` |
| `state_params.domain` | string | Имя тенанта |
| `state_params.resourceRealm` | string | `valo.ngcloud.ru` |
| `state_params.s3Storage` | integer | Объём хранилища (GB) |
| `state_params.expireAt` | string | Дата истечения демо-периода (~2 месяца) |
| `state_params.valoOptions` | string (json) | Настройки VALO |
| `vault_secrets` | map (sensitive) | Секретов нет |
| `vault_fields` | list | `[]` |

```hcl
output "valo_url" { value = nubes_valo_tenant.tenant.state_out_flat["urlConnect"] }
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
> **Демо-период**
> Тенант создаётся в тестовом режиме (`isDemo = "true"`) сроком на 2 месяца
> (`expireAt`). Перевод в промышленную эксплуатацию — вне Terraform.
>
- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание тенанта |
| `modify` | ✅ `apply` | Меняются `valo_options`, `s3_storage` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |


See the [Terraform Registry attribute reference](../reference/resources/valo_tenant.md).
