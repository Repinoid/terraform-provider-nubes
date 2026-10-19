---
page_title: "nubes_vc_mgmt_sthutrval_cluster"
description: "nubes_vc_mgmt_sthutrval_cluster TEST resource schema for Nubes Cloud."
subcategory: "Other TEST resources"
---

# `nubes_vc_mgmt_sthutrval_cluster` Resource

TEST API resource. Its schema is shown below; service-level narrative documentation is not currently available.

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_configuration` | `object` | Required | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `cluster_configuration` | `object` | Required | — |
| `control_plane_configuration` | `object` | Required | — |
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
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
| `worker_configuration` | `string` | Optional, Computed | Не изменяется после создания. |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_ip_list_api` | `string` | Optional, Computed | Access IP List для подключения к Kubernetes API. Указывается массив. Если передан пустой массив, доступ выделяется всем<br/>Пример: **[&quot;0.0.0.0/0&quot;]** или **[]** или **[&quot;178.4.13.2&quot;, &quot;178.4.14.0/24&quot;]** |
| `access_ip_list_ingress` | `string` | Optional, Computed | Access IP List для подключения к Ingress. Указывается массив. Если передан пустой массив, доступ выделяется всем<br/>Пример: **[&quot;0.0.0.0/0&quot;]** или **[]** или **[&quot;178.4.13.2&quot;, &quot;178.4.14.0/24&quot;]** |
| `need_external_address_api` | `bool` | Optional, Computed | Выделяет внешний IP для доступа к kubernetes API. Пока нельзя выставить false |
| `need_external_address_ingress` | `bool` | Optional, Computed | Выделяет внешний IP для доступа к Ingress |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `app_version` | `string` | Optional, Computed | — |

### Nested attributes: `control_plane_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `count` | `number` | Optional, Computed | — |
| `sizing_disk` | `number` | Optional, Computed | — |
| `sizing_policy` | `string` | Required | Доступные sizing Policy представлены исходя из ресурсной платформы, которая была выбрана |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cluster_name` | `string` | Optional, Computed | — |
| `network_provider` | `string` | Required | — |
| `provider_vdc` | `string` | Required | Приписка `2.8` - Скорость vCPU - 2800MHz<br/>И так далее...<br/> |
| `resource_realm` | `string` | Optional, Computed | Облако в котором будет развёрнут тенант |
