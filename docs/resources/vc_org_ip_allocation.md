---
page_title: "nubes_vc_org_ip_allocation"
description: "Manages external IP allocation on an existing Cloud Director organization."
subcategory: "VMware"
---

# `nubes_vc_org_ip_allocation` Resource

Manages external IP allocation on an existing Cloud Director organization.

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `id` | `string` | Computed | No | — |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Не снимать аллокацию IP при `destroy` (по умолчанию `false` — квота обнуляется, `count=0` по каждому элементу). |
| `organization` | `string` | Required | No | Организация, на которой выделяются внешние IP: имя из ЛК (например `organ`) или её UUID. |
| `vip_configure` | `string` | Required | No | JSON-массив аллокаций: `[{"name":"internet-ipv4-v1","count":"3"}]`. Значение перезаписывает текущую аллокацию целиком. `count` — строка. Порядок ключей и форматирование не важны (сравнение смысловое). Снять аллокацию (`[]`) через этот атрибут **нельзя** — только удалением ресурса (`destroy`). |
