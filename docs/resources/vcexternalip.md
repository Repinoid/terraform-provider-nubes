---
page_title: "nubes_vcexternalip"
description: "nubes_vcexternalip TEST resource schema for Nubes Cloud."
subcategory: "VMware"
---

# `nubes_vcexternalip` Resource

Service: Публичные IP адреса.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/vcexternalip/).

## Schema

| Name | Type | Plan behavior | Details |
|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | — |
| `dnat_create` | `bool` | Optional, Computed | Необходимо создать DNAT-правило. |
| `from_service_cloud_edge_name` | `string` | Optional, Computed | Edge составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_cloud_edge_scope` | `string` | Optional, Computed | Не изменяется после создания. |
| `from_service_cloud_org_name` | `string` | Optional, Computed | Org составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_cloud_vdc_name` | `string` | Optional, Computed | Vdc составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_cloud_vmware_url` | `string` | Optional, Computed | Api-Url составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_namespace` | `string` | Required | Namespace составной услуги (для отображения правил FW) Не изменяется после создания. |
| `from_service_vdc_group_name` | `string` | Optional, Computed | Не изменяется после создания. |
| `id` | `string` | Computed | — |
| `internal_addr_access` | `string` | Required | Внутренний адрес для связки ната. Требует dnatCreate=true или snatCreate=true |
| `internal_port_access` | `string` | Optional, Computed | Порт для связки ната с адресом. Несколько значений можно передавать через пробел. Если передать all, то dnat будет на все порты. Требует dnatCreate=true<br/> |
| `ip_space_name` | `string` | Required | Тип ip Space (internet-ipv4-v1 / internet-antiddos-v1 / internet-no-antiddos-v1) |
| `keep_on_destroy` | `bool` | Optional, Computed | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | — |
| `resource_name` | `string` | Required | — |
| `resource_realm` | `string` | Optional, Computed | Платформа для развертывания Не изменяется после создания. |
| `service_uid` | `string` | Required | — |
| `snat_create` | `bool` | Optional, Computed | Необходимо создать SNAT-правило.<br/> |
| `state_out` | `["map","string"]` | Computed | — |
| `state_out_flat` | `["map","string"]` | Computed | — |
| `state_params` | `["map","string"]` | Computed | — |
| `state_params_flat` | `["map","string"]` | Computed | — |
| `vault_fields` | `["list","string"]` | Computed | — |
| `vault_secrets` | `["map","string"]` | Computed | **Sensitive.**  |
| `vault_url` | `string` | Computed | — |
| `vault_user_path` | `string` | Computed | — |
