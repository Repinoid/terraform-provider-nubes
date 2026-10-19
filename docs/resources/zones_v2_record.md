---
page_title: "nubes_zones_v2_record"
description: "nubes_zones_v2_record TEST resource schema for Nubes Cloud."
subcategory: "Сеть"
---

# `nubes_zones_v2_record` Resource

Service: Управление DNS.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/zones_v2/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `record_input` | `string` | Optional, Computed | Валидация производится только после запуска операции. Если при указанном типе будет указано некорректно значение, это будет описано в логе ошибки операции |
| `record_name` | `string` | Required | Имя зоны в конце проставляется автоматически Изменение требует пересоздания. |
| `record_ttl` | `number` | Optional, Computed | Time-to-Leave записи. Раз в какое время запись будет обновлена |
| `record_type` | `string` | Optional, Computed | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
| `zones_v2_id` | `string` | Required | — |
