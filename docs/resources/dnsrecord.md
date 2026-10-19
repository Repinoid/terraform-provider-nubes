---
page_title: "nubes_dnsrecord"
description: "nubes_dnsrecord TEST resource schema for Nubes Cloud."
subcategory: "Сеть"
---

# `nubes_dnsrecord` Resource

Service: DNS запись.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/dnsrecord/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `record_input` | `string` | Required | Значение. Необходима валидация, в зависимости от типа записи. Должно быть заменено на UUID-внешнийIP |
| `record_name` | `string` | Required | Имя записи. Имя зоны в конце проставляется автоматически Не изменяется после создания. |
| `record_ttl` | `number` | Optional, Computed | Time-to-Leave записи. Раз в какое время запись будет обновлена Не изменяется после создания. |
| `record_type` | `string` | Required | Тип DNS-записи Не изменяется после создания. |
| `resource_name` | `string` | Required | — |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
| `zone_uid` | `string` | Optional | UUID экземпляра услуги DNS-зоны Не изменяется после создания. |
