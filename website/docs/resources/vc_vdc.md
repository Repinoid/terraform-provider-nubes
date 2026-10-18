# `nubes_vc_vdc` — Виртуальный датацентр (vDC)

Логически изолированная среда в VMware Cloud Director для размещения и управления
виртуальными машинами, сетевой инфраструктурой и ресурсами хранения.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/vc_vdc.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_vc_org" "org" {
  resource_name = "my-org"
}

resource "nubes_vc_vdc" "vdc" {
  resource_name = "my-vdc"

  organization_uid = nubes_vc_org.org.id
  provider_vdc     = "Intel Cascade Lake 3.1" # «Кластер ресурсов» из ЛК
  network_provider = "msk1"                   # «Зона доступности» из ЛК

  storage_config = [{
    name = "SSD SDS" # имя политики из ЛК
    size = 2000
  }]

  cpu_guaranteed = 0   # 0 | 50 | 80 (% резервирования)
  cpu_allocated  = 80  # 1…1000 (шт)
  mem_allocated  = 200 # 1…4096 (GB)
  suffix_name    = "vdc1"
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `provider_vdc` | string | **да** | — | Provider VDC из ЛК | — |
| `network_provider` | string | **да** | — | Сетевой провайдер из ЛК | — |
| `organization_uid` | string (uuid) | нет | — | UUID организации (сервис 19) | ссылка на `nubes_vc_org` |
| `cpu_guaranteed` | integer | нет | `0` | % резервирования CPU | `0`, `50`, `80` |
| `cpu_allocated` | integer | нет | `80` | Выделенные CPU (шт) | 1…1000 |
| `mem_allocated` | integer | нет | `200` | Выделенная память (GB) | 1…4096 |
| `suffix_name` | string | нет | `<случайное>` | Суффикс имени | длина 1…16 |

### `storage_config` *(список)*

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `name` | string | **да** | — | Имя хранилища (из конфигурации ресурсной платформы) |
| `size` | integer | нет | `2000` | Размер (GB) |

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
| `create` | ✅ `apply` | Создание vDC |
| `modify` | ✅ `apply` | Меняются `cpu_allocated`, `mem_allocated`, `storage_config` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/vc_vdc.md).
