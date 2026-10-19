---
page_title: "nubes_kafka"
description: "nubes_kafka TEST resource schema for Nubes Cloud."
subcategory: "Очереди"
---

# `nubes_kafka` Resource

Service: ApacheKafka.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/kafka/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_configuration` | `object` | Required | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `cluster_configuration` | `object` | Required | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `resource_name` | `string` | Required | — |
| `startup_configuration` | `object` | Required | — |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `allocation_external` | `bool` | Optional, Computed | — |
| `master_access_list` | `string` | Optional, Computed | Необходимо настраивать, когда зарезервирован внешний адрес<br/>Если передан пустой массив, доступ выделяется всем |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `disk` | `number` | Optional, Computed | — |
| `memory` | `number` | Optional, Computed | — |
| `replicas` | `number` | Required | — |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `resource_realm` | `string` | Required | Кластер Kubernetes, на котором будет развернут экземпляр |
