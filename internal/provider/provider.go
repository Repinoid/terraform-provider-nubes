// Package provider — регистрация Terraform-провайдера Nubes.
//
// Определяет:
//   - New() — создание провайдера с конфигурацией (api_token, api_endpoint)
//   - Resources() — список всех ресурсов из resources_gen.AllResources()
//   - Schema() — схема конфигурации провайдера
package provider

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
	"github.com/repinoid/terraform-provider-nubes/internal/resources_core"
	"github.com/repinoid/terraform-provider-nubes/internal/resources_gen"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &NubesProvider{}

type NubesProvider struct {
	version string
}

// NubesProviderModel — конфигурация провайдера из HCL-блока provider {}.
// Каждое поле соответствует атрибуту в Schema() ниже.
// Поля Optional: если не указаны — берутся из env-переменных или defaults.
type NubesProviderModel struct {
	ApiEndpoint types.String `tfsdk:"api_endpoint"`
	ApiToken    types.String `tfsdk:"api_token"`
	// Insecure отключает проверку TLS-сертификата сервера.
	// НЕ использовать в продакшн! Только для dev-стендов с самоподписанным сертом.
	// Может быть задан также через env NUBES_INSECURE=true.
	Insecure types.Bool `tfsdk:"insecure"`
	// LogLevel задаёт глобальный уровень вывода этапов операции.
	// "none" (default) — молчать. "info" — имя этапа + длительность.
	// "debug" — info + детали подэтапов (StageMsg без timestamp-мусора).
	// Может быть переопределён атрибутом log_level в отдельном ресурсе.
	LogLevel types.String `tfsdk:"log_level"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NubesProvider{version: version}
	}
}

func (p *NubesProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nubes"
	resp.Version = p.version
}

func (p *NubesProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_endpoint": schema.StringAttribute{
				MarkdownDescription: "API Gateway endpoint for Nubes Cloud",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "API authentication token",
				Optional:            true,
				Sensitive:           true,
			},
			// insecure: отключает проверку TLS-сертификата Nubes API.
			// По умолчанию false — сертификат проверяется (безопасно).
			// Устанавливать true только на dev-стенде с самоподписанным сертом.
			// Альтернатива: env NUBES_INSECURE=true (не требует правки .tf файлов).
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Disable TLS certificate verification. Use only for dev environments with self-signed certs. Can also be set via NUBES_INSECURE env var.",
				Optional:            true,
			},
			// log_level: уровень вывода информации об этапах операций.
			// Значение применяется ко всем ресурсам, если не переопределено в ресурсе.
			"log_level": schema.StringAttribute{
				MarkdownDescription: "Logging level for operation stages. Values: `none` (default, silent), `info` (stage name + duration), `debug` (info + detailed sub-step messages). Can be overridden per resource.",
				Optional:            true,
			},
		},
	}
}

func (p *NubesProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config NubesProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ⛔ LEGACY: deck-api.ngcloud.ru ЗАКРЫВАЕТСЯ. Использовать ТОЛЬКО lk-api-gateway.
	// По умолчанию — Gateway. Переопределяется через NUBES_API_ENDPOINT.
	apiEndpoint := "https://lk-api-gateway.ngcloud.ru/api/v1/svc"
	apiToken := ""

	if !config.ApiEndpoint.IsNull() {
		apiEndpoint = config.ApiEndpoint.ValueString()
	}

	if !config.ApiToken.IsNull() {
		apiToken = strings.TrimSpace(config.ApiToken.ValueString())
	} else {
		if token := os.Getenv("NUBES_API_TOKEN"); token != "" {
			apiToken = strings.TrimSpace(token)
		}
	}

	// --- TLS: определяем нужно ли пропустить проверку сертификата ---
	// Приоритет: config.insecure > env NUBES_INSECURE > false (безопасный default).
	// Nubes Cloud API использует валидный TLS-сертификат от доверенного CA,
	// поэтому в продакшн InsecureSkipVerify должен быть false.
	// true оставлен только для совместимости с dev-стендами без нормального сертификата.
	insecureSkipVerify := false
	if !config.Insecure.IsNull() && !config.Insecure.IsUnknown() {
		// Явно задано в provider {} блоке HCL
		insecureSkipVerify = config.Insecure.ValueBool()
	} else if os.Getenv("NUBES_INSECURE") == "true" {
		// Задано через переменную окружения (удобно для CI/CD без правки .tf файлов)
		insecureSkipVerify = true
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   60 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecureSkipVerify,
			MinVersion:         tls.VersionTLS12,
		},
	}

	// log_level: оставлен для обратной совместимости, но этапы выводятся ВСЕГДА.
	// Значение используется только для StageMsg (debug).
	logLevel := "info"
	if !config.LogLevel.IsNull() && !config.LogLevel.IsUnknown() {
		if v := strings.TrimSpace(config.LogLevel.ValueString()); v != "" {
			logLevel = v
		}
	}

	client := &core.UniversalClient{
		HttpClient: &http.Client{
			Transport: transport,
			Timeout:   300 * time.Second,
		},
		ApiEndpoint:       apiEndpoint,
		ApiToken:          apiToken,
		ProviderVersion:   p.version,
		OperationTimeouts: core.DefaultOperationTimeouts(),
		LogLevel:          logLevel,
	}

	if len(operationTimeoutsConfigEmbedded) > 0 {
		timeouts, err := core.LoadOperationTimeoutsFromBytes(operationTimeoutsConfigEmbedded)
		if err != nil {
			resp.Diagnostics.AddError("Ошибка загрузки таймаутов операций", err.Error())
			return
		}
		client.OperationTimeouts = timeouts
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *NubesProvider) Resources(ctx context.Context) []func() resource.Resource {
	resources := resources_gen.AllResources()
	resources = append(resources, resources_core.NewServiceOperationResource)
	// Ресурсы-модификаторы для операций, которых нет в create-схеме ресурсов-инстансов.
	// Организация и шлюз создаются вручную в ЛК, поэтому адресуются по uid, а не ссылкой на ресурс.
	resources = append(resources, resources_core.NewOrgIpAllocationResource)
	resources = append(resources, resources_core.NewNsxtSnatResource)
	return resources
}

func (p *NubesProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
