# `nubes_vc_nsxt` Resource

Service: Сетевой шлюз периметра (Edge).

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/vc_nsxt/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/vc_nsxt/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `id` | `string` | Computed | No | — |
| `ip_space_name` | `string` | Optional, Computed | No | Имя ip Space для внешнего IP |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `need_enable_avi` | `bool` | Optional, Computed | No | Включение Load Balancer |
| `operation_timeout` | `string` | Optional | No | — |
| `qos_profile` | `string` | Optional, Computed | No | Параметр пока не работает<br/>Должен возвращать поле из ресурсной платформы типа vc из .vcd.hardware.gatewayQoSProfiles<br/><br/>Как по идее должен работать.<br/>Нужно заполнить или CFS vdcUid, или vdcGroupUid<br/>Идеально конечно проверять `if (params.vdcType == &quot;vdc&quot; && params.vdcUid != &quot;&quot;)` или `if (params.vdcType == &quot;vdcGroup&quot; && params.vdcGroupUid != &quot;&quot;)`<br/><br/>Далее надо пойти по стейту vdc -&gt; org -&gt; resPlatform, взять gatewayQoSProfiles<br/>Или пойти по стейту vdcgroup -&gt; vdc -&gt; org -&gt; resPlatform, взять gatewayQoSProfiles<br/><br/>Регулярку могу написать (Виталя) |
| `reconcile_trigger` | `string` | Optional | No | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | No | — |
| `routed_net_configuration` | `object` | Required | No | — |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
| `vdc_group_uid` | `string` | Optional | No | UUID Услуги `Группа виртуальных датацентров (groupvDC)` Не изменяется после создания. |
| `vdc_type` | `string` | Optional, Computed | No | Тип родительской услуги Не изменяется после создания. |
| `vdc_uid` | `string` | Optional | No | UUID Услуги `Виртуальный датацентр (vDC)` Не изменяется после создания. |
| `virtual_services_count` | `number` | Optional, Computed | No | Кол-во виртуальных сервисов, которые **резервируются** на AVI |

### Nested attributes: `routed_net_configuration`

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `ip_addr_pool` | `string` | Optional, Computed | No | **Параметр не изменяемый**<br/>Адрессный пул, которая будет присвоена routed-сети. Описывается как (10.10.10.0/24). Маска 24 обязательна. (.1) - шлюз. (.2-.254) - Под адресацию для ВМ |
| `main_dns` | `string` | Optional, Computed | No | Подставляется автоматически в resolv.conf при создании виртуальных машин |
| `second_dns` | `string` | Optional, Computed | No | Подставляется автоматически в resolv.conf при создании виртуальных машин |
