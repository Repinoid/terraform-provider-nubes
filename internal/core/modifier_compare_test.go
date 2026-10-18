package core

import (
	"context"
	"net/http"
	"testing"
)

func TestModifierRawValuesEqual_Scalars(t *testing.T) {
	tests := []struct {
		name            string
		wanted, current string
		want            bool
	}{
		{"bool совпал", "false", "false", true},
		{"bool регистр", "True", "true", true},
		{"bool расхождение", "true", "false", false},
		{"int с пробелом", " 3 ", "3", true},
		{"int расхождение", "1", "3", false},
		{"null и пустое", "null", "", true},
		{"string совпал", "QoS-100Mbit", "QoS-100Mbit", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := modifierRawValuesEqual(tt.wanted, tt.current); got != tt.want {
				t.Errorf("modifierRawValuesEqual(%q, %q) = %v, want %v", tt.wanted, tt.current, got, tt.want)
			}
		})
	}
}

func TestModifierRawValuesEqual_JSON(t *testing.T) {
	// JSON: порядок ключей не значим
	if !modifierRawValuesEqual(
		`{"mainDns":"81.22.46.22","ipAddrPool":"10.10.102.0/24"}`,
		`{"ipAddrPool":"10.10.102.0/24","mainDns":"81.22.46.22"}`) {
		t.Errorf("JSON с разным порядком ключей должен считаться равным")
	}
	// JSON-массив: порядок элементов значим
	if !modifierRawValuesEqual(`[{"name":"a"}]`, `[{"name":"a"}]`) {
		t.Errorf("эквивалентные JSON-массивы должны быть равны")
	}
	if modifierRawValuesEqual(`[{"name":"a"},{"name":"b"}]`, `[{"name":"b"},{"name":"a"}]`) {
		t.Errorf("JSON-массивы с другим порядком не должны быть равны")
	}
}

func TestModifierDesiredEqualsLive(t *testing.T) {
	newClient := func(liveParams string) (*UniversalClient, func()) {
		return makeTestClient(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"instance":{"state":{"params":` + liveParams + `}}}`))
		})
	}

	t.Run("совпало с live", func(t *testing.T) {
		c, cleanup := newClient(`{"ipSpaceName":"internet-ipv4-v1"}`)
		defer cleanup()
		equal, err := c.modifierDesiredEqualsLive(context.Background(), "inst-1", map[string]string{"ipSpaceName": "internet-ipv4-v1"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if !equal {
			t.Errorf("expected equal when live matches desired")
		}
	})

	t.Run("отличается от live", func(t *testing.T) {
		c, cleanup := newClient(`{"ipSpaceName":"no-needed"}`)
		defer cleanup()
		equal, err := c.modifierDesiredEqualsLive(context.Background(), "inst-1", map[string]string{"ipSpaceName": "internet-ipv4-v1"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if equal {
			t.Errorf("expected not equal when live differs")
		}
	})

	t.Run("код не найден в live — не пропускаем", func(t *testing.T) {
		c, cleanup := newClient(`{"otherParam":"x"}`)
		defer cleanup()
		equal, err := c.modifierDesiredEqualsLive(context.Background(), "inst-1", map[string]string{"ipSpaceName": "internet-ipv4-v1"})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if equal {
			t.Errorf("must not report equal when code is absent in live")
		}
	})

	t.Run("live недоступен — ошибка, не пропускаем", func(t *testing.T) {
		c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer cleanup()
		equal, err := c.modifierDesiredEqualsLive(context.Background(), "inst-1", map[string]string{"ipSpaceName": "internet-ipv4-v1"})
		if err == nil {
			t.Errorf("expected error when live is unavailable")
		}
		if equal {
			t.Errorf("must not report equal when live is unavailable")
		}
	})
}
