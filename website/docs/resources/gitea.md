# `nubes_gitea` — Gitea

Лёгкий self-hosted хостинг Git-репозиториев с веб-интерфейсом (репозитории,
организации, команды, CI/CD-интеграции). Работает поверх существующей базы
PostgreSQL (`nubes_postgres`), отдаёт адрес веб-интерфейса.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/gitea.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
> **Gitea создаётся поверх существующего PostgreSQL**
> Для развёртывания необходима уже развёрнутая база PostgreSQL (`nubes_postgres`).
> `psql_uid` в схеме помечен как необязательный, однако на практике обязателен
> (см. ниже в Argument Reference). Пользователь `techGitea` и БД `gitea` создаются
> платформой автоматически.
>
---

## Пример использования

```hcl
resource "nubes_gitea" "git" {
  resource_name = "my-gitea"

  # ⚠️ обязателен: UUID существующего PostgreSQL (сервис 90)
  psql_uid = nubes_postgres.db.id

  resource_cpu       = 500  # 300…4000
  resource_memory    = 512  # 256…8192
  resource_disk      = 1    # 1…999
  resource_instances = 1    # 1…7

  domain = "my-gitea"

  ini_server_params = {
    "server.ROOT_URL" = "https://my-gitea.example.com/"
  }
}
```

> ⚠️ `psql_uid` в схеме помечен как необязательный, однако без него создание не
> запускается — указывайте UUID существующей базы PostgreSQL.
>
> ⚠️ `resource_realm` в схеме отсутствует — платформа подставляет кластер сама.

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `psql_uid` | string (uuid) | нет *(фактически да)* | — | UUID услуги k8sPgsql (сервис 90). Пользователь `techGitea` и БД `gitea` создаются автоматически | ссылка на `nubes_postgres` |
| `domain` | string | нет | `gitea-1` | FQDN или имя; при FQDN проверяется A-запись DNS | DNS-имя/FQDN |
| `resource_cpu` | integer | нет | `500` | Квота CPU пода (milicores) | 300…4000 |
| `resource_memory` | integer | нет | `512` | Квота памяти пода (MB) | 256…8192 |
| `resource_disk` | integer | нет | `1` | Квота диска (GB) | 1…999 |
| `resource_instances` | integer | нет | `1` | Количество узлов | 1…7 |
| `ingress_extra_annotations` | map | нет | — | Доп. аннотации ingress | — |
| `ini_server_params` | map | нет | — | Параметры `[server]` в `app.ini` | — |

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
| `state_out` / `state_out_flat` | map | Ответ API (адрес веб-интерфейса, `instanceUid`, `monitoring`) |
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

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ✅ `apply` | Меняются `resource_cpu`, `resource_memory`, `resource_disk`, `resource_instances`, `ingress_extra_annotations`, `ini_server_params` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |


See the [Terraform Registry attribute reference](../reference/resources/gitea.md).
