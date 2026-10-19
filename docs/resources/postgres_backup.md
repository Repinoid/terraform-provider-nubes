---
page_title: "nubes_postgres_backup"
description: "nubes_postgres_backup TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_postgres_backup` Resource

Service: PostgreSQL.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/postgres/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `postgres_id` | `string` | Required | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
