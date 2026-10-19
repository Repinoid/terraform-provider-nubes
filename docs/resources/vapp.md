---
page_title: "nubes_vapp"
description: "nubes_vapp TEST resource schema for Nubes Cloud."
subcategory: "VMware"
---

# `nubes_vapp` Resource

Service: Виртуальный каталог ВМ (vApp).

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/vapp/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `nsxt_uid` | `string` | Optional | UUID Услуги `Сетевой шлюз периметра (Edge)` Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | — |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | — |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vapp_name` | `string` | Optional, Computed | Уникальное в рамках организации Не изменяется после создания. |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
| `vdc_uid` | `string` | Optional | UUID Услуги `Виртуальный датацентр (vDC)` Не изменяется после создания. |
