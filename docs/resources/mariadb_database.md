---
page_title: "nubes_mariadb_database"
description: "nubes_mariadb_database TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_mariadb_database` Resource

Service: Mariadb.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/mariadb/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `db_name` | `string` | Optional, Computed | Пример: `db01` Не изменяется после создания. Изменение требует пересоздания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mariadb_id` | `string` | Required | — |
| `operation_timeout` | `string` | Optional | — |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
