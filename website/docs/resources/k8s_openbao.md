# `nubes_k8s_openbao` — Vault

Централизованное защищённое хранилище секретов (ключей, паролей, конфигураций) на базе
OpenBao с разграничением доступа по ролям. Отдаёт адрес хранилища и ключи распечатки
(unseal) с root-токеном.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/k8s_openbao.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

Минимальный (обязательные параметры):

```hcl
resource "nubes_k8s_openbao" "vault" {
  resource_name  = "my-openbao"
  resource_realm = "k8s-4-ext-nubes-ru"
}
```

Типовой (экономные ресурсы):

```hcl
resource "nubes_k8s_openbao" "vault" {
  resource_name     = "my-openbao"
  operation_timeout = "20m"

  resource_realm     = "k8s-4-ext-nubes-ru"
  resource_cpu       = 250 # 250…8000 (milicores)
  resource_memory    = 256 # 256…8192 (MB)
  resource_disk      = 1   # 1…999 (GB)
  resource_instances = 1   # 1…7 (по умолчанию 3)
  app_version        = "2.5.2"
  domain             = "my-openbao"
}
```

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `resource_realm` | string | **да** | — | K8s-кластер: `k8s-4-ext-nubes-ru`, `k8s-4-ext-nubes-ru` | список |
| `resource_cpu` | integer | нет | `500` | Квота CPU пода (milicores) | 250…8000 |
| `resource_memory` | integer | нет | `512` | Квота памяти пода (MB) | 256…8192 |
| `resource_disk` | integer | нет | `10` | Квота диска (GB) | 1…999 |
| `resource_instances` | integer | нет | `3` | Количество узлов | 1…7 |
| `app_version` | string | нет | `2.5.2` | Версия приложения | `2.5.2` |
| `domain` | string | нет | `bao-<случайное>` | FQDN или имя | DNS-имя/FQDN |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `20m`) |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить вместо удаления |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out.connectionUrl` | string | Адрес хранилища: `https://<domain>.openbaok8s-<uid>.dev.nubes.ru` |
| `state_out.instanceUid` | string | UUID инстанса |
| `state_out_flat["monitoring.logs"]` | string | Ссылка на логи (Grafana) |
| `state_out_flat["monitoring.resourceMetrics"]` | string | Ссылка на метрики (Grafana) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Ключи распечатки: `unsealKey_0`…`unsealKey_4`, `rootKey` |
| `vault_fields` | list | `["unsealKey_0", …, "unsealKey_4", "rootKey"]` |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |

```hcl
output "openbao_url"  { value = nubes_k8s_openbao.vault.state_out_flat["connectionUrl"] }
output "openbao_root" { value = nubes_k8s_openbao.vault.vault_secrets["rootKey"] }
```

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

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
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Подресурсы

### `nubes_k8s_openbao_kv`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `k8s_openbao_id` | string (uuid) | **да** | — | UUID инстанса OpenBao |
| `kv_name` | string | нет | `kvname-<случайное>` | Имя KV-хранилища |

### `nubes_k8s_openbao_policy`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `k8s_openbao_id` | string (uuid) | **да** | — | UUID инстанса OpenBao |
| `policy_name` | string | нет | `policy-<случайное>` | Имя политики |
| `policy_settings` | map (список) | **да** | — | Правила политики |

`policy_settings` (элемент списка):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `path` | string | **да** | — | Путь, к которому применяется правило |
| `capabilities` | string | нет | `["read", "list"]` | Список возможностей (JSON-массив) |

### `nubes_k8s_openbao_secrets`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `k8s_openbao_id` | string (uuid) | **да** | — | UUID инстанса OpenBao |
| `kv_name` | string | **да** | — | Имя KV-хранилища |
| `secret_path` | string | нет | `path/to/mysecret` | Путь секрета |
| `secret_data` | map | **да** | — | Содержимое секрета (ключ → значение) |

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание инстанса |
| `modify` | ✅ `apply` | Меняются `resource_instances`, `resource_cpu`, `resource_memory`, `resource_disk`, `app_version` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `create_kv` / `delete_kv` | ✅ через `nubes_k8s_openbao_kv` | Управление KV |
| `create_policy` / `update_policy` / `delete_policy` | ✅ через `nubes_k8s_openbao_policy` | Управление политиками |
| `create_secrets` / `update_secrets` / `delete_secrets` / `recovery_secrets` | нет | Недоступны из Terraform |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/k8s_openbao.md).
