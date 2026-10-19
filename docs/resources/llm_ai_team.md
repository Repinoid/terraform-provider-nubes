---
page_title: "nubes_llm_ai_team"
description: "nubes_llm_ai_team TEST resource schema for Nubes Cloud."
subcategory: "Приложения"
---

# `nubes_llm_ai_team` Resource

Service: ML Улей.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/llm_ai/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `llm_ai_id` | `string` | Required | — |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `models` | `string` | Required | — |
| `operation_timeout` | `string` | Optional | — |
| `rpm_limit` | `number` | Optional, Computed | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
| `team_name` | `string` | Required | Разрешены только латинские буквы и цифры без пробелов Изменение требует пересоздания. |
| `tpm_limit` | `number` | Optional, Computed | — |
