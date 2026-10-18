package resources_core

import "github.com/repinoid/terraform-provider-nubes/internal/core/jsonutil"

// JSONStringsEquivalent compares JSON strings ignoring whitespace and map key order.
//
// Реэкспорт каноничной реализации из core/jsonutil: функции вынесены в нейтральный
// пакет, чтобы их мог импортировать и core (без цикла core → resources_core → core).
// Оставляем обёртку, чтобы не трогать существующие вызовы в resources_core.
func JSONStringsEquivalent(a, b string) bool {
	return jsonutil.JSONStringsEquivalent(a, b)
}
