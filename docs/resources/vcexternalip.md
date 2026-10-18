---
page_title: "nubes_vcexternalip"
description: "nubes_vcexternalip TEST resource schema for Nubes Cloud."
subcategory: "VMware"
---

# `nubes_vcexternalip` Resource

Service: Публичные IP адреса.

[Full TEST service guide](https://tf-docs.nodejsk8s.dev.nubes.ru/nubes-test/vcexternalip/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `adopt_existing_on_create` | `bool` | Optional, Computed | No | — |
| `dnat_create` | `bool` | Optional, Computed | No | Необходимо создать DNAT-правило. |
| `from_service_cloud_edge_name` | `string` | Optional, Computed | No | Edge составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_cloud_edge_scope` | `string` | Optional, Computed | No | Не изменяется после создания. |
| `from_service_cloud_org_name` | `string` | Optional, Computed | No | Org составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_cloud_vdc_name` | `string` | Optional, Computed | No | Vdc составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_cloud_vmware_url` | `string` | Optional, Computed | No | Api-Url составной услуги (указывать, если resourceRealm связан с CloudDirector) Не изменяется после создания. |
| `from_service_namespace` | `string` | Required | No | Namespace составной услуги (для отображения правил FW) Не изменяется после создания. |
| `from_service_vdc_group_name` | `string` | Optional, Computed | No | Не изменяется после создания. |
| `id` | `string` | Computed | No | — |
| `internal_addr_access` | `string` | Required | No | Внутренний адрес для связки ната. Требует dnatCreate=true или snatCreate=true |
| `internal_port_access` | `string` | Optional, Computed | No | Порт для связки ната с адресом. Несколько значений можно передавать через пробел. Если передать all, то dnat будет на все порты. Требует dnatCreate=true<br/> |
| `ip_space_name` | `string` | Required | No | Тип ip Space (internet-ipv4-v1 / internet-antiddos-v1 / internet-no-antiddos-v1) |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Режим state_only при destroy: ресурс не удаляется и не меняется в облаке, только убирается из состояния. Имеет приоритет над suspend_on_destroy. false = обычное удаление. |
| `log_level` | `string` | Optional | No | Operation stages log level: none (default), info, debug. Overrides provider-level log_level. |
| `operation_timeout` | `string` | Optional | No | — |
| `resource_name` | `string` | Required | No | — |
| `resource_realm` | `string` | Optional, Computed | No | Платформа для развертывания Не изменяется после создания. |
| `service_uid` | `string` | Required | No | — |
| `snat_create` | `bool` | Optional, Computed | No | Необходимо создать SNAT-правило.<br/> |
| `state_out` | `["map","string"]` | Computed | No | — |
| `state_out_flat` | `["map","string"]` | Computed | No | — |
| `state_params` | `["map","string"]` | Computed | No | — |
| `state_params_flat` | `["map","string"]` | Computed | No | — |
| `vault_fields` | `["list","string"]` | Computed | No | — |
| `vault_secrets` | `["map","string"]` | Computed | Yes | — |
| `vault_url` | `string` | Computed | No | — |
| `vault_user_path` | `string` | Computed | No | — |
