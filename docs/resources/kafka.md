---
page_title: "nubes_kafka Resource"
description: "Service: ApacheKafka."
subcategory: "Очереди"
---

# `nubes_kafka` Resource

Service: ApacheKafka.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/kafka/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/kafka/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `cluster_configuration` | `object` | Required | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `resource_name` | `string` | Required | No | — |
| `startup_configuration` | `object` | Required | No | — |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `allocation_external` | `bool` | Optional, Computed | No | — |
| `master_access_list` | `string` | Optional, Computed | No | Необходимо настраивать, когда зарезервирован внешний адрес<br/>Если передан пустой массив, доступ выделяется всем |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cpu` | `number` | Optional, Computed | No | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | No | — |
| `memory` | `number` | Optional, Computed | No | — |
| `replicas` | `number` | Required | No | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `resource_realm` | `string` | Optional, Computed | No | Кластер Kubernetes, на котором будет развернут экземпляр |
