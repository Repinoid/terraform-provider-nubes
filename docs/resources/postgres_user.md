---
page_title: "nubes_postgres_user"
description: "nubes_postgres_user TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_postgres_user` Resource

Service: PostgreSQL.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/postgres/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mtls_access` | `bool` | Optional, Computed | Подключение пользователя без ввода паролей<br/>Политика применяется **ко всем базам** внутри Postgresql<br/><br/>RBAC на базы настраивать вручную Не изменяется после создания. Изменение требует пересоздания. |
| `operation_timeout` | `string` | Optional | — |
| `password` | `string` | Computed | **Sensitive.** Пароль пользователя, который генерирует платформа. Провайдер читает его из Vault родительского инстанса сразу после создания пользователя. |
| `postgres_id` | `string` | Required | — |
| `role` | `string` | Required | Выбрать роль из списка Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
| `username` | `string` | Required | Не изменяется после создания. Изменение требует пересоздания. |
