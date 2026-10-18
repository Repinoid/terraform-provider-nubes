---
page_title: "nubes_postgres_user Resource"
description: "Service: PostgreSQL."
subcategory: "Базы данных"
---

# `nubes_postgres_user` Resource

Service: PostgreSQL.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/postgres/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/postgres/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mtls_access` | `bool` | Optional, Computed | No | Подключение пользователя без ввода паролей<br/>Политика применяется **ко всем базам** внутри Postgresql<br/><br/>RBAC на базы настраивать вручную Не изменяется после создания. Изменение требует пересоздания. |
| `operation_timeout` | `string` | Optional | No | — |
| `password` | `string` | Computed | Yes | Пароль пользователя, который генерирует платформа. Провайдер читает его из Vault родительского инстанса сразу после создания пользователя. |
| `postgres_id` | `string` | Required | No | — |
| `role` | `string` | Required | No | Выбрать роль из списка Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
| `username` | `string` | Required | No | Не изменяется после создания. Изменение требует пересоздания. |
