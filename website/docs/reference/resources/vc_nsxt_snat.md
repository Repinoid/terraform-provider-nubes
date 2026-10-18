# `nubes_vc_nsxt_snat` Resource

Service: Сетевой шлюз периметра (Edge).

[Full service guide](https://repinoid.github.io/terraform-provider-nubes/resources/vc_nsxt/).

[Service manual](https://repinoid.github.io/terraform-provider-nubes/guides/services/vc_nsxt/).

## Schema

| Name | Type | Plan behavior | Sensitive | Description |
|---|---|---|---|---|
| `id` | `string` | Computed | No | — |
| `ip_space_name` | `string` | Required | No | Имя ipSpace для внешнего IP (SNAT). Значение `no-needed` выключает SNAT. Имя должно быть выделено на организации (см. `nubes_vc_org_ip_allocation`). |
| `keep_on_destroy` | `bool` | Optional, Computed | No | Не выключать SNAT при `destroy` (по умолчанию `false` — отправляется `ipSpaceName = "no-needed"`). |
| `nsxt_uid` | `string` | Required | No | UUID существующей услуги «Сетевой шлюз периметра (Edge)». |
