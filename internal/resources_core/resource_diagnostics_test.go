package resources_core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

func makeShouldRemoveClient(handler http.HandlerFunc) (*core.UniversalClient, func()) {
	srv := httptest.NewServer(handler)
	c := &core.UniversalClient{
		HttpClient:  srv.Client(),
		ApiEndpoint: srv.URL,
	}
	return c, func() { srv.Close() }
}

func TestShouldRemoveFromState(t *testing.T) {
	cases := []struct {
		name       string
		handler    http.HandlerFunc
		wantRemove bool
		wantErr    bool
	}{
		{
			name: "deleted статус",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"instance":{"instanceUid":"u1","explainedStatus":"deleted","isDeleted":true}}`))
			},
			wantRemove: true,
			wantErr:    false,
		},
		{
			name: "isDeleted=true без статуса deleted",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"instance":{"instanceUid":"u1","explainedStatus":"running","isDeleted":true}}`))
			},
			wantRemove: true,
			wantErr:    false,
		},
		{
			name: "HTTP 404 not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`"Instance not found or no permission to view"`))
			},
			wantRemove: true,
			wantErr:    false,
		},
		{
			name: "живой инстанс",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"instance":{"instanceUid":"u1","explainedStatus":"running","isDeleted":false}}`))
			},
			wantRemove: false,
			wantErr:    false,
		},
		{
			name: "HTTP 500 — НЕ удалять",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"message":"boom"}`))
			},
			wantRemove: false,
			wantErr:    true,
		},
		{
			name: "HTTP 403 — НЕ удалять",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`"permission denied"`))
			},
			wantRemove: false,
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, cleanup := makeShouldRemoveClient(tc.handler)
			defer cleanup()

			remove, err := ShouldRemoveFromState(context.Background(), client, "u1")
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected err=%v, got err=%v", tc.wantErr, err)
			}
			if remove != tc.wantRemove {
				t.Fatalf("expected remove=%v, got remove=%v", tc.wantRemove, remove)
			}
		})
	}
}

func TestShouldRemoveFromState_NilClient(t *testing.T) {
	remove, err := ShouldRemoveFromState(context.Background(), nil, "u1")
	if err != nil {
		t.Fatalf("expected nil error for nil client, got %v", err)
	}
	if remove {
		t.Fatalf("expected remove=false for nil client")
	}
}

func TestShouldRemoveFromState_EmptyID(t *testing.T) {
	client, cleanup := makeShouldRemoveClient(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("should not make request for empty id")
	})
	defer cleanup()

	remove, err := ShouldRemoveFromState(context.Background(), client, "  ")
	if err != nil {
		t.Fatalf("expected nil error for empty id, got %v", err)
	}
	if remove {
		t.Fatalf("expected remove=false for empty id")
	}
}
