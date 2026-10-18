# `nubes_nextcloud_backup` Resource

Service: Nextcloud.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/nextcloud/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/nextcloud/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `nextcloud_id` | `string` | Required | No | — |
| `operation_timeout` | `string` | Optional | No | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
| `startup_configuration` | `string` | Required | No | Не изменяется после создания. Изменение требует пересоздания. |
