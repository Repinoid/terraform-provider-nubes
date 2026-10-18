---
page_title: "nubes_kafka_user"
description: "nubes_kafka_user TEST resource schema for Nubes Cloud."
subcategory: "Очереди"
---

# `nubes_kafka_user` Resource

Service: ApacheKafka.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/kafka/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_hosts` | `string` | Optional, Computed | No | Пример: `192.168.1.0/24` или `*` |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `group` | `string` | Optional, Computed | No | Пример: `order-consumers` или `*` |
| `id` | `string` | Computed | No | — |
| `kafka_id` | `string` | Required | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `name_topic` | `string` | Required | No | Пример: `orders-events`. Использовать маску нельзя(`orders*`), только полное название или `*` |
| `operation_timeout` | `string` | Optional | No | — |
| `operations` | `string` | Required | No | Пример: `Create, Describe, Read, Write` |
| `password` | `string` | Computed | Yes | Пароль пользователя, который генерирует платформа. Провайдер читает его из Vault родительского инстанса сразу после создания пользователя. |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
| `username` | `string` | Required | No | Пример: `username-admin` Изменение требует пересоздания. |
