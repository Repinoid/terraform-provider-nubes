---
page_title: "nubes_zones_v2"
description: "nubes_zones_v2 TEST resource schema for Nubes Cloud."
subcategory: "Сеть"
---

# `nubes_zones_v2` Resource

Service: Управление DNS.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/zones_v2/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `records` | `string` | Optional, Computed | No | — |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
| `zone_expire` | `number` | Optional, Computed | No | Время в течение которого вторичный DNS будет пытаться завершить синхронизацию зоны с первичным Не изменяется после создания. |
| `zone_kind` | `string` | Optional, Computed | No | Тип зоны Не изменяется после создания. |
| `zone_masters_for_slave` | `string` | Optional, Computed | No | Использовать если выбран тип SLAVE. IP адреса в формате через пробел для настройки Master-серверов Не изменяется после создания. |
| `zone_name` | `string` | Required | No | Имя зоны Не изменяется после создания. |
| `zone_refresh` | `number` | Optional, Computed | No | Время ожидания ответа вторичного DNS перед запросом SOA-записи с первичных серверов Не изменяется после создания. |
| `zone_retry` | `number` | Optional, Computed | No | Время в течение которого сервер повторно попробует синхронизировать информацию с первичных DNS-серверов Не изменяется после создания. |
| `zone_ttl` | `number` | Optional, Computed | No | Время, в течение которого информация будет кешироваться другими DNS-серверами Не изменяется после создания. |
