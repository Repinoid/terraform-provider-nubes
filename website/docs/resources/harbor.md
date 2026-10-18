# `nubes_harbor` — Container Registry

Приватный Docker-реестр (Harbor) с контролем доступа и репликацией. Хранит и раздаёт
контейнерные образы; после создания отдаёт адрес для `docker login`.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/harbor.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> У этого ресурса `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`),
> а не удаляет. Подробности и флаги — в разделе [«Жизненный цикл»](#zhiznennyy-cikl)
> и в общем руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

Минимальный (обязательные параметры):

```hcl
resource "nubes_harbor" "registry" {
  resource_name  = "my-harbor"
  resource_realm = "k8s-4-ext-nubes-ru"
}
```

Типовой (S3-хранилище + дефолты проставлены явно):

```hcl
resource "nubes_harbor" "registry" {
  resource_name  = "my-harbor"
  resource_realm = "k8s-4-ext-nubes-ru"

  resource_cpu       = 1000 # min 500, max 4000 (milicores)
  resource_memory    = 1024 # min 512, max 8192 (MB)
  resource_instances = 1     # min 1, max 7
  domain             = "my-harbor"

  auth_type = "local" # local | keycloak
  pv_type   = "s3"    # s3 | local

  # ⚠️ при pv_type = "s3" (по умолчанию) ОБЯЗАТЕЛЬНО укажите UUID S3:
  s3_uid = nubes_s3.storage.id
}
```

---

## Argument Reference

Параметры создания ресурса:

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `resource_realm` | string | **да** | — | Платформа/K8s-кластер для развёртывания | `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` |
| `domain` | string | нет | `harbor-1` | FQDN или имя; при FQDN проверяется A-запись DNS | DNS-имя/FQDN |
| `resource_cpu` | integer | нет | `1000` | Квота CPU пода (milicores) | 500…4000 |
| `resource_memory` | integer | нет | `1024` | Квота памяти пода (MB) | 512…8192 |
| `resource_instances` | integer | нет | `1` | Количество узлов | 1…7 |
| `pv_type` | string | нет | `s3` | Тип хранилища образов | `s3`, `local` |
| `s3_uid` | string (uuid) | **да при `pv_type=s3`** | — | UUID экземпляра S3 (сервис 12) | ссылка на `nubes_s3` |
| `resource_disk` | integer | нет | `5` | Размер диска — нужен при `pv_type=local` | 5…4096 |
| `auth_type` | string | нет | `local` | Тип авторизации после создания | `local`, `keycloak` |

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
| `vault_secrets` | map (sensitive) | Секреты из Vault: `adminPass` |
| `vault_url` | string | Адрес Vault для ресурса |
| `vault_user_path` | string | Путь пользователя в Vault |
| `vault_fields` | list | Список доступных ключей секретов (`["adminPass"]`) |

Пример получения адреса и пароля:

```hcl
output "harbor_url"  { value = nubes_harbor.registry.state_out_flat["connectionUrl"] }
output "harbor_pass" { value = nubes_harbor.registry.vault_secrets["adminPass"] }
```

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

⚠️ **Главное отличие от классического Terraform:** `terraform destroy` **не удаляет** инстанс,
а **усыпляет** его (`suspend_on_destroy = true` по умолчанию). Данные сохраняются.

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy` | `true` | Инстанс переводится в `suspend` (не удаляется) |
| `suspend_on_destroy` | `false` | Выполняется операция `delete` — инстанс удаляется в облаке |
| `keep_on_destroy` | `false` | Если `true` — ресурс вообще не трогается в облаке (state_only), приоритетнее suspend |

Возврат из `suspend`: повторный `terraform apply` с `adopt_existing_on_create = true` —
провайдер найдёт инстанс по имени и выполнит `resume` (тот же `id`).

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
Прочее:

- Инстанс в статусе **«not created»** (после упавшего create) блокирует имя и **не удаляется
  через API** — чистится вручную в Личном кабинете.
- `reconcile` запускается ресурсом `nubes_service_operation` (операция присутствует в API сервиса).

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ✅ `apply` (изменение параметров) | Меняются `resource_cpu`, `resource_memory`, `resource_instances`, `resource_disk` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/harbor.md).
