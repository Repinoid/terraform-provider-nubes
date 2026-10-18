# `nubes_dnsrecord` — DNS запись

DNS-запись внутри существующей зоны (`nubes_zones_v2`). Сервис создаёт запись
заданного типа (A, CNAME, MX, TXT и др.) и отдаёт итоговый FQDN и адреса NS-серверов.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/dnsrecord.md).
>
> **Специфика Nubes: `destroy` для этого ресурса **удаляет** запись**
> `terraform destroy` удаляет запись из облака и из Terraform-state. Усыпление
> (`suspend`) у этого сервиса не поддерживается.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_zones_v2" "zone" {
  resource_name = "my-zone"
  zone_name     = "example.com"
  records       = jsonencode([])
}

resource "nubes_dnsrecord" "www" {
  resource_name = "my-www-record"

  # ⚠️ обязателен: UUID существующей DNS-зоны (сервис 109)
  zone_uid = nubes_zones_v2.zone.id

  record_name  = "www"
  record_type  = "A"
  record_input = "10.0.0.1"
  record_ttl   = 3600 # 30…604800
}
```

> ⚠️ `zone_uid` в схеме помечен как необязательный, однако без него создание не
> запускается — указывайте UUID существующей DNS-зоны.

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `record_type` | string | **да** | — | Тип DNS-записи | `A`, `AAAA`, `CNAME`, `MX`, `PTR`, `SOA`, `SRV`, `TXT`, `LUA`, `NS`, `SSHFP` |
| `record_name` | string | **да** | — | Имя записи; имя зоны добавляется автоматически | — |
| `record_input` | string | **да** | — | Значение записи (формат зависит от типа) | валидируется по типу записи |
| `zone_uid` | string (uuid) | нет *(фактически да)* | — | UUID экземпляра DNS-зоны (сервис 109) | ссылка на `nubes_zones_v2` |
| `record_ttl` | integer | нет | `3600` | Time-to-live записи (сек) | 30…604800 |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `5m`) |
| `log_level` | string | нет | `none` | Логирование стадий операции: `none`, `info`, `debug` |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply (=import) |

---

## Attribute Reference

Вычисляемые атрибуты (только чтение):

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.record` | string | Итоговый FQDN записи: `<record_name>.<zone_name>` |
| `state_out.dnsServers` | string (json) | Адреса NS-серверов зоны |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_params.recordName` / `recordType` / `recordInput` / `recordTTL` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault (обычно пусто) |
| `vault_fields` | list | `[]` |

```hcl
output "www_fqdn" { value = nubes_dnsrecord.www.state_out_flat["record"] }
```

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

`terraform destroy` удаляет запись из облака и из Terraform-state. Усыпление (`suspend`)
у этого сервиса не поддерживается — атрибута `suspend_on_destroy` нет.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — запись остаётся в облаке (`state_only`); удаление в облаке не выполняется |

- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание записи |
| `modify` | ✅ `apply` | Изменяется только `record_input` |
| `delete` | ✅ `destroy` | Удаление записи |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/dnsrecord.md).
