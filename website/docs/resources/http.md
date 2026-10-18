# `nubes_http` — Простой HTTP контейнер

Разворачивает произвольный Docker-образ из реестра и публикует его по указанному
домену. Отдаёт адрес опубликованного контейнера.

> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
resource "nubes_http" "app" {
  resource_name  = "my-http"
  resource_realm = "k8s-4-ext-nubes-ru"

  domain        = "my-http"
  registry_path = "nexus-sa.tst.nubes.ru/docker-nubes/jolt"

  resource_cpu       = 100 # 100…4000 (milicores)
  resource_memory    = 200 # 128…8192 (MB)
  resource_instances = 1   # 1…7

  json_env = {
    PORT = "8080"
  }

  ingress_annotations = {
    "nginx.ingress.kubernetes.io/proxy-body-size" = "10m"
  }
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `resource_realm` | string | **да** | — | Платформа: `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` | список |
| `domain` | string | нет | `http-1` | FQDN или имя; при FQDN проверяется A-запись DNS | DNS-имя/FQDN |
| `registry_path` | string | нет | `nexus-sa.tst.nubes.ru/docker-nubes/jolt` | Путь до образа, который ставится как deployment | — |
| `resource_cpu` | integer | нет | `100` | Квота CPU пода (milicores) | 100…4000 |
| `resource_memory` | integer | нет | `200` | Квота памяти пода (MB) | 128…8192 |
| `resource_instances` | integer | нет | `1` | Количество узлов | 1…7 |
| `json_env` | map | нет | — | Переменные окружения контейнера | — |
| `ingress_annotations` | map | нет | `{}` | Аннотации ingress | — |

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
| `state_out.connectionUrl` | string | Адрес опубликованного контейнера |
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
| `create` | ✅ `apply` | Развёртывание контейнера |
| `modify` | ✅ `apply` | Меняются `resource_cpu`, `resource_memory`, `resource_instances`, `registry_path`, `json_env`, `ingress_annotations` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |
| `redeploy` | ✅ `apply` (при смене `git_revision`) | Пересборка из git (поле-триггер) |


See the [Terraform Registry attribute reference](../reference/resources/http.md).
