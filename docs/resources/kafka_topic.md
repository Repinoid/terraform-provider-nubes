---
page_title: "nubes_kafka_topic"
description: "nubes_kafka_topic TEST resource schema for Nubes Cloud."
subcategory: "Очереди"
---

# `nubes_kafka_topic` Resource

Service: ApacheKafka.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/kafka/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `kafka_id` | `string` | Required | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: подресурс не удаляется и не меняется в облаке, только убирается из состояния. Нужен там, где родительский инстанс при destroy не удаляется (например suspend_on_destroy = true) и объект должен остаться внутри него. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `name_topic` | `string` | Required | No | Пример: `orders-events` Изменение требует пересоздания. |
| `operation_timeout` | `string` | Optional | No | — |
| `partitions` | `number` | Required | No | Число партиций можно увеличить, но уменьшить их нельзя. |
| `replicas` | `number` | Required | No | Количество реплик не может привышать количество реплик брокера |
| `skip_missing_on_delete` | `bool` | Optional, Computed | No | — |
