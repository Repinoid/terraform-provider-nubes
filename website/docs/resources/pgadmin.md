# `nubes_pgadmin` — pgAdmin

Веб-инструмент с графическим интерфейсом для управления базами данных PostgreSQL.
Отдаёт адрес веб-консоли.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/pgadmin.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_pgadmin" "admin" {
  resource_name     = "my-pgadmin"
  operation_timeout = "15m"

  resource_realm = "k8s-4-ext-nubes-ru"

  # login должен быть email
  login    = "admin@example.com"
  password = var.pgadmin_password # sensitive

  resource_cpu    = 500 # 500…4000
  resource_memory = 512 # 512…8192
  resource_disk   = 1   # 1…999
  domain          = "my-pgadmin"
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `resource_realm` | string | **да** | — | Платформа: `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` | список |
| `login` | string | **да** (sensitive) | — | Логин для входа в pgAdmin | **email** (`login@email.net`), длина 3…255 |
| `password` | string | **да** (sensitive) | — | Пароль для входа в pgAdmin | — |
| `domain` | string | нет | `pgadmin-1` | FQDN или имя; при FQDN проверяется A-запись DNS | уникален |
| `resource_cpu` | integer | нет | `500` | Квота CPU пода (milicores) | 500…4000 |
| `resource_memory` | integer | нет | `512` | Квота памяти пода (MB) | 512…8192 |
| `resource_disk` | integer | нет | `1` | Квота диска (GB) | 1…999 |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `15m`) |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить вместо удаления |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.connectionUrl` | string | Адрес веб-консоли: `https://<domain>.pgadmink8s.dev.nubes.ru` |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_out_flat["monitoring.logs"]` | string | Ссылка на логи (Grafana) |
| `state_out_flat["monitoring.resourceMetrics"]` | string | Ссылка на метрики (Grafana) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API (`password` замаскирован `*******`) |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

```hcl
output "pgadmin_url" { value = nubes_pgadmin.admin.state_out_flat["connectionUrl"] }
```

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy` | `true` | Инстанс переводится в `suspend` (не удаляется) |
| `suspend_on_destroy` | `false` | Выполняется операция `delete` — инстанс удаляется в облаке |
| `keep_on_destroy` | `false` | Если `true` — ресурс вообще не трогается в облаке (state_only) |

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
- `redeploy` доступен из Terraform через поле-триггер `git_revision`.

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ✅ `apply` | Меняются `resource_cpu`, `resource_memory`, `resource_disk` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `redeploy` | ✅ `apply` (при смене `git_revision`) | Пересборка из git (поле-триггер) |


See the [Terraform Registry attribute reference](../reference/resources/pgadmin.md).
