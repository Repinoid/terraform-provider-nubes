---
page_title: "nubes_kafka_user"
description: "nubes_kafka_user TEST resource schema for Nubes Cloud."
subcategory: "Очереди"
---

# `nubes_kafka_user` Resource

Service: ApacheKafka.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/kafka/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `access_hosts` | `string` | Optional, Computed | Пример: `192.168.1.0/24` или `*` |
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `group` | `string` | Optional, Computed | Пример: `order-consumers` или `*` |
| `id` | `string` | Computed | — |
| `kafka_id` | `string` | Required | — |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `name_topic` | `string` | Required | Пример: `orders-events`. Использовать маску нельзя(`orders*`), только полное название или `*` |
| `operation_timeout` | `string` | Optional | — |
| `operations` | `string` | Required | Пример: `Create, Describe, Read, Write` |
| `password` | `string` | Computed | **Sensitive.** Пароль пользователя, который генерирует платформа. Провайдер читает его из Vault родительского инстанса сразу после создания пользователя. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | — |
| `username` | `string` | Required | Пример: `username-admin` Изменение требует пересоздания. |
