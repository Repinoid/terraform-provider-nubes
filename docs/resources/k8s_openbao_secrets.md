---
page_title: "nubes_k8s_openbao_secrets"
description: "nubes_k8s_openbao_secrets TEST resource schema for Nubes Cloud."
subcategory: "K8s"
---

# `nubes_k8s_openbao_secrets` Resource

Service: Vault.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/k8s_openbao/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `k8s_openbao_id` | `string` | Required | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `kv_name` | `string` | Required | Выберите из списка Не изменяется после создания. Изменение требует пересоздания. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `secret_data` | `string` | Required | Хотя бы один секрет необходимо добавить Не изменяется после создания. Изменение требует пересоздания. |
| `secret_path` | `string` | Optional, Computed | Пример: `path` или `road/to/path` Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
