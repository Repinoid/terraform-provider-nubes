# `nubes_zones_v2` — Управление DNS

Управляет DNS-зоной домена в инфраструктуре `ngcloud.ru`: создаёт зону, поддерживает
записи, удаляет зону и синхронизирует её с фактическим состоянием авторитетных
DNS-серверов (`ns3.ngcloud.ru`).

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/zones_v2.md).
>
> **Специфика Nubes: `destroy` для этого ресурса **удаляет** зону**
> `terraform destroy` удаляет зону из облака и из Terraform-state. Усыпление (`suspend`)
> у этого сервиса не поддерживается. Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в
> общем руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_zones_v2" "zone" {
  resource_name = "my-zone"
  zone_name     = "example.com"
  records       = jsonencode([]) # ⚠️ обязателен

  # resource_realm = "powerdns"   # дефолт (фиксирован)
  zone_kind    = "NATIVE" # NATIVE | SLAVE
  zone_ttl     = 86400
  zone_refresh = 86400
  zone_retry   = 7200
  zone_expire  = 3600000
}

resource "nubes_zones_v2_record" "www" {
  zones_v2_id  = nubes_zones_v2.zone.id
  record_name  = "www"
  record_type  = "A"
  record_input = jsonencode([{ record = "10.0.0.1" }])
  record_ttl   = 3600
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `zone_name` | string | **да** | — | Имя зоны | длина 1…64 |
| `records` | string (json) | **да** | — | Записи зоны (JSON-массив) | ⚠️ `[]` нормализуется платформой |
| `zone_kind` | string | нет | `NATIVE` | Тип зоны | `NATIVE`, `SLAVE` |
| `zone_masters_for_slave` | string | нет | — | IP мастер-серверов через пробел (для `SLAVE`) | — |
| `zone_ttl` | integer | нет | `86400` | Время кеширования записи другими DNS | ≥ 600 |
| `zone_refresh` | integer | нет | `86400` | Интервал проверки SOA на первичных серверах | ≥ 600 |
| `zone_retry` | integer | нет | `7200` | Интервал повторной попытки синхронизации | ≥ 600 |
| `zone_expire` | integer | нет | `3600000` | Время, за которое вторичный DNS должен завершить синхронизацию | ≥ 3600 |
| `resource_realm` | string | нет | `powerdns` | «Платформа» — DNS-сервис (фиксировано) | `powerdns` |

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
| `state_out.dnsServers` | string (json) | Адреса NS-серверов зоны: `["185.247.187.83", "81.22.46.43"]` |
| `state_out.zoneSerial` | string | Серийный номер зоны (SOA) |
| `state_out.zoneSoaServer` | string | SOA-сервер: `ns.ngcloud.ru` |
| `state_params.records` | string (json) | Записи зоны |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault (обычно пусто) |
| `vault_fields` | list | Список доступных ключей секретов |

```hcl
output "ns_servers" {
  value = jsondecode(nubes_zones_v2.zone.state_out_flat["dnsServers"])
}
```

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

`terraform destroy` удаляет зону из облака и из Terraform-state. Усыпление (`suspend`)
у этого сервиса не поддерживается — атрибута `suspend_on_destroy` нет.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — зона остаётся в облаке (`state_only`); удаление в облаке не выполняется |


---

## Подресурсы

### `nubes_zones_v2_record`

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `zones_v2_id` | string (uuid) | **да** | — | UUID DNS-зоны | — |
| `record_name` | string | **да** | — | Имя записи (без имени зоны) | длина 1…64 |
| `record_type` | string | нет | `A` | Тип записи | `A`, `AAAA`, `CNAME`, `MX`, `SRV`, `TXT`, `LUA`, `NS`, `SSHFP` |
| `record_ttl` | integer | нет | `3600` | Время жизни записи (сек) | 30…84000 |
| `record_input` | map (список) | нет | `[{}]` | Значение записи | `record` — **да** |

Идентификатор: `<zone_uid>:record:<hash>`.

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание зоны |
| `modify` | ✅ `apply` | Меняется `records` |
| `delete` | ✅ `destroy` | Удаление зоны |
| `create_record` / `modify_record` / `delete_record` | ✅ через `nubes_zones_v2_record` | Управление записями |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/zones_v2.md).
