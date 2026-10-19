---
page_title: "nubes_k8s_openbao"
description: "nubes_k8s_openbao TEST resource schema for Nubes Cloud."
subcategory: "K8s"
---

# `nubes_k8s_openbao` Resource

Service: Vault.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/k8s_openbao/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `app_version` | `string` | Optional, Computed | Выбрать из списка |
| `domain` | `string` | Optional, Computed | Если будет указан FQDN, то происходят проверки DNS на валидность А-записи, или будет добавлено в зону resourceRealm Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_cpu` | `number` | Optional, Computed | Пример: `500` |
| `resource_disk` | `number` | Optional, Computed | Пример: `10` |
| `resource_instances` | `number` | Optional, Computed | Пример: `3` |
| `resource_memory` | `number` | Optional, Computed | Пример: `2048` |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Required | Кластер Kubernetes в котором будет развёрнуто приложение Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
