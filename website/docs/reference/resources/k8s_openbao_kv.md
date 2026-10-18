# `nubes_k8s_openbao_kv` Resource

Service: Vault.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/k8s_openbao/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/k8s_openbao/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `k8s_openbao_id` | `string` | Required | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `kv_name` | `string` | Optional, Computed | No | Будет создан `v2 тип движка`<br/><br/>Уникальное в рамках услуги Не изменяется после создания. Изменение требует пересоздания. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
