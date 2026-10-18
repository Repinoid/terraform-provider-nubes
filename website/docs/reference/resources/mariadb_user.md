# `nubes_mariadb_user` Resource

Service: Mariadb.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/mariadb/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/mariadb/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_hosts` | `string` | Optional, Computed | No | Пример: `%` Не изменяется после создания. Изменение требует пересоздания. |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `db_name` | `string` | Optional, Computed | No | Необходимо указать массив существующих баз Не изменяется после создания. Изменение требует пересоздания. |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mariadb_id` | `string` | Required | No | — |
| `operation_timeout` | `string` | Optional | No | — |
| `password` | `string` | Optional, Computed | Yes | Секрет будет записан в хранилище секретов и не будет отображаться публично Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
| `username` | `string` | Optional, Computed | No | Пример: `myusername` Не изменяется после создания. Изменение требует пересоздания. |
