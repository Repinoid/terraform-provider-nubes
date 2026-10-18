---
page_title: "nubes_vapp Resource"
description: "Service: Виртуальный каталог ВМ (vApp)."
subcategory: "VMware"
---

# `nubes_vapp` Resource

Service: Виртуальный каталог ВМ (vApp).

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/vapp/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/vapp/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `nsxt_uid` | `string` | Optional | No | UUID Услуги `Сетевой шлюз периметра (Edge)` Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | No | — |
| `resource_name` | `string` | Required | No | — |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vapp_name` | `string` | Optional, Computed | No | Уникальное в рамках организации Не изменяется после создания. |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
| `vdc_uid` | `string` | Optional | No | UUID Услуги `Виртуальный датацентр (vDC)` Не изменяется после создания. |
