---
page_title: "nubes_vc_nsxt_snat"
description: "Manages SNAT on an existing Cloud Director edge gateway."
subcategory: "VMware"
---

# `nubes_vc_nsxt_snat` Resource

Manages SNAT on an existing Cloud Director edge gateway.

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `id` | `string` | Computed | No | — |
| `ip_space_name` | `string` | Required | No | Имя ipSpace для внешнего IP (SNAT). Значение `no-needed` выключает SNAT. Имя должно быть выделено на организации (см. `nubes_vc_org_ip_allocation`). |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Не выключать SNAT при `destroy` (по умолчанию `false` — отправляется `ipSpaceName = "no-needed"`). |
| `nsxt_uid` | `string` | Required | No | UUID существующей услуги «Сетевой шлюз периметра (Edge)». |
