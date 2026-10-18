# `nubes_vc_vm_v3` Resource

Service: Виртуальная машина.

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/vc_vm_v3/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/vc_vm_v3/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `access_ip_list` | `string` | Optional, Computed | No | Формат массива<br/>Каким адресам будет предоставлен доступ до вм<br/><br/>Требует выделенного внешнего IP<br/>Если ничего не передано, то выставлено [&quot;0.0.0.0/0&quot;] |
| `access_port_list` | `string` | Required | No | Порты **22/tcp**, **10050/tcp** проставляются автоматически<br/><br/>Формат следующий:<br/>- **22/tcp** - доступ только для 22 по **tcp**<br/>- **53/udp** - доступ только для 53 по **udp**<br/>- **80/all** - доступ для 80 по **tcp/udp**<br/><br/>Формат массива:<br/>- [{&quot;port&quot;:&quot;22&quot;,&quot;type&quot;:&quot;all&quot;}] |
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `cloud_init` | `string` | Optional, Computed | No | Формат файла cloud-init<br/>```<br/>#cloud-config<br/><br/>manage_etc_hosts: true<br/>package_update: true<br/>```<br/><br/>Если будут указана установка пакетов/прочего, доступ к ВМ **может быть предоставлен раньше чем выполнится cloud-init** Не изменяется после создания. |
| `git_revision` | `string` | Optional | No | Git revision (commit hash/tag). Changing this triggers redeploy instead of modify. |
| `id` | `string` | Computed | No | — |
| `image_vm` | `string` | Required | No | Не изменяется после создания. |
| `ip_space_name` | `string` | Required | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | No | — |
| `same_snat` | `bool` | Optional, Computed | No | Если выставить **false**, приложение будет опубликовано за общим SNAT<br/>Если выставить **true**, приложение будет опубликовано за внешним IP, который будет создан в процессе разворачивания услуги<br/><br/>Если выставить **true** - общий SNAT становится необязательным |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `suspend_on_destroy` | `bool` | Optional, Computed | No | — |
| `user_login` | `string` | Optional, Computed | No | Учётка SSH Не изменяется после создания. |
| `user_public_key` | `string` | Required | No | Публичная часть SSH-Ключа для пользователя<br/>```<br/>$ ssh-keygen <br/>$ cat ~/.ssh/id_ed25519.pub<br/>```<br/><br/> Не изменяется после создания. |
| `vapp_uid` | `string` | Optional | No | UI display_name for service_id=26 (mapped to UUID in core) Не изменяется после создания. |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
| `vm_cpu` | `number` | Optional, Computed | No | шт |
| `vm_disk` | `number` | Optional, Computed | No | Основной диск у vm зависит от образа (см. документацию)<br/>Если необходимо дополнительное место для vm, введите колличество Gb в эту форму |
| `vm_name` | `string` | Optional, Computed | No | Не изменяется после создания. |
| `vm_ram` | `number` | Optional, Computed | No | Gb |
