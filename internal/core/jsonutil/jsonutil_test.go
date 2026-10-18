package jsonutil

import "testing"

// FIX(uuid-case) регрессия: облако отдаёт один и тот же UUID то в lowercase (vdcUid),
// то в UPPERCASE (nsxtUid у кластера Штурвала). Сравнение plan vs actual должно
// игнорировать регистр UUID — и для одиночного значения, и внутри JSON.
func TestJSONStringsEquivalentUUIDCase(t *testing.T) {
	cases := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{
			name: "uuid внутри JSON: разный регистр — эквивалентно",
			a:    `{"nsxtUid":"2c37fed1-e8f8-4a84-8434-7851c7c8b5d6","clusterName":"shturval-dev-00"}`,
			b:    `{"clusterName":"shturval-dev-00","nsxtUid":"2C37FED1-E8F8-4A84-8434-7851C7C8B5D6"}`,
			want: true,
		},
		{
			name: "вложенный JSON и смешанный регистр",
			a:    `{"a":{"b":["2C37FED1-E8F8-4A84-8434-7851C7C8B5D6"]}}`,
			b:    `{"a":{"b":["2c37fed1-e8f8-4a84-8434-7851c7c8b5d6"]}}`,
			want: true,
		},
		{
			name: "разные UUID — не эквивалентно",
			a:    `{"nsxtUid":"2c37fed1-e8f8-4a84-8434-7851c7c8b5d6"}`,
			b:    `{"nsxtUid":"86a01033-adac-424f-8296-aa10a001d6ed"}`,
			want: false,
		},
		{
			name: "разные значения не-UUID — не эквивалентно",
			a:    `{"clusterName":"shturval-dev-00"}`,
			b:    `{"clusterName":"shturval-dev-01"}`,
			want: false,
		},
		{
			name: "числа и bool остаются различимыми",
			a:    `{"exVip":true,"count":3}`,
			b:    `{"exVip":true,"count":"3"}`,
			want: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := JSONStringsEquivalent(tc.a, tc.b); got != tc.want {
				t.Fatalf("JSONStringsEquivalent(%s, %s) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestLowercaseUUIDsInText(t *testing.T) {
	in := `{"nsxtUid":"2C37FED1-E8F8-4A84-8434-7851C7C8B5D6","name":"internet-ipv4-v1"}`
	want := `{"nsxtUid":"2c37fed1-e8f8-4a84-8434-7851c7c8b5d6","name":"internet-ipv4-v1"}`
	if got := LowercaseUUIDsInText(in); got != want {
		t.Fatalf("LowercaseUUIDsInText() = %s, want %s", got, want)
	}

	// Текст без UUID не меняется (в т.ч. имя ipSpace с дефисами и цифрами).
	plain := `{"ipSpaceName":"internet-ipv4-v1","qosProfile":"QoS-100Mbit"}`
	if got := LowercaseUUIDsInText(plain); got != plain {
		t.Fatalf("LowercaseUUIDsInText() изменил текст без UUID: %s", got)
	}
}
