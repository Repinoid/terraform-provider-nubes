# `nubes_superset` — Apache Superset

BI-инструмент Apache Superset как сервис (k8s-operator). Авторизация — 2FA Keycloak
(для доменных логинов `@mgmt.nubes.ru`). Отдаёт адрес веб-интерфейса.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/superset.md).
>
> **Специфика Nubes: `destroy` для этого ресурса **удаляет** инстанс**
> `terraform destroy` удаляет инстанс из облака и из Terraform-state. Усыпление
> (`suspend`) у этого сервиса не поддерживается. Подробности — в разделе
> [«Жизненный цикл»](#zhiznennyy-cikl) и в
> общем руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
> **Обязательные параметры**
> - `emails` — учётная запись, получающая админские права; должна существовать
>   в Keycloak (дефолтный `testuser1` не подходит).
> - `domain` — дефолтный `superset-1` может быть занят, задавайте уникальный.
>
>
---

## Пример использования

```hcl
resource "nubes_superset" "bi" {
  resource_name  = "my-superset"
  resource_realm = "k8s-4-ext-nubes-ru"

  domain = "my-superset"
  emails = "admin@example.com" # ⚠️ реальный email из Keycloak

  resource_cpu       = 500 # 100…4000 (milicores)
  resource_memory    = 512 # 128…8192 (MB)
  resource_disk      = 1   # 1…999 (GB)
  resource_instances = 1   # 1…7
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `resource_realm` | string | **да** | — | Платформа: `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` | список |
| `domain` | string | нет | `superset-1` | FQDN или имя; при FQDN проверяется A-запись DNS | должен быть уникален |
| `emails` | string | нет | `testuser1` | Учётная запись, получающая админские права (email, для доменных — без суффикса) | ⚠️ должна существовать в Keycloak |
| `resource_cpu` | integer | нет | `500` | Квота CPU пода (milicores) | 100…4000 |
| `resource_memory` | integer | нет | `512` | Квота памяти пода (MB) | 128…8192 |
| `resource_disk` | integer | нет | `1` | Квота диска (GB) | 1…999 |
| `resource_instances` | integer | нет | `1` | Количество узлов | 1…7 |

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
| `state_out` / `state_out_flat` | map | Ответ API (адрес веб-интерфейса, `instanceUid`, `monitoring`) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

`terraform destroy` удаляет инстанс из облака и из Terraform-state. Усыпление
(`suspend`) у этого сервиса не поддерживается — атрибута `suspend_on_destroy` нет.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `keep_on_destroy` | `false` | Если `true` — ресурс остаётся в облаке (`state_only`); удаление в облаке не выполняется |

- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `delete` | ✅ `destroy` | Удаление инстанса |

> ℹ️ Доступны только `create` и `delete` — изменение параметров не поддерживается.


See the [Terraform Registry attribute reference](../reference/resources/superset.md).
