package resources_core

import "testing"

// Регрессия на баг: canonicalParamsString раньше брала значение из map по
// lowercase-ключу, из-за чего значение всегда было пустым, а id подресурса
// (BuildSubresourceID) становился константой в пределах родителя.
func TestCanonicalParamsString(t *testing.T) {
	cases := []struct {
		name   string
		params map[string]string
		want   string
	}{
		{"empty", map[string]string{}, ""},
		{"single camelCase", map[string]string{"recordName": "www"}, "recordname=www"},
		{"sorted by normalized key", map[string]string{"bKey": "2", "aKey": "1"}, "akey=1;bkey=2"},
		{"value preserved", map[string]string{"recordName": "api"}, "recordname=api"},
	}
	for _, c := range cases {
		got := canonicalParamsString(c.params)
		if got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}

func TestBuildSubresourceID_uniquePerIdentity(t *testing.T) {
	a := BuildSubresourceID("uid", "record", map[string]string{"recordName": "www"})
	b := BuildSubresourceID("uid", "record", map[string]string{"recordName": "api"})
	if a == b {
		t.Errorf("id подресурса должен зависеть от identity-параметра, но совпал: %q", a)
	}
}
