package provider

import _ "embed"

// operation_timeouts.json — ДЕФОЛТНЫЙ конфиг таймаутов операций для сборки из
// исходников (без него падает go build / go test в provider/).
//
// Это НЕ сгенерированный артефакт и он СОЗНАТЕЛЬНО лежит в git:
// при релизе скрипт TOOLS/scripts/03_build_and_upload_provider.sh копирует
// провайдер во временный каталог и ПОДМЕНЯЕТ там этот файл профильным
// (TOOLS/config/<profile>/operation_timeouts.json) — рабочее дерево не меняется.
//
// Менять значения для релиза нужно в TOOLS/config/<profile>/, а не здесь.
//
//go:embed operation_timeouts.json
var operationTimeoutsConfigEmbedded []byte
