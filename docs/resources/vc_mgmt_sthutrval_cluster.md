---
page_title: "nubes_vc_mgmt_sthutrval_cluster Resource"
description: "## Schema"
subcategory: "General"
---

# `nubes_vc_mgmt_sthutrval_cluster` Resource

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_configuration` | `object` | Required | No | — |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `cluster_configuration` | `object` | Required | No | — |
| `control_plane_configuration` | `object` | Required | No | — |
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
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
| `worker_configuration` | `string` | Optional, Computed | No | Не изменяется после создания. |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_ip_list_api` | `string` | Optional, Computed | No | Access IP List для подключения к Kubernetes API. Указывается массив. Если передан пустой массив, доступ выделяется всем<br/>Пример: **[&quot;0.0.0.0/0&quot;]** или **[]** или **[&quot;178.4.13.2&quot;, &quot;178.4.14.0/24&quot;]** |
| `access_ip_list_ingress` | `string` | Optional, Computed | No | Access IP List для подключения к Ingress. Указывается массив. Если передан пустой массив, доступ выделяется всем<br/>Пример: **[&quot;0.0.0.0/0&quot;]** или **[]** или **[&quot;178.4.13.2&quot;, &quot;178.4.14.0/24&quot;]** |
| `need_external_address_api` | `bool` | Optional, Computed | No | Выделяет внешний IP для доступа к kubernetes API. Пока нельзя выставить false |
| `need_external_address_ingress` | `bool` | Optional, Computed | No | Выделяет внешний IP для доступа к Ingress |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `app_version` | `string` | Optional, Computed | No | — |

### Nested attributes: `control_plane_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `count` | `number` | Optional, Computed | No | — |
| `sizing_disk` | `number` | Optional, Computed | No | — |
| `sizing_policy` | `string` | Required | No | Доступные sizing Policy представлены исходя из ресурсной платформы, которая была выбрана |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cluster_name` | `string` | Optional, Computed | No | — |
| `network_provider` | `string` | Required | No | — |
| `provider_vdc` | `string` | Required | No | Приписка `2.8` - Скорость vCPU - 2800MHz<br/>И так далее...<br/> |
| `resource_realm` | `string` | Optional, Computed | No | Облако в котором будет развёрнут тенант |
