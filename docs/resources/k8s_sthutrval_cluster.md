---
page_title: "nubes_k8s_sthutrval_cluster"
description: "nubes_k8s_sthutrval_cluster TEST resource schema for Nubes Cloud."
subcategory: "K8s"
---

# `nubes_k8s_sthutrval_cluster` Resource

Service: Kubernetes кластер Штурвал.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/k8s_sthutrval_cluster/).

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
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
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
| `worker_configuration` | `string` | Optional, Computed | No | — |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_ip_list_api` | `string` | Optional, Computed | No | Access IP List для подключения к Kubernetes API. Указывается массив. Если передан пустой массив, доступ выделяется всем |
| `access_ip_list_ingress` | `string` | Optional, Computed | No | Access IP List для подключения к Ingress. Указывается массив. Если передан пустой массив, доступ выделяется всем |
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
| `sizing_disk` | `number` | Optional, Computed | No | Дисковая политика берётся из услуги VDC, первый в списке |
| `sizing_policy` | `string` | Optional, Computed | No | Доступные sizing Policy можно обновлять через услугу `Виртуальный датацентр` |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `cluster_name` | `string` | Optional, Computed | No | — |
| `ex_ingress` | `bool` | Optional, Computed | No | — |
| `ex_local_csi` | `bool` | Optional, Computed | No | — |
| `ex_logging` | `bool` | Optional, Computed | No | А также перенаправление логов в Loki<br/>Без этого логи сервисов в ЛК отображаться не будут |
| `ex_monitoring` | `bool` | Optional, Computed | No | А также перенаправление логов в VictoriaMetrics<br/>Без этого метрики сервисов в ЛК отображаться не будут |
| `ex_named_csi` | `bool` | Optional, Computed | No | — |
| `ex_update` | `bool` | Optional, Computed | No | — |
| `ex_vip` | `bool` | Optional, Computed | No | — |
| `nsxt_uid` | `string` | Required | No | — |
| `vdc_uid` | `string` | Required | No | — |
