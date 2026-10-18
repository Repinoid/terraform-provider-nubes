# `nubes_vc_vdc_group` — Группа датацентров

Логически объединяет несколько виртуальных датацентров (vDC) в единую группу, позволяя
растягивать SDN между разными vDC и централизованно управлять сетевой связностью.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/vc_vdc_group.md).
>
> **Специфика Nubes: `destroy` для этого ресурса **удаляет** группу**
> `terraform destroy` удаляет группу из облака и из Terraform-state. Усыпление (`suspend`)
> у этого сервиса не поддерживается. Подробности — в разделе
> [«Жизненный цикл»](#zhiznennyy-cikl) и в общем руководстве
> [«Как работает провайдер»](../guides/provider-behavior.md).
>
> **Перед созданием группы**
> У входящего vDC должна быть выполнена операция `reconcile` (её можно запустить из Terraform ресурсом `nubes_service_operation` у ресурса vDC).
> В противном случае создание отклоняется на валидации.
>
---

## Пример использования

```hcl
resource "nubes_vc_vdc_group" "group" {
  resource_name = "my-vdc-group"

  vdc_uid = nubes_vc_vdc.this.id # UUID существующего vDC (сервис 21)
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса |
| `vdc_uid` | string (uuid) | нет | — | UUID существующего vDC (сервис 21) |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
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

`terraform destroy` удаляет группу из облака и из Terraform-state. Усыпление (`suspend`)
у этого сервиса не поддерживается — атрибута `suspend_on_destroy` нет.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — ресурс остаётся в облаке (`state_only`); удаление в облаке не выполняется |

> **Требуется `reconcile` у vDC**
> Платформа отказывает в создании группы, пока у существующего vDC не выполнена
> операция `reconcile`. Из Terraform её запускает ресурс `nubes_service_operation`
> у ресурса vDC (`nubes_vc_vdc`).
>
- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.

---

## Подресурсы

### `nubes_vc_vdc_group_vdc`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `vc_vdc_group_id` | string (uuid) | **да** | — | UUID группы датацентров |

Добавляет vDC в группу; удаление подресурса убирает vDC из группы.

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание группы |
| `delete` | ✅ `destroy` | Удаление группы |
| `add_vdc` / `remove_vdc` | ✅ через `nubes_vc_vdc_group_vdc` | Управление составом группы |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/vc_vdc_group.md).
