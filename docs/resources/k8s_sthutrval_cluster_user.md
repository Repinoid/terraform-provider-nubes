---
page_title: "nubes_k8s_sthutrval_cluster_user"
description: "nubes_k8s_sthutrval_cluster_user TEST resource schema for Nubes Cloud."
subcategory: "K8s"
---

# `nubes_k8s_sthutrval_cluster_user` Resource

Service: Kubernetes кластер Штурвал.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/k8s_sthutrval_cluster/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `admins` | `string` | Optional, Computed | No | Не изменяется после создания. Изменение требует пересоздания. |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `k8s_sthutrval_cluster_id` | `string` | Required | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
