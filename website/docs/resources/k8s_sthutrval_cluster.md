# `nubes_k8s_sthutrval_cluster` — Kubernetes кластер Штурвал

Разворачивает полноценный Kubernetes-кластер «Штурвал» на VMware Cloud Director
(control plane + worker-группы) и отдаёт доступ к API и ingress.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/k8s_sthutrval_cluster.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
> ⚠️ Перед созданием убедитесь, что выполнены условия развёртывания из инструкции
> сервиса: в частности, **1 кластер = 1 vDC**, vDC и NSX-T должны быть готовы.

---

## Пример использования

```hcl
resource "nubes_k8s_sthutrval_cluster" "k8s" {
  resource_name     = "my-k8s"
  operation_timeout = "60m"

  startup_configuration = {
    vdc_uid      = nubes_vc_vdc.this.id
    nsxt_uid     = "nsxt-uid"
    cluster_name = "my-k8s"
  }
  control_plane_configuration = {
    sizing_policy = "Intel Broadwell 2.4" # значение из ЛК
    sizing_disk   = 50
    count         = 1 # 1 | 3 | 5
  }
  worker_configuration = {
    group_name    = "workers"
    sizing_policy = "Intel Broadwell 2.4" # значение из ЛК
    sizing_disk   = 50
    count         = 2
    autoscale     = false
    label_deck    = true
  }
  access_configuration = {
    need_external_address_api     = true
    need_external_address_ingress = true
  }
  cluster_configuration = {
    app_version = "2.14.0" # 2.14.0 | 2.13.1 | 2.12.1 | 2.11.0
  }
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
| `vdc_uid` | string | **да** | — | UUID существующего vDC |
| `nsxt_uid` | string | **да** | — | UUID NSX-T |
| `cluster_name` | string | нет | `example-<случайное>` | Имя кластера |
| `ex_logging` | bool | нет | `true` | Включить логирование |
| `ex_monitoring` | bool | нет | `true` | Включить мониторинг |
| `ex_local_csi` | bool | нет | `true` | Локальный CSI |
| `ex_vip` | bool | нет | `true` | Виртуальный IP |
| `ex_update` | bool | нет | `true` | Автообновление |
| `ex_ingress` | bool | нет | `true` | Ingress |
| `ex_named_csi` | bool | нет | `true` | Именованный CSI |

### `control_plane_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `sizing_policy` | string | нет | — | Sizing policy из ЛК |
| `sizing_disk` | integer | нет | `50` | Размер диска (GB) |
| `count` | integer | нет | `1` | Число узлов: `1`, `3`, `5` |

### `worker_configuration` *(список групп воркеров)*

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `group_name` | string | нет | `workers-<случайное>` | Имя группы |
| `sizing_policy` | string | **да** | — | Sizing policy из ЛК |
| `sizing_disk` | integer | нет | `50` | Размер диска (GB) |
| `count` | integer | нет | `2` | Число узлов |
| `autoscale` | bool | нет | `false` | Включить автоскейлинг |
| `autoscale_min` | integer | нет | `2` | Минимум узлов |
| `autoscale_max` | integer | нет | `3` | Максимум узлов |
| `label_deck` | bool | нет | `true` | Метка для Deck |

### `access_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `need_external_address_api` | bool | нет | `true` | Внешний адрес для API |
| `access_ip_list_api` | string (json) | нет | `[]` | Список разрешённых IP для API |
| `need_external_address_ingress` | bool | нет | `true` | Внешний адрес для ingress |
| `access_ip_list_ingress` | string (json) | нет | `[]` | Список разрешённых IP для ingress |

### `cluster_configuration`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `app_version` | string | нет | `2.14.0` | Версия: `2.14.0`, `2.13.1`, `2.12.1`, `2.11.0` |

Общие атрибуты (есть у всех ресурсов):

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `operation_timeout` | string | нет | — | Максимальное время ожидания операции (напр. `60m`) |
| `log_level` | string | нет | `none` | Логирование стадий: `none`, `info`, `debug` |
| `suspend_on_destroy` | bool | нет | `true` | При destroy — усыпить вместо удаления |
| `keep_on_destroy` | bool | нет | `false` | При destroy — ничего не делать в облаке (state_only) |
| `adopt_existing_on_create` | bool | нет | `false` | Авто-подхват существующего инстанса при apply |

---

## Attribute Reference

| Имя | Тип | Описание |
|---|---|---|
| `id` | string | UUID инстанса |
| `state_out` / `state_out_flat` | map | Ответ API (endpoint API/ingress, `instanceUid`, `monitoring`) |
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
| `keep_on_destroy` | `false` | Если `true` — ресурс вообще не трогается в облаке (state_only) |

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
- `reconcile` запускается из Terraform ресурсом `nubes_service_operation` (`operation = "reconcile"`, уникальный `run_id`).

---

## Подресурсы

### `nubes_k8s_sthutrval_cluster_user`

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `k8s_sthutrval_cluster_id` | string (uuid) | **да** | — | UUID кластера |
| `admins` | map (список) | нет | `[{}]` | Администраторы кластера |

`admins` (элемент списка):

| Имя | Тип | Обязательный | Описание |
|---|---|---|---|
| `email` | string | **да** | Email администратора |

> **Общие атрибуты подресурсов**
> `adopt_existing_on_create`, `keep_on_destroy`, `skip_missing_on_delete`,
> `operation_timeout`, `log_level`.
>
---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание кластера |
| `modify` | ✅ `apply` | Меняются `cluster_configuration`, `control_plane_configuration`, `worker_configuration`, `access_configuration` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Усыпление |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Пробуждение |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `create_user` / `delete_user` | ✅ через `nubes_k8s_sthutrval_cluster_user` | Управление пользователями |
| `reconcile` | ✅ через `nubes_service_operation` | Синхронизация состояния (универсальная операция) |


See the [Terraform Registry attribute reference](../reference/resources/k8s_sthutrval_cluster.md).
