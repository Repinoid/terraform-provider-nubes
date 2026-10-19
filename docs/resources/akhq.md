---
page_title: "nubes_akhq"
description: "nubes_akhq TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_akhq` Resource

Service: Akhq.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/akhq/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `domain` | `string` | Required | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `kafka_uid` | `string` | Optional | Выбрать кластер из списка |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `resource_cpu` | `number` | Optional, Computed | Пример: `1000` Не изменяется после создания. |
| `resource_instances` | `number` | Optional, Computed | Не изменяется после создания. |
| `resource_memory` | `number` | Optional, Computed | Пример: `1024` Не изменяется после создания. |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Required | Выбрать платформу из списка Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
