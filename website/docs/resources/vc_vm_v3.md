# `nubes_vc_vm_v3` — Виртуальная машина

Разворачивает виртуальную машину в vApp (виртуальном каталоге ВМ) с заданными CPU,
памятью, диском, образом и правилами доступа. Поддерживает `cloud-init` и SSH-ключ.

> **Инструкция по сервису**
> Отдельная страница: [руководство по сервису](../guides/services/vc_vm_v3.md).
>
> **Специфика Nubes: `destroy` ≠ удаление**
> `terraform destroy` по умолчанию **усыпляет** инстанс (`suspend`), а не удаляет.
> Подробности — в разделе [«Жизненный цикл»](#zhiznennyy-cikl) и в общем
> руководстве [«Как работает провайдер»](../guides/provider-behavior.md).
>
---

## Пример использования

```hcl
variable "vapp_uid" {
  type        = string
  description = "UUID или имя существующей vApp из ЛК"
}

variable "ip_space_name" {
  type        = string
  description = "Имя IP-space из ЛК; укажите no-needed, если внешний IP не требуется"
}
```

```hcl
resource "nubes_vc_vm_v3" "vm" {
  resource_name = "my-vm"

  vapp_uid   = var.vapp_uid
  vm_name    = "web01"
  image_vm   = "Ubuntu_22-20G" # значение из Личного кабинета

  vm_cpu  = 1  # 1…64
  vm_ram  = 1  # 1…256 (GB)
  vm_disk = 20 # GB (доп. диск; основной зависит от образа)

  ip_space_name    = var.ip_space_name
  access_ip_list   = jsonencode(["192.0.2.10/32"]) # замените своим CIDR
  user_login       = "myuser"
  user_public_key  = file(pathexpand("~/.ssh/id_ed25519.pub"))

  cloud_init = <<-EOT
    #cloud-config
    packages:
      - nginx
  EOT

  # Порты 22/tcp и 10050/tcp платформа добавляет автоматически.
  access_port_list = jsonencode([
    { port = "80", type = "tcp" },
  ])

  same_snat = false

  # Изменение этих маркеров в существующем ресурсе запускает операции.
  git_revision      = "redeploy-1"
  reconcile_trigger = "reconcile-1"
}
```

`vapp_uid` должен указывать на существующую vApp: задайте её UUID или имя из ЛК,
либо используйте ссылку на ресурс `nubes_vapp`, объявленный в этой конфигурации.
`access_ip_list` ограничивает адреса, которым разрешён доступ к ВМ; не оставляйте
широкое правило `0.0.0.0/0`, если доступ всем адресам не требуется.

---

## Argument Reference

| Имя | Тип | Обязательный | По умолчанию | Описание | Ограничения |
|---|---|---|---|---|---|
| `resource_name` | string | **да** | — | Имя инстанса | — |
| `ip_space_name` | string | **да** | — | Имя IP-space из ЛК; инструкция сервиса допускает `no-needed`, если внешний IP не требуется | — |
| `image_vm` | string | **да** | — | Имя образа ВМ из ЛК (например, `Ubuntu_22-20G`); после создания не изменяется | — |
| `user_public_key` | string | **да** | — | Публичный SSH-ключ пользователя; после создания не изменяется | — |
| `access_port_list` | string (json) | **да** | — | Разрешённые порты (JSON-массив `{port,type}`); `22/tcp` и `10050/tcp` платформа добавляет сама | — |
| `vapp_uid` | string | **да при создании** | — | UUID или имя существующей vApp (сервис 26); после создания не изменяется | ссылка на `nubes_vapp` или значение из ЛК |
| `vm_name` | string | нет | `web01` | Имя ВМ; после создания не изменяется | — |
| `vm_cpu` | integer | нет | `1` | Количество CPU (шт) | 1…64 |
| `vm_ram` | integer | нет | `1` | Объём RAM (GB) | 1…256 |
| `vm_disk` | integer | нет | — | Дополнительный диск (GB); основной зависит от образа | — |
| `access_ip_list` | string (json) | нет | `["0.0.0.0/0"]` | Разрешённые адреса (JSON-массив) | CIDR |
| `git_revision` | string | нет | — | Триггер `redeploy`: изменение значения запускает операцию. Само значение не передаётся как ревизия Git в API | — |
| `reconcile_trigger` | string | нет | — | Триггер `reconcile`: изменение или очистка значения запускает операцию при `Update`; при создании операция не вызывается | — |
| `cloud_init` | string | нет | — | Cloud-init конфигурация | — |
| `user_login` | string | нет | `myuser` | Логин пользователя; после создания не изменяется | — |
| `same_snat` | bool | нет | `false` | Использовать общий SNAT | `false`, `true` |

Хотя `vapp_uid` помечен в Terraform-схеме как `Optional`, провайдер требует его
при создании и завершит `Create` ошибкой, если значение не передано.

### `access_port_list` *(обязателен; JSON-строка)*

Задаётся как `jsonencode` массива объектов `{ port, type }`. Порты **22/tcp** и
**10050/tcp** платформа проставляет автоматически.

| Имя | Тип | Обязательный | По умолчанию | Описание |
|---|---|---|---|---|
| `port` | string | нет | `22` | Порт |
| `type` | string | нет | `all` | Протокол: `tcp`, `udp`, `all` |

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
| `state_out` / `state_out_flat` | map | Ответ API (адреса ВМ, `instanceUid`, `monitoring`) |
| `state_params` / `state_params_flat` | map | Параметры, отправленные в API |
| `vault_secrets` | map (sensitive) | Секреты из Vault |
| `vault_url` / `vault_user_path` | string | Адрес/путь Vault |
| `vault_fields` | list | Список доступных ключей секретов |

---

## Жизненный цикл <a id="zhiznennyy-cikl"></a>

| Флаг | Дефолт | Поведение при `destroy` |
|---|---|---|
| `suspend_on_destroy = true` | `true` | ВМ переводится в `suspend`, а не удаляется |
| `suspend_on_destroy = false` | `false` | Провайдер запрашивает `delete`; операция должна быть доступна для текущего состояния ВМ |
| `keep_on_destroy = false` | `false` | Обычное поведение: применяется `suspend_on_destroy` |
| `keep_on_destroy = true` | — | Облачная операция не вызывается; ресурс удаляется только из Terraform state. Этот флаг приоритетнее `suspend_on_destroy` |

Возврат из `suspend`: повторный `terraform apply` с `adopt_existing_on_create = true`
— провайдер найдёт ВМ по имени и выполнит `resume`.

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
- Изменение `git_revision` запускает `redeploy`; значение поля — только триггер, не параметр ревизии Git. Если одновременно меняются параметры ВМ, сначала выполняется `modify`, затем `redeploy`.
- Изменение или очистка `reconcile_trigger` запускает `reconcile` из `Update` этого же ресурса. Для этого не нужны дополнительные ресурсы. При одновременном изменении обоих триггеров порядок: `modify` → `redeploy` → `reconcile`. Create сам по себе триггеры не исполняет.
- Наличие операции в документации не гарантирует её доступность для конкретного инстанса: провайдер проверяет `availableOperations` во время выполнения.

---

## Operations

| Операция | Доступно из Terraform | Комментарий |
|---|---|---|
| `create` | ✅ `apply` | Создание ВМ |
| `modify` | ✅ `apply` | Меняются `vm_cpu`, `vm_ram`, `vm_disk`, `ip_space_name`, `access_ip_list`, `access_port_list`, `same_snat` |
| `suspend` | ✅ `destroy` (при `suspend_on_destroy=true`) | Остановка ВМ |
| `resume` | ✅ `apply` (при `adopt_existing_on_create=true`) | Запуск ВМ |
| `delete` | ✅ `destroy` (`suspend_on_destroy=false`) | Удаление инстанса в облаке. Только для работающего (`running`); усыплённый — вручную в ЛК |
| `reconcile` | ✅ `apply` (при изменении `reconcile_trigger`) | Операция вызывается из `Update` этого ресурса |
| `redeploy` | ✅ `apply` (при смене `git_revision`) | Операция вызывается по изменению поля-триггера; само значение `git_revision` не отправляется в API как ревизия |


See the [Terraform Registry attribute reference](../reference/resources/vc_vm_v3.md).
