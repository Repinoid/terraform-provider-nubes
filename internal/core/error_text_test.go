package core

import (
	"strings"
	"testing"
)

// TestFormatAPIError — текст ошибки API должен ВСЕГДА быть понятным юзеру.
//
// Платформа отдаёт минимум четыре разные формы тела (см. NOTES/30_analysis/
// ERROR_MATRIX_dns_zones_2026-10-16.md): если их не разбирать, в сообщение
// уйдёт сырой JSON. Тест фиксирует все формы + гарантию «текст есть всегда».
func TestFormatAPIError(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		body       string
		wantSubstr []string // все подстроки обязаны присутствовать
		wantAbsent []string // этих подстрок быть НЕ должно
	}{
		{
			name:       "RFC7807-подобная форма платформы + HTML-энтити",
			status:     400,
			body:       `{"DETAIL":"Invalid format integer &gt; 0 recordTTL:","TYPE":"about:blank","TITLE":"Bad Request"}`,
			wantSubstr: []string{"400", "Bad Request", "Invalid format integer > 0 recordTTL"},
			wantAbsent: []string{"DETAIL", "&gt;", "TYPE"},
		},
		{
			name:       "шлюз: вложенный error{code,message}",
			status:     401,
			body:       `{"error":{"code":"AUTHENTICATION_REQUIRED","message":"Требуется авторизация."},"requestId":"abc"}`,
			wantSubstr: []string{"401", "AUTHENTICATION_REQUIRED", "Требуется авторизация."},
			wantAbsent: []string{"requestId"},
		},
		{
			name:       "legacy-прокси: ERROR + DETAIL",
			status:     500,
			body:       `{"ERROR":"Internal error","DETAIL":"boom"}`,
			wantSubstr: []string{"500", "Internal error", "boom"},
		},
		{
			name:       "нижний регистр: message",
			status:     400,
			body:       `{"message":"Поле обязательно"}`,
			wantSubstr: []string{"400", "Поле обязательно"},
		},
		{
			name:       "список типов с энтити",
			status:     400,
			body:       `{"TITLE":"Bad Request","DETAIL":"Значение отсутствует в списке &#39;A,AAAA&#39;"}`,
			wantSubstr: []string{"Значение отсутствует в списке 'A,AAAA'"},
			wantAbsent: []string{"&#39;"},
		},
		{
			name:       "не JSON — сырое тело сохраняем",
			status:     502,
			body:       `Bad Gateway <html>`,
			wantSubstr: []string{"502", "Bad Gateway <html>"},
		},
		{
			name:       "пустое тело — берём текст статуса",
			status:     503,
			body:       ``,
			wantSubstr: []string{"503", "Service Unavailable"},
		},
		{
			name:       "массив строк errors",
			status:     400,
			body:       `{"errors":["поле A пустое","поле B пустое"]}`,
			wantSubstr: []string{"400", "поле A пустое", "поле B пустое"},
		},
		{
			name:       "массив объектов errors",
			status:     400,
			body:       `{"errors":[{"message":"Неверный формат"}]}`,
			wantSubstr: []string{"400", "Неверный формат"},
		},
		{
			name:       "Google-style error.errors",
			status:     403,
			body:       `{"error":{"errors":[{"message":"Доступ запрещён"}]}}`,
			wantSubstr: []string{"403", "Доступ запрещён"},
		},
		{
			name:       "заглавная Message (регистронезависимый fallback)",
			status:     400,
			body:       `{"Message":"Поле обязательно"}`,
			wantSubstr: []string{"400", "Поле обязательно"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := formatAPIError(tc.status, []byte(tc.body))
			if err == nil {
				t.Fatalf("ожидалась ошибка, получен nil")
			}
			got := err.Error()
			if strings.TrimSpace(got) == "" {
				t.Fatalf("текст ошибки ПУСТ — юзер не увидит причину")
			}
			for _, want := range tc.wantSubstr {
				if !strings.Contains(got, want) {
					t.Errorf("в тексте нет %q:\n  получено: %s", want, got)
				}
			}
			for _, absent := range tc.wantAbsent {
				if strings.Contains(got, absent) {
					t.Errorf("в тексте лишнее %q:\n  получено: %s", absent, got)
				}
			}
		})
	}
}

// TestParamLabel — в текстах ошибок должен быть код параметра из спеки,
// а не бесполезный числовой ID.
func TestParamLabel(t *testing.T) {
	params := []universalCfsParam{
		{SvcOperationCfsParamId: 779, Code: "recordType", Name: "Тип записи"},
		{SvcOperationCfsParamId: 778, Code: "recordTTL"},
		{SvcOperationCfsParamId: 777, Name: "Только имя"},
	}

	if got := paramLabelByID(params, 779); got != "recordType (Тип записи)" {
		t.Errorf("779: получено %q", got)
	}
	if got := paramLabelByID(params, 778); got != "recordTTL" {
		t.Errorf("778: получено %q", got)
	}
	if got := paramLabelByID(params, 777); got != "Только имя" {
		t.Errorf("777: получено %q", got)
	}
	// Неизвестный ID — честный fallback, а не пустая строка.
	if got := paramLabelByID(params, 999); got != "id 999" {
		t.Errorf("999: получено %q", got)
	}
	// Совпадающие code/name не дублируем.
	if got := paramLabel(universalCfsParam{SvcOperationCfsParamId: 1, Code: "cpu", Name: "cpu"}); got != "cpu" {
		t.Errorf("code==name: получено %q", got)
	}
}
