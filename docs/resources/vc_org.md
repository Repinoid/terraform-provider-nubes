---
page_title: "nubes_vc_org"
description: "nubes_vc_org TEST resource schema for Nubes Cloud."
subcategory: "VMware"
---

# `nubes_vc_org` Resource

Service: Организация в Cloud Director.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/vc_org/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `organization_suffix` | `string` | Optional, Computed | Не изменяется после создания. |
| `organization_type` | `string` | Optional, Computed | Доступ к организации предоставляется только при выборе `iaas` Не изменяется после создания. |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Optional, Computed | Облако в котором будет развёрнут тенант Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `v_ip_configure` | `string` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
