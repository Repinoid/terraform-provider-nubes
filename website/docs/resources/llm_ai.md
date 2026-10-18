# `nubes_llm_ai` — LLM AI

Доступ к большой языковой модели через OpenAI-совместимый API. Сервис создаёт
экземпляр LLM-платформы и отдаёт адрес для запросов к чат-модели.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/llm_ai.md).
>
> **Особенность сервиса: `destroy` сразу удаляет**
> У этого сервиса нет операции `suspend`. `terraform destroy` **удаляет** инстанс
> в облаке. Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_llm_ai" "llm" {
  resource_name = "my-llm"
}
```

Типовые подресурсы (секреты и команда):

```hcl
resource "nubes_llm_ai_secrets" "key" {
  llm_ai_id = nubes_llm_ai.llm.id
  key_name  = "app-key"
}

resource "nubes_llm_ai_team" "team" {
  llm_ai_id = nubes_llm_ai.llm.id
  team_name = "app-team"
  models    = jsonencode(["gpt-4o-mini"])
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса |

### `startup_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_realm` | string | нет | `llmAi.aillm.ru` | Платформа LLM (фиксирована) |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.connectionUrl` | string | Адрес API чат-модели (`https://api.aillm.ru/v1/chat/completions`) |
| `state_out_flat["connectionUrl"]` | string | То же в плоском виде |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API (`resourceRealm`) |
| `vault_secrets` | map (sensitive) | Секреты из Vault (пусто) |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов (пусто) |

> **Точка подключения**
> `state_out.connectionUrl` = `https://api.aillm.ru/v1/chat/completions`
> (OpenAI-совместимый endpoint). Секретов у инстанса нет — ключи создаются
> подресурсом `nubes_llm_ai_secrets`.
>
---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — ресурс вообще не трогается в облаке (state_only) |
| `suspend_on_destroy` | — | ❌ у сервиса нет операции `suspend`: `destroy` **удаляет** инстанс |

- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Подресурсы

Ключи и команды создаются подресурсами `nubes_llm_ai_secrets` и `nubes_llm_ai_team`;
их создание требует предварительной операции `reconcile` инстанса (запускается
ресурсом `nubes_service_operation`, `operation = "reconcile"`).

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ⛔ нет | Операция модификации инстанса недоступна |
| `delete` | ✅ `destroy` | Удаление инстанса в облаке (усыпления нет) |
| `create_secrets` / `update_secrets` / `delete_secrets` | ✅ через `nubes_llm_ai_secrets` | Управление ключами |
| `create_team` / `modify_team` / `delete_team` | ✅ через `nubes_llm_ai_team` | Управление командами |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/llm_ai.md).
