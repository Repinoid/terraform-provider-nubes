---
page_title: "nubes_mariadb_user"
description: "nubes_mariadb_user TEST resource schema for Nubes Cloud."
subcategory: "Базы данных"
---

# `nubes_mariadb_user` Resource

Service: Mariadb.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/mariadb/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `db_name` | `string` | Optional, Computed | Необходимо указать массив существующих баз Не изменяется после создания. Изменение требует пересоздания. |
| `id` | `string` | Computed | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `mariadb_id` | `string` | Required | — |
| `operation_timeout` | `string` | Optional | — |
| `password` | `string` | Optional, Computed | **Sensitive.** Секрет будет записан в хранилище секретов и не будет отображаться публично Не изменяется после создания. Изменение требует пересоздания. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
| `username` | `string` | Optional, Computed | Пример: `myusername` Не изменяется после создания. Изменение требует пересоздания. |
