package resources_core

import "testing"

// Проверяем нормализацию vIPConfigure: свежая орга отдаёт `[{}]` (пустой элемент) —
// это «не выделено», а не аллокация (NOTES/30_analysis/HAR_FRESH_CREATE_2026-09-24.md).
func TestParseVipConfigure_EmptyAndBroken(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want int
	}{
		{"пустая строка", "", 0},
		{"пустой массив", "[]", 0},
		{"пустой элемент (свежая орга)", "[{}]", 0},
		{"только name без count", `[{"name":"internet-ipv4-v1"}]`, 1},
		{"элемент без name", `[{"count":"3"}]`, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseVipConfigure(tc.raw)
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if len(got) != tc.want {
				t.Fatalf("получено %d элементов, ожидалось %d (%+v)", len(got), tc.want, got)
			}
		})
	}
}

func TestParseVipConfigure_CountAsString(t *testing.T) {
	// count у платформы может прийти и числом (в старых стейтах) — нормализуем в строку.
	got, err := parseVipConfigure(`[{"name":"internet-ipv4-v1","count":4}]`)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if len(got) != 1 || got[0].Count != "4" {
		t.Fatalf("ожидался count=\"4\", получено %+v", got)
	}
}

func TestFormatVipConfigure_Canonical(t *testing.T) {
	got := formatVipConfigure([]vipAllocation{{Name: "internet-ipv4-v1", Count: "3"}})
	want := `[{"name":"internet-ipv4-v1","count":"3"}]`
	if got != want {
		t.Fatalf("получено %q, ожидалось %q", got, want)
	}
	if empty := formatVipConfigure(nil); empty != "[]" {
		t.Fatalf("для пустого списка ожидалось \"[]\", получено %q", empty)
	}
}

// importVipConfigure/importIpSpaceName — выбор значения Required-атрибутов при импорте
// (B7: Required не должен оставаться null после import).
func TestImportVipConfigure(t *testing.T) {
	cases := []struct {
		name string
		live map[string]string
		want string
	}{
		{"live-аллокация", map[string]string{"vIPConfigure": `[{"name":"internet-ipv4-v1","count":"3"}]`}, `[{"name":"internet-ipv4-v1","count":"3"}]`},
		{"нет ключа", map[string]string{}, "[]"},
		{"nil", nil, "[]"},
		{"свежая орга [{}]", map[string]string{"vIPConfigure": `[{}]`}, "[]"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := importVipConfigure(tc.live); got != tc.want {
				t.Fatalf("получено %q, ожидалось %q", got, tc.want)
			}
		})
	}
}

func TestImportIpSpaceName(t *testing.T) {
	cases := []struct {
		name string
		live map[string]string
		want string
	}{
		{"live-значение", map[string]string{"ipSpaceName": "internet-ipv4-v1"}, "internet-ipv4-v1"},
		{"нет ключа", map[string]string{}, "no-needed"},
		{"nil", nil, "no-needed"},
		{"только пробелы", map[string]string{"ipSpaceName": "   "}, "no-needed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := importIpSpaceName(tc.live); got != tc.want {
				t.Fatalf("получено %q, ожидалось %q", got, tc.want)
			}
		})
	}
}

// Сравнение смысловое: `jsonencode` сортирует ключи по алфавиту (count раньше name),
// но для нас это то же самое значение — переписывать state нельзя (иначе ложный дрейф).
func TestVipAllocationsEqual_OrderInsensitive(t *testing.T) {
	a, err := parseVipConfigure(`[{"name":"internet-ipv4-v1","count":"3"}]`)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	b, err := parseVipConfigure(`[{"count":"3","name":"internet-ipv4-v1"}]`) // так отдаёт jsonencode
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !vipAllocationsEqual(a, b) {
		t.Fatal("значения должны считаться равными несмотря на порядок ключей")
	}
}

func TestVipAllocationsEqual_Differences(t *testing.T) {
	base, _ := parseVipConfigure(`[{"name":"internet-ipv4-v1","count":"3"}]`)
	otherCount, _ := parseVipConfigure(`[{"name":"internet-ipv4-v1","count":"2"}]`)
	otherName, _ := parseVipConfigure(`[{"name":"internet-antiddos-v1","count":"3"}]`)
	empty, _ := parseVipConfigure(`[{}]`)

	if vipAllocationsEqual(base, otherCount) {
		t.Fatal("разный count должен считаться разными значениями")
	}
	if vipAllocationsEqual(base, otherName) {
		t.Fatal("разное имя ipSpace должно считаться разными значениями")
	}
	if vipAllocationsEqual(base, empty) {
		t.Fatal("пустая аллокация должна отличаться от непустой")
	}
}

func TestParseVipConfigure_RoundTripIsStable(t *testing.T) {
	raw := `[{"name":"internet-ipv4-v1","count":"4"}]`
	items, err := parseVipConfigure(raw)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if again := formatVipConfigure(items); again != raw {
		t.Fatalf("round-trip не стабилен: %q → %q", raw, again)
	}
}

func TestParseVipConfigure_InvalidJSON(t *testing.T) {
	if _, err := parseVipConfigure(`{"name":"x"}`); err == nil {
		t.Fatal("ожидалась ошибка на объект вместо массива")
	}
}
