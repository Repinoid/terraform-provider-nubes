# `nubes_akhq` — Akhq

Веб-интерфейс для Apache Kafka: просмотр топиков, consumer-групп, сообщений, ACL,
состояния брокеров и метаданных. Работает поверх существующего кластера Kafka
(`nubes_kafka`), отдаёт ссылку на веб-UI.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/akhq.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> У этого ресурса `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`),
> а не удаляет. Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
> **AKHQ создаётся поверх существующего Kafka**
> Для развёртывания необходим уже развёрнутый кластер Kafka (`nubes_kafka`).
> `kafka_uid` в схеме помечен как необязательный, однако на практике обязателен
> (см. ниже в Argument Reference).
>
---

## Пример использования

```hcl
resource "nubes_kafka" "cluster" {
  resource_name = "my-kafka"

  startup_configuration = { resource_realm = "k8s-4-ext-nubes-ru" }
  cluster_configuration = {
    cpu      = 2000
    memory   = 2048
    disk     = 1
    replicas = 1
  }
}

resource "nubes_akhq" "ui" {
  resource_name  = "my-akhq"
  resource_realm = "k8s-4-ext-nubes-ru" # та же платформа, что у Kafka
  domain         = "my-akhq"

  # ⚠️ обязателен на практике: UUID существующего кластера Kafka (сервис 116)
  kafka_uid = nubes_kafka.cluster.id

  resource_cpu       = 500  # min 500
  resource_memory    = 512  # min 512
  resource_instances = 1    # min 1
}
```

> ⚠️ `kafka_uid` в схеме помечен как необязательный, однако на практике без него
> инстанс не создаётся — указывайте UUID существующего кластера Kafka.

---

## Argument Reference

Параметры создания ресурса:

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `resource_realm` | string | **да** | — | Платформа/K8s-кластер для развёртывания | `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` |
| `domain` | string | **да** | — | FQDN или имя; при FQDN проверяется A-запись DNS | DNS-имя/FQDN |
| `kafka_uid` | string (uuid) | нет *(фактически да)* | — | UUID экземпляра Kafka (сервис 116) | ссылка на `nubes_kafka` |
| `resource_cpu` | integer | нет | `500` | Квота CPU пода (milicores) | ≥ 500 |
| `resource_memory` | integer | нет | `512` | Квота памяти пода (MB) | ≥ 512 |
| `resource_instances` | integer | нет | `1` | Количество узлов | ≥ 1 |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `15m`) |
| `log_level` | string | нет | `none` | Логирование стадий операции: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить (`true`) вместо удаления (`false`) |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only); приоритетнее suspend |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply (=import) |

---

## Attribute Reference

Вычисляемые атрибуты (только чтение):

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out` | map | Ответ API: `connectionUrl`, `instanceUid`, `monitoring` (json) |
| `state_out_flat` | map | То же плоско: `connectionUrl`, `monitoring.logs`, `monitoring.resourceMetrics` |
| `state_params` | map | Параметры, отправленные в API |
| `state_params_flat` | map | То же плоско |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` | string | Адрес Vault для ресурса |
| `vault_user_path` | string | Путь пользователя в Vault |
| `vault_fields` | list | Список доступных ключей секретов |

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

⚠️ `terraform destroy` **не удаляет** инстанс, а **усыпляет** его
(`suspend_on_destroy = true` по умолчанию). Данные сохраняются.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy` | `true` | Инстанс переводится в `suspend` (не удаляется) |
| `suspend_on_destroy` | `false` | Выполняется операция `delete` — инстанс удаляется в облаке |
| `keep_on_destroy` | `false` | Если `true` — облако не трогается, ресурс только убирается из состояния (`state_only`); приоритетнее `suspend_on_destroy` |

Возврат из `suspend`: повторный `terraform apply` с `adopt_existing_on_create = true`
— провайдер найдёт инстанс по имени и выполнит `resume` (тот же `id`).

> **Зачем нужно усыпление**
> `suspend` — временная остановка: ресурсы высвобождаются, а инстанс можно
> вернуть в работу в любой момент повторным `apply`.
>
> **Удаление усыплённого инстанса**
> Удалить усыплённый инстанс можно, когда он пробыл в `suspend` **14 дней**.
>
> ⚠️ Из Terraform это невозможно — операция `delete` для усыплённого инстанса
> недоступна. Удаляйте его вручную в Личном кабинете (ЛК).
>
- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.
- `reconcile` **недоступен** — операция отсутствует в API сервиса.

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ✅ `apply` (изменение параметров) | Меняется `kafka_uid` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |


See the [Terraform Registry attribute reference](../reference/resources/akhq.md).
