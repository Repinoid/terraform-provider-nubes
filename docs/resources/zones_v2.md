---
page_title: "nubes_zones_v2"
description: "nubes_zones_v2 TEST resource schema for Nubes Cloud."
subcategory: "Сеть"
---

# `nubes_zones_v2` Resource

Service: Управление DNS.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/zones_v2/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `records` | `string` | Optional, Computed | — |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Optional, Computed | Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
| `zone_expire` | `number` | Optional, Computed | Время в течение которого вторичный DNS будет пытаться завершить синхронизацию зоны с первичным Не изменяется после создания. |
| `zone_kind` | `string` | Optional, Computed | Тип зоны Не изменяется после создания. |
| `zone_masters_for_slave` | `string` | Optional, Computed | Использовать если выбран тип SLAVE. IP адреса в формате через пробел для настройки Master-серверов Не изменяется после создания. |
| `zone_name` | `string` | Required | Имя зоны Не изменяется после создания. |
| `zone_refresh` | `number` | Optional, Computed | Время ожидания ответа вторичного DNS перед запросом SOA-записи с первичных серверов Не изменяется после создания. |
| `zone_retry` | `number` | Optional, Computed | Время в течение которого сервер повторно попробует синхронизировать информацию с первичных DNS-серверов Не изменяется после создания. |
| `zone_ttl` | `number` | Optional, Computed | Время, в течение которого информация будет кешироваться другими DNS-серверами Не изменяется после создания. |
