---
page_title: "nubes_k8s_openbao Resource"
description: "Service: Vault."
subcategory: "K8s"
---

# `nubes_k8s_openbao` Resource

Service: Vault.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/k8s_openbao/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/k8s_openbao/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `app_version` | `string` | Optional, Computed | No | Выбрать из списка |
| `domain` | `string` | Optional, Computed | No | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_cpu` | `number` | Optional, Computed | No | Пример: `500` |
| `resource_disk` | `number` | Optional, Computed | No | Пример: `10` |
| `resource_instances` | `number` | Optional, Computed | No | Пример: `3` |
| `resource_memory` | `number` | Optional, Computed | No | Пример: `2048` |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Кластер Kubernetes в котором будет развёрнуто приложение Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
