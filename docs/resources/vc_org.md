---
page_title: "nubes_vc_org Resource"
description: "Service: Организация в Cloud Director."
subcategory: "VMware"
---

# `nubes_vc_org` Resource

Service: Организация в Cloud Director.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/vc_org/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/vc_org/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `organization_type` | `string` | Optional, Computed | No | Доступ к организации предоставляется только при выборе `iaas` Не изменяется после создания. |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Облако в котором будет развёрнут тенант Не изменяется после создания. |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `v_ip_configure` | `string` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
