---
page_title: "nubes_k8s_velero"
description: "nubes_k8s_velero TEST resource schema for Nubes Cloud."
subcategory: "K8s"
---

# `nubes_k8s_velero` Resource

Service: Velero.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/k8s_velero/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `backup_configuration` | `object` | Required | — |
| `cluster_configuration` | `object` | Required | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `node_agent_configuration` | `object` | Required | — |
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

### Nested attributes: `backup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `s3_bucket_uid` | `string` | Required | Экземпляр S3-Бакета обязателен<br/>Учётная запись будет создана автоматически, если её нету |

### Nested attributes: `cluster_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `memory` | `number` | Optional, Computed | Указывается в Megabytes |
| `replicas` | `number` | Optional, Computed | — |

### Nested attributes: `node_agent_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `cpu` | `number` | Optional, Computed | Указывается в Milicores<br/>`1000` Mili == `1` Ядро |
| `memory` | `number` | Optional, Computed | Указывается в Megabytes |

### Nested attributes: `startup_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `resource_realm` | `string` | Required | — |
