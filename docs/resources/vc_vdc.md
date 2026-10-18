---
page_title: "nubes_vc_vdc Resource"
description: "Service: Виртуальный датацентр (vDC)."
subcategory: "VMware"
---

# `nubes_vc_vdc` Resource

Service: Виртуальный датацентр (vDC).

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/vc_vdc/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/vc_vdc/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `cpu_allocated` | `number` | Optional, Computed | No | шт |
| `cpu_guaranteed` | `number` | Optional, Computed | No | % резервирования Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mem_allocated` | `number` | Optional, Computed | No | gb |
| `network_provider` | `string` | Required | No | Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | No | — |
| `organization_uid` | `string` | Optional | No | Для создания VDC предварительно необходимо создать организацию через услугу `Организация в Cloud Director` Не изменяется после создания. |
| `provider_vdc` | `string` | Required | No | Приписка `2.8` - Скорость vCPU - 2800MHz<br/>И так далее...<br/> Не изменяется после создания. |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | No | — |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `storage_config` | `string` | Optional, Computed | No | Первая политика в списке указывается как &quot;По умолчанию&quot; |
| `suffix_name` | `string` | Optional, Computed | No | Если вписать произвольное значение, оно будет добавлено в конец имени VDC<br/>**Нельзя изменить после создания** Не изменяется после создания. |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
