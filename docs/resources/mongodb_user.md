---
page_title: "nubes_mongodb_user"
description: "nubes_mongodb_user TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_mongodb_user` Resource

Service: MongoDB.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/mongodb/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `db_name` | `string` | Required | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mongodb_id` | `string` | Required | — |
| `operation_timeout` | `string` | Optional | — |
| `password` | `string` | Computed | **Sensitive.** Пароль пользователя, который генерирует платформа. Провайдер читает его из Vault родительского инстанса сразу после создания пользователя. |
| `role` | `string` | Required | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
| `username` | `string` | Required | Изменение требует пересоздания. |
