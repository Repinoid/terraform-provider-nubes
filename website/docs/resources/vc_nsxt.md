# `nubes_vc_nsxt` — Сетевой шлюз периметра (Edge)

Создаёт периметровый сетевой шлюз (Edge) в рамках одного виртуального датацентра (vDC)
или группы vDC. При создании автоматически разворачивается routed-сеть с адресным
пространством `10.10.102.0/24`.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/vc_nsxt.md).
>
> **Специфика Nubes: `destroy` для этого ресурса **удаляет** шлюз**
> `terraform destroy` удаляет шлюз из облака и из Terraform-state. Усыпление (`suspend`)
> у этого сервиса не поддерживается. Подробности — в разделе
> [«Жизненный цикл»](#zhiznennyy-cikl) и в общем руководстве
> [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_vc_nsxt" "edge" {
  resource_name = "my-edge"

  vdc_uid  = nubes_vc_vdc.this.id # UUID виртуального датацентра (сервис 21)
  vdc_type = "vdc"                # vdc | vdcGroup

  need_enable_avi         = false
  virtual_services_count  = 1 # 1…4

  routed_net_configuration = {
    ip_addr_pool = "10.10.102.0/24"
    main_dns     = "81.22.46.22"
    second_dns   = "185.247.187.77"
  }
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `vdc_uid` | string | нет | — | UUID услуги «Виртуальный датацентр (vDC)» (сервис 21) | ссылка на `nubes_vc_vdc` |
| `vdc_group_uid` | string | нет | — | UUID услуги «Группа виртуальных датацентров» (сервис 29) | ссылка на `nubes_vc_vdc_group` |
| `vdc_type` | string | нет | `vdc` | Тип родительской услуги | `vdc`, `vdcGroup` |
| `need_enable_avi` | bool | нет | `false` | Включить Load Balancer (AVI) | `false`, `true` |
| `virtual_services_count` | integer | нет | `1` | Количество виртуальных сервисов, резервируемых на AVI | 1…4 |
| `qos_profile` | string | нет | — | QoS-профиль (⚠️ **пока не работает** — параметр не возвращается платформой) | — |

### `routed_net_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `ip_addr_pool` | string | нет | `10.10.102.0/24` | Адресное пространство routed-сети (CIDR) |
| `main_dns` | string | нет | `81.22.46.22` | Основной DNS |
| `second_dns` | string | нет | `185.247.187.77` | Дополнительный DNS |

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

`terraform destroy` удаляет шлюз из облака и из Terraform-state. Усыпление (`suspend`)
у этого сервиса не поддерживается — атрибута `suspend_on_destroy` нет.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — ресурс остаётся в облаке (`state_only`); удаление в облаке не выполняется |

- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание шлюза |
| `modify` | ✅ `apply` | Меняются `need_enable_avi`, `virtual_services_count`, `ip_space_name`, `qos_profile`, `routed_net_configuration` |
| `delete` | ✅ `destroy` | Удаление шлюза |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/vc_nsxt.md).
