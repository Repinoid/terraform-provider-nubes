---
page_title: "nubes_zones_v2_record Resource"
description: "Service: Управление DNS."
subcategory: "Сеть"
---

# `nubes_zones_v2_record` Resource

Service: Управление DNS.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/zones_v2/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/zones_v2/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `record_input` | `string` | Optional, Computed | No | Валидация производится только после запуска операции. Если при указанном типе будет указано некорректно значение, это будет описано в логе ошибки операции |
| `record_name` | `string` | Required | No | Имя зоны в конце проставляется автоматически Изменение требует пересоздания. |
| `record_ttl` | `number` | Optional, Computed | No | Time-to-Leave записи. Раз в какое время запись будет обновлена |
| `record_type` | `string` | Optional, Computed | No | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
| `zones_v2_id` | `string` | Required | No | — |
