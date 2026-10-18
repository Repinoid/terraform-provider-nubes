---
page_title: "nubes_llm_ai_secrets Resource"
description: "Service: ML Улей."
subcategory: "Приложения"
---

# `nubes_llm_ai_secrets` Resource

Service: ML Улей.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/llm_ai/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/llm_ai/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `key_name` | `string` | Required | No | Разрешены только латинские буквы и цифры без пробелов Не изменяется после создания. Изменение требует пересоздания. |
| `llm_ai_id` | `string` | Required | No | — |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `models` | `string` | Required | No | Не изменяется после создания. Изменение требует пересоздания. |
| `operation_timeout` | `string` | Optional | No | — |
| `rpm_limit` | `number` | Optional, Computed | No | Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
| `team_name` | `string` | Required | No | Не изменяется после создания. Изменение требует пересоздания. |
| `tpm_limit` | `number` | Optional, Computed | No | Не изменяется после создания. Изменение требует пересоздания. |
