# `nubes_flask` — Web-сервер с фреймворком Flask

Разворачивает веб-приложение на Flask из Git-репозитория: платформа собирает образ
и публикует его по указанному домену. Отдаёт адрес приложения.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/flask.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_flask" "app" {
  resource_name = "my-flask"

  startup_configuration = {
    resource_realm = "k8s-4-ext-nubes-ru"
  }
  cluster_configuration = {
    cpu      = 200
    memory   = 256
    replicas = 1 # 1 | 3 | 5 | 7
  }
  access_configuration = {
    domain = "my-flask"
  }
  app_configuration = {
    git_path  = "https://github.com/example/my-flask-app.git"
    version   = "3.12"
    health_path = "/health"
  }

  json_env = {
    FLASK_ENV = "production"
  }
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса |
| `json_env` | map | нет | `{}` | Переменные окружения приложения |

### `startup_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `resource_realm` | string | **да** | — | Платформа: `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` |

### `cluster_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `cpu` | integer | нет | `200` | Квота CPU пода (milicores) |
| `memory` | integer | нет | `256` | Квота памяти пода (MB) |
| `replicas` | integer | нет | `1` | Число реплик: `1`, `3`, `5`, `7` |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `domain` | string | нет | `flask-<случайное>` | FQDN или имя приложения |

### `app_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `git_path` | string | нет | `https://github.com/Foxyhhd/Baldurs-Gate-test.git` | Git-репозиторий приложения |
| `health_path` | string | нет | — | Путь health-check |
| `version` | string | нет | `3.12` | Версия Python |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить вместо удаления |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.connectionUrl` | string | Адрес приложения |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_out_flat["monitoring.logs"]` | string | Ссылка на логи (Grafana) |
| `state_out_flat["monitoring.resourceMetrics"]` | string | Ссылка на метрики (Grafana) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy` | `true` | Инстанс переводится в `suspend` (не удаляется) |
| `suspend_on_destroy` | `false` | Выполняется операция `delete` — инстанс удаляется в облаке |
| `keep_on_destroy` | `false` | Если `true` — облако не трогается, ресурс только убирается из состояния (`state_only`); приоритетнее `suspend_on_destroy` |

Возврат из `suspend`: повторный `terraform apply` с `adopt_existing_on_create = true`
— провайдер найдёт инстанс по имени и выполнит `resume`.

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
- `redeploy` запускается из Terraform сменой поля `git_revision`; `reconcile` — ресурсом `nubes_service_operation` (`operation = "reconcile"`).

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Развёртывание приложения |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `app_configuration`, `json_env` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |
| `redeploy` | ✅ `apply` (при смене `git_revision`) | Пересборка из git (поле-триггер) |


See the [Terraform Registry attribute reference](../reference/resources/flask.md).
