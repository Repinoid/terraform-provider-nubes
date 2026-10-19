---
page_title: "nubes_k8s_sthutrval_cluster"
description: "nubes_k8s_sthutrval_cluster TEST resource schema for Nubes Cloud."
subcategory: "K8s"
---

# `nubes_k8s_sthutrval_cluster` Resource

Service: Kubernetes кластер Штурвал.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/k8s_sthutrval_cluster/).

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
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
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
| `worker_configuration` | `string` | Optional, Computed | — |

### Nested attributes: `access_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_ip_list_api` | `string` | Optional, Computed | Access IP List для подключения к Kubernetes API. Указывается массив. Если передан пустой массив, доступ выделяется всем |
| `access_ip_list_ingress` | `string` | Optional, Computed | Access IP List для подключения к Ingress. Указывается массив. Если передан пустой массив, доступ выделяется всем |
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
| `sizing_disk` | `number` | Optional, Computed | Дисковая политика берётся из услуги VDC, первый в списке |
| `sizing_policy` | `string` | Optional, Computed | Доступные sizing Policy можно обновлять через услугу `Виртуальный датацентр` |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cluster_name` | `string` | Optional, Computed | — |
| `ex_ingress` | `bool` | Optional, Computed | — |
| `ex_local_csi` | `bool` | Optional, Computed | — |
| `ex_logging` | `bool` | Optional, Computed | А также перенаправление логов в Loki<br/>Без этого логи сервисов в ЛК отображаться не будут |
| `ex_monitoring` | `bool` | Optional, Computed | А также перенаправление логов в VictoriaMetrics<br/>Без этого метрики сервисов в ЛК отображаться не будут |
| `ex_named_csi` | `bool` | Optional, Computed | — |
| `ex_update` | `bool` | Optional, Computed | — |
| `ex_vip` | `bool` | Optional, Computed | — |
| `nsxt_uid` | `string` | Required | — |
| `vdc_uid` | `string` | Required | — |
