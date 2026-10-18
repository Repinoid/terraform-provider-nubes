---
page_title: "nubes_vc_org_user Resource"
description: "Service: Организация в Cloud Director."
subcategory: "VMware"
---

# `nubes_vc_org_user` Resource

Service: Организация в Cloud Director.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/vc_org/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/vc_org/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `admins` | `string` | Optional, Computed | No | Не изменяется после создания. Изменение требует пересоздания. |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
| `vc_org_id` | `string` | Required | No | — |
