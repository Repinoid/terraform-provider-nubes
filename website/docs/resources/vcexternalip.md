# `nubes_vcexternalip` — Публичные IP адреса

Резервирует публичный IP-адрес и привязывает его к существующему сервису правилами
DNAT/SNAT.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/vcexternalip.md).
>
> **Специфика Nubes: `destroy` для этого ресурса **освобождает** IP**
> `terraform destroy` освобождает публичный IP (удаляет его в облаке и из
> Terraform-state). Усыпление (`suspend`) у этого сервиса не поддерживается.
> Подробности — в разделе
> [«Жизненный цикл»](#zhiznennyy-cikl) и в общем руководстве
> [«Как работает провайдер»](../guides/provider-behavior.md).
>
> **Обычно эта услуга создаётся автоматически**
> Как правило, внешние IP резервируются автоматически при развёртывании более сложных
> услуг. Резервировать IP вручную можно только после создания Organization, vDC, vApp
> и Edge.
>
---

## Пример использования

```hcl
resource "nubes_vcexternalip" "ip" {
  resource_name = "my-external-ip"

  # ⚠️ все четыре параметра обязательны
  service_uid           = nubes_postgres.db.id # UUID сервиса-получателя
  from_service_namespace = "my-namespace"
  internal_addr_access   = "10.0.0.15"
  ip_space_name          = "internet-no-antiddos-v1"

  dnat_create           = true
  internal_port_access  = "5432"
  snat_create           = false
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса |
| `service_uid` | string | **да** | — | UUID сервиса, к которому привязывается IP |
| `from_service_namespace` | string | **да** | — | Namespace составной услуги (для отображения правил FW) |
| `internal_addr_access` | string | **да** | — | Внутренний адрес для связки NAT (требует `dnat_create` или `snat_create`) |
| `ip_space_name` | string | **да** | — | Тип IP-space: `internet-ipv4-v1` / `internet-antiddos-v1` / `internet-no-antiddos-v1` |
| `dnat_create` | bool | нет | `true` | Создать DNAT-правило |
| `internal_port_access` | string | нет | `5432` | Порт для связки NAT; несколько значений — через пробел, `all` — все порты (требует `dnat_create=true`) |
| `snat_create` | bool | нет | — | Создать SNAT-правило |
| `resource_realm` | string | нет | `ngcloud.ru` | Платформа для развёртывания |
| `from_service_cloud_edge_name` | string | нет | — | Edge составной услуги (если realm связан с CloudDirector) |
| `from_service_cloud_vdc_name` | string | нет | — | vDC составной услуги |
| `from_service_cloud_org_name` | string | нет | — | Org составной услуги |
| `from_service_cloud_vmware_url` | string | нет | — | API-URL составной услуги |
| `from_service_cloud_edge_scope` | string | нет | — | Edge scope составной услуги |
| `from_service_vdc_group_name` | string | нет | — | Имя группы vDC составной услуги |

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
| `state_out` / `state_out_flat` | map | Ответ API (внешний адрес, `instanceUid`) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

`terraform destroy` освобождает публичный IP (удаление в облаке и из Terraform-state).
Усыпление (`suspend`) у этого сервиса не поддерживается — атрибута `suspend_on_destroy` нет.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — ресурс остаётся в облаке (`state_only`); удаление в облаке не выполняется |

- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Резервирование IP и создание правил NAT |
| `modify` | ✅ `apply` | Меняются `service_uid`, `ip_space_name`, `dnat_create`, `snat_create`, `internal_port_access`, `internal_addr_access` |
| `delete` | ✅ `destroy` | Освобождение IP |


See the [Terraform Registry attribute reference](../reference/resources/vcexternalip.md).
