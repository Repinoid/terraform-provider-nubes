---
page_title: "nubes_dnsrecord"
description: "nubes_dnsrecord TEST resource schema for Nubes Cloud."
subcategory: "Сеть"
---

# `nubes_dnsrecord` Resource

Service: DNS запись.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/dnsrecord/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `record_input` | `string` | Required | No | Значение. Необходима валидация, в зависимости от типа записи. Должно быть заменено на UUID-внешнийIP |
| `record_name` | `string` | Required | No | Имя записи. Имя зоны в конце проставляется автоматически Не изменяется после создания. |
| `record_ttl` | `number` | Optional, Computed | No | Time-to-Leave записи. Раз в какое время запись будет обновлена Не изменяется после создания. |
| `record_type` | `string` | Required | No | Тип DNS-записи Не изменяется после создания. |
| `resource_name` | `string` | Required | No | — |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
| `zone_uid` | `string` | Optional | No | UUID экземпляра услуги DNS-зоны Не изменяется после создания. |
