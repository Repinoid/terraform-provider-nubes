# `nubes_k8s_openbao_secrets` Resource

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
| `kv_name` | `string` | Required | No | Выберите из списка Не изменяется после создания. Изменение требует пересоздания. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `secret_data` | `string` | Required | No | Хотя бы один секрет необходимо добавить Не изменяется после создания. Изменение требует пересоздания. |
| `secret_path` | `string` | Optional, Computed | No | Пример: `path` или `road/to/path` Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
