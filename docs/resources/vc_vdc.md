---
page_title: "nubes_vc_vdc"
description: "nubes_vc_vdc TEST resource schema for Nubes Cloud."
subcategory: "VMware"
---

# `nubes_vc_vdc` Resource

Service: Виртуальный датацентр (vDC).

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/vc_vdc/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `cpu_allocated` | `number` | Optional, Computed | шт |
| `cpu_guaranteed` | `number` | Optional, Computed | % резервирования Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mem_allocated` | `number` | Optional, Computed | gb |
| `network_provider` | `string` | Required | Не изменяется после создания. |
| `operation_timeout` | `string` | Optional | — |
| `organization_uid` | `string` | Optional | Для создания VDC предварительно необходимо создать организацию через услугу `Организация в Cloud Director` Не изменяется после создания. |
| `provider_vdc` | `string` | Required | Приписка `2.8` - Скорость vCPU - 2800MHz<br/>И так далее...<br/> Не изменяется после создания. |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | — |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `storage_config` | `string` | Optional, Computed | Первая политика в списке указывается как &quot;По умолчанию&quot; |
| `suffix_name` | `string` | Optional, Computed | Если вписать произвольное значение, оно будет добавлено в конец имени VDC<br/>**Нельзя изменить после создания** Не изменяется после создания. |
| `suspend_on_destroy` | `bool` | Optional, Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
