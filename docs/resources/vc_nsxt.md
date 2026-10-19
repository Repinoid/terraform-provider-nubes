---
page_title: "nubes_vc_nsxt"
description: "nubes_vc_nsxt TEST resource schema for Nubes Cloud."
subcategory: "VMware"
---

# `nubes_vc_nsxt` Resource

Service: Сетевой шлюз периметра (Edge).

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/vc_nsxt/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `id` | `string` | Computed | — |
| `ip_space_name` | `string` | Optional, Computed | Имя ip Space для внешнего IP |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `need_enable_avi` | `bool` | Optional, Computed | Включение Load Balancer |
| `operation_timeout` | `string` | Optional | — |
| `qos_profile` | `string` | Optional, Computed | Параметр пока не работает<br/>Должен возвращать поле из ресурсной платформы типа vc из .vcd.hardware.gatewayQoSProfiles<br/><br/>Как по идее должен работать.<br/>Нужно заполнить или CFS vdcUid, или vdcGroupUid<br/>Идеально конечно проверять `if (params.vdcType == &quot;vdc&quot; && params.vdcUid != &quot;&quot;)` или `if (params.vdcType == &quot;vdcGroup&quot; && params.vdcGroupUid != &quot;&quot;)`<br/><br/>Далее надо пойти по стейту vdc -&gt; org -&gt; resPlatform, взять gatewayQoSProfiles<br/>Или пойти по стейту vdcgroup -&gt; vdc -&gt; org -&gt; resPlatform, взять gatewayQoSProfiles<br/><br/>Регулярку могу написать (Виталя) |
| `reconcile_trigger` | `string` | Optional | Строка-триггер: изменение или очистка значения запускает операцию reconcile при Update. При Create операция не вызывается. |
| `resource_name` | `string` | Required | — |
| `routed_net_configuration` | `object` | Required | — |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
| `vdc_group_uid` | `string` | Optional | UUID Услуги `Группа виртуальных датацентров (groupvDC)` Не изменяется после создания. |
| `vdc_type` | `string` | Optional, Computed | Тип родительской услуги Не изменяется после создания. |
| `vdc_uid` | `string` | Optional | UUID Услуги `Виртуальный датацентр (vDC)` Не изменяется после создания. |
| `virtual_services_count` | `number` | Optional, Computed | Кол-во виртуальных сервисов, которые **резервируются** на AVI |

### Nested attributes: `routed_net_configuration`

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `ip_addr_pool` | `string` | Optional, Computed | **Параметр не изменяемый**<br/>Адрессный пул, которая будет присвоена routed-сети. Описывается как (10.10.10.0/24). Маска 24 обязательна. (.1) - шлюз. (.2-.254) - Под адресацию для ВМ |
| `main_dns` | `string` | Optional, Computed | Подставляется автоматически в resolv.conf при создании виртуальных машин |
| `second_dns` | `string` | Optional, Computed | Подставляется автоматически в resolv.conf при создании виртуальных машин |
