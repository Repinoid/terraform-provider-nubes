// Terraform-провайдер Nubes Cloud.
//
// Регистрирует провайдер в hashicorp/terraform-plugin-framework.
// Все ресурсы подключаются через resources_gen.AllResources().
//
// Сборка: go build -o terraform-provider-nubes .
// Версия задаётся через ldflags: -ldflags "-X main.version=X.Y.Z"
package main

import (
	"context"
	"log"

	"github.com/repinoid/terraform-provider-nubes/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// address задаётся через ldflags при сборке: -ldflags "-X main.address=host/ns/name"
// Значение по умолчанию — только для dev-сборок без скрипта.
var (
	version string = "0.0.0-dev"
	address string = "registry.terraform.io/repinoid/nubes"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: address,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err)
	}
}
