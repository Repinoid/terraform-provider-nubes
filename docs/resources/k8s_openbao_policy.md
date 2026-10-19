---
page_title: "nubes_k8s_openbao_policy"
description: "nubes_k8s_openbao_policy TEST resource schema for Nubes Cloud."
subcategory: "K8s"
---

# `nubes_k8s_openbao_policy` Resource

Service: Vault.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/k8s_openbao/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `k8s_openbao_id` | `string` | Required | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `policy_name` | `string` | Optional, Computed | Сделано для дальнейшего удобного использования Не изменяется после создания. Изменение требует пересоздания. |
| `policy_settings` | `string` | Required | Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
