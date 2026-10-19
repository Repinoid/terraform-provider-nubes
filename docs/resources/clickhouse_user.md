---
page_title: "nubes_clickhouse_user"
description: "nubes_clickhouse_user TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_clickhouse_user` Resource

Service: ClickHouse.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/clickhouse/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_hosts` | `string` | Optional, Computed | Пример: `::/0` Не изменяется после создания. Изменение требует пересоздания. |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `clickhouse_id` | `string` | Required | — |
| `db_name` | `string` | Optional, Computed | Не изменяется после создания. Изменение требует пересоздания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `password` | `string` | Computed | **Sensitive.** Пароль пользователя, который генерирует платформа. Провайдер читает его из Vault родительского инстанса сразу после создания пользователя. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
| `username` | `string` | Optional, Computed | Пример: `myusername` Не изменяется после создания. Изменение требует пересоздания. |
