package core

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNormalizeValue_EmptyString(t *testing.T) {
	param := universalCfsParam{DataType: "string"}
	result := normalizeUniversalValueV6("", param)
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}

func TestNormalizeValue_NullString(t *testing.T) {
	param := universalCfsParam{DataType: "string"}
	result := normalizeUniversalValueV6("null", param)
	if result != "" {
		t.Errorf("expected empty string for 'null' input, got: %q", result)
	}
}

func TestNormalizeValue_MapType(t *testing.T) {
	param := universalCfsParam{DataType: "map"}
	result := normalizeUniversalValueV6("", param)
	if result != "{}" {
		t.Errorf("expected '{}' for map dataType, got: %q", result)
	}
}

func TestNormalizeValue_ArrayType(t *testing.T) {
	param := universalCfsParam{DataType: "array"}
	result := normalizeUniversalValueV6("", param)
	if result != "[]" {
		t.Errorf("expected '[]' for array dataType, got: %q", result)
	}
}

func TestNormalizeValue_TrimSpace(t *testing.T) {
	param := universalCfsParam{DataType: "string"}
	result := normalizeUniversalValueV6("  hello  ", param)
	if result != "hello" {
		t.Errorf("expected 'hello' after trim, got: %q", result)
	}
}

func TestIsRetryable(t *testing.T) {
	retryable := []int{
		http.StatusUnauthorized,       // 401 — транзиентный отказ Gateway
		http.StatusTooManyRequests,    // 429
		http.StatusServiceUnavailable, // 503
		http.StatusBadGateway,         // 502
		http.StatusGatewayTimeout,     // 504
	}
	for _, code := range retryable {
		if !isRetryable(code) {
			t.Errorf("expected isRetryable(%d)=true", code)
		}
	}

	notRetryable := []int{
		http.StatusBadRequest,          // 400
		http.StatusForbidden,           // 403
		http.StatusNotFound,            // 404
		http.StatusInternalServerError, // 500
	}
	for _, code := range notRetryable {
		if isRetryable(code) {
			t.Errorf("expected isRetryable(%d)=false", code)
		}
	}
}

// Idempotent-пропуск НЕ должен создавать операцию: иначе останется «черновик»
// (POST /instanceOperations без run), и следующий modify будет ждать idle до таймаута.
func TestRunInstanceOperationByIdempotent_SkipsWithoutCreatingOperation(t *testing.T) {
	var opPosts atomic.Int32
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/instances/inst-1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"instance":{"instanceUid":"inst-1","serviceId":22,"explainedStatus":"running","isDeleted":false,"operationIsPending":false,"operationIsInProgress":false,"availableOperations":[{"svcOperationId":207,"operation":"modify"}],"state":{"params":{"ipSpaceName":"internet-ipv4-v1"}}}}`))
		case r.Method == http.MethodPost && r.URL.Path == "/instanceOperations":
			opPosts.Add(1)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"instanceOperationUid":"op-x"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})
	defer cleanup()

	err := c.RunInstanceOperationUniversalByIdempotent(context.Background(), "inst-1", "modify",
		map[string]string{"ipSpaceName": "internet-ipv4-v1"})
	if err != nil {
		t.Fatalf("expected nil error on idempotent skip, got: %v", err)
	}
	if n := opPosts.Load(); n != 0 {
		t.Fatalf("operation must NOT be created when live already matches, got %d POST(s)", n)
	}
}

// ===== Q5: partial-result tests (instance uid при ошибке после создания) =====
// Инстанс уже создан (POST /instances отдал Location), но следующая операция падает:
// uid должен вернуться ВМЕСТЕ с ошибкой, иначе облачный инстанс осиротеет.
func TestCreateGenericInstance_KeepsUIDWhenOperationCreateFails(t *testing.T) {
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/instances":
			w.Header().Set("Location", "./inst-9")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"instanceUid":"inst-9"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/instanceOperations":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"ERROR":"boom"}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})
	defer cleanup()

	uid, err := c.CreateGenericInstanceUniversalV6(context.Background(), 21, "vm-1", nil)
	if err == nil {
		t.Fatal("expected error from failed operation")
	}
	if uid != "inst-9" {
		t.Fatalf("expected instanceUid preserved on error, got %q", uid)
	}
}

// Инстанс НЕ создан: uid обязан быть пустым (иначе запишем в state несуществующий ресурс).
func TestCreateGenericInstance_EmptyUIDWhenInstanceCreateFails(t *testing.T) {
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ERROR":"bad request"}`))
	})
	defer cleanup()

	uid, err := c.CreateGenericInstanceUniversalV6(context.Background(), 21, "vm-1", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if uid != "" {
		t.Fatalf("expected empty uid before creation, got %q", uid)
	}
}

// ===== waitForOperationFinish tests =====
func makeTestClient(handler http.HandlerFunc) (*UniversalClient, func()) {
	srv := httptest.NewServer(handler)
	var buf bytes.Buffer
	c := &UniversalClient{
		HttpClient:   srv.Client(),
		ApiEndpoint:  srv.URL,
		PollInterval: 10 * time.Millisecond,
		StagesWriter: &buf,
	}
	return c, func() { srv.Close() }
}

func makeStagesJSON(stages []opStage, dtFinish string, isSuccessful bool, errorLog string) string {
	dtFinishField := "null"
	if dtFinish != "" {
		dtFinishField = `"` + dtFinish + `"`
	}
	isOk := "true"
	if !isSuccessful {
		isOk = "false"
	}
	stagesJSON := ""
	for i, s := range stages {
		if i > 0 {
			stagesJSON += ","
		}
		dtF := "null"
		if s.DtFinish != nil {
			dtF = `"` + *s.DtFinish + `"`
		}
		isS := "null"
		if s.DtFinish != nil {
			if s.IsSuccessful {
				isS = "true"
			} else {
				isS = "false"
			}
		}
		stagesJSON += `{"instanceOperationStageUid":"uid-` + s.InstanceOperationStageUid + `",` +
			`"stage":"` + s.Stage + `",` +
			`"isSuccessful":` + isS + `,` +
			`"dtFinish":` + dtF + `,` +
			`"duration":` + ftoa(s.Duration) + `}`
	}
	errLog := "null"
	if errorLog != "" {
		errLog = `"` + errorLog + `"`
	}
	return `{"instanceOperation":{"dtFinish":` + dtFinishField + `,"isSuccessful":` + isOk +
		`,"errorLog":` + errLog + `,"isInProgress":false,"isPending":false,"duration":1.5,"stages":[` + stagesJSON + `]}}`
}

func ftoa(f float64) string {
	s := strings.TrimRight(strings.TrimRight(
		strings.Replace(
			strings.Replace(
				time.Duration(f*float64(time.Second)).String(), "s", "", 1), "m", ".", 1),
		"0"), ".")
	if s == "" {
		return "0"
	}
	return s
}

func TestWaitForOperationFinish_AllStagesAtOnce(t *testing.T) {
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(makeStagesJSON([]opStage{
			{InstanceOperationStageUid: "1", Stage: "1. Validation", IsSuccessful: true, DtFinish: strPtr("2026-01-01T00:00:01Z"), Duration: 1.2},
			{InstanceOperationStageUid: "2", Stage: "2. Main", IsSuccessful: true, DtFinish: strPtr("2026-01-01T00:00:02Z"), Duration: 2.3},
		}, "2026-01-01T00:00:03Z", true, "")))
	})
	defer cleanup()

	err := c.waitForOperationFinish(context.Background(), "test-op-uid", 5*time.Second)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestWaitForOperationFinish_StagesSequential(t *testing.T) {
	var callCount atomic.Int32
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		n := callCount.Add(1)
		w.Header().Set("Content-Type", "application/json")
		switch n {
		case 1:
			// poll 1: stage 1 pending, operation in progress
			w.Write([]byte(makeStagesJSON([]opStage{
				{InstanceOperationStageUid: "1", Stage: "1. Validation", DtFinish: nil, Duration: 0.5},
			}, "", false, "")))
		case 2:
			// poll 2: stage 1 done, stage 2 pending
			w.Write([]byte(makeStagesJSON([]opStage{
				{InstanceOperationStageUid: "1", Stage: "1. Validation", IsSuccessful: true, DtFinish: strPtr("2026-01-01T00:00:01Z"), Duration: 1.2},
				{InstanceOperationStageUid: "2", Stage: "2. Main", DtFinish: nil, Duration: 0.3},
			}, "", false, "")))
		default:
			// poll 3: all done
			w.Write([]byte(makeStagesJSON([]opStage{
				{InstanceOperationStageUid: "1", Stage: "1. Validation", IsSuccessful: true, DtFinish: strPtr("2026-01-01T00:00:01Z"), Duration: 1.2},
				{InstanceOperationStageUid: "2", Stage: "2. Main", IsSuccessful: true, DtFinish: strPtr("2026-01-01T00:00:02Z"), Duration: 2.3},
			}, "2026-01-01T00:00:03Z", true, "")))
		}
	})
	defer cleanup()

	err := c.waitForOperationFinish(context.Background(), "test-op-uid", 5*time.Second)
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestWaitForOperationFinish_StageFails(t *testing.T) {
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(makeStagesJSON([]opStage{
			{InstanceOperationStageUid: "1", Stage: "1. Validation", IsSuccessful: true, DtFinish: strPtr("2026-01-01T00:00:01Z"), Duration: 1.2},
			{InstanceOperationStageUid: "2", Stage: "2. Main", IsSuccessful: false, DtFinish: strPtr("2026-01-01T00:00:02Z"), Duration: 0.1},
		}, "2026-01-01T00:00:02Z", false, "disk error")))
	})
	defer cleanup()

	err := c.waitForOperationFinish(context.Background(), "test-op-uid", 5*time.Second)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "disk error") {
		t.Errorf("expected error to contain 'disk error', got: %v", err)
	}
}

func TestWaitForOperationFinish_RetryOn503(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow retry test in short mode")
	}
	var callCount atomic.Int32
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		n := callCount.Add(1)
		if n <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(makeStagesJSON([]opStage{
			{InstanceOperationStageUid: "1", Stage: "1. Done", IsSuccessful: true, DtFinish: strPtr("2026-01-01T00:00:01Z"), Duration: 0.5},
		}, "2026-01-01T00:00:01Z", true, "")))
	})
	defer cleanup()

	// doRequest retry delays: 2s→4s→8s — используем короткий RetryBaseDelay для теста
	c.RetryBaseDelay = 10 * time.Millisecond

	err := c.waitForOperationFinish(context.Background(), "test-op-uid", 5*time.Second)
	if err != nil {
		t.Fatalf("expected nil error after retry, got: %v", err)
	}
}

func TestWaitForOperationFinish_Timeout(t *testing.T) {
	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// always in progress, never finishes
		w.Write([]byte(makeStagesJSON([]opStage{
			{InstanceOperationStageUid: "1", Stage: "1. Running", DtFinish: nil, Duration: 0.5},
		}, "", false, "")))
	})
	defer cleanup()

	err := c.waitForOperationFinish(context.Background(), "test-op-uid", 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "не завершилась за установленный таймаут") {
		t.Errorf("expected timeout message, got: %v", err)
	}
}

func TestCreateGenericInstanceUniversalV6_ContinuesWhenOpDetailsGETFails(t *testing.T) {
	var postedParams = make(map[int]string)
	var pollCount atomic.Int32

	c, cleanup := makeTestClient(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/instances":
			w.Header().Set("Location", "./inst-1")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"instanceUid":"inst-1"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/instanceOperations":
			w.Header().Set("Location", "./op-1")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{"instanceOperationUid":"op-1"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/instanceOperations/op-1" && strings.Contains(r.URL.RawQuery, "fields=cfsParams"):
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"boom"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/instanceOperationCfsParams":
			var req genericParamReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Fatalf("failed to decode param request: %v", err)
			}
			postedParams[req.SvcOperationCfsParamId] = req.ParamValue
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		case r.Method == http.MethodGet && r.URL.Path == "/instanceOperations/op-1/validate-cfs":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"ok":true}`))
		case r.Method == http.MethodPost && r.URL.Path == "/instanceOperations/op-1/run":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		case r.Method == http.MethodGet && r.URL.Path == "/instanceOperations/op-1" && strings.Contains(r.URL.RawQuery, "fields=dtFinish"):
			if pollCount.Add(1) != 1 {
				t.Fatalf("unexpected extra polling request: %s", r.URL.String())
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"instanceOperation":{"dtFinish":"2026-01-01T00:00:03Z","isSuccessful":true,"errorLog":null,"isInProgress":false,"isPending":false,"duration":1.5,"stages":[{"instanceOperationStageUid":"1","stage":"1. Validation","isSuccessful":true,"dtFinish":"2026-01-01T00:00:01Z","duration":1.0}]}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/instances/inst-1":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"instance":{"instanceUid":"inst-1","explainedStatus":"running","isDeleted":false,"operationIsPending":false,"operationIsInProgress":false}}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})
	defer cleanup()

	c.PollInterval = 10 * time.Millisecond
	c.RetryBaseDelay = 10 * time.Millisecond

	params := map[int]string{
		30:  "c3f81e4b-7a2b-4e1d-9e0a-123456789abc",
		335: "fast-2.8",
		361: `[{"name":"fast","size":200}]`,
		397: "0",
		557: "8",
		558: "32",
		746: "default",
	}

	if _, err := c.CreateGenericInstanceUniversalV6(context.Background(), 21, "f12vdc-2", params); err != nil {
		t.Fatalf("expected create to succeed, got error: %v", err)
	}

	if len(postedParams) != len(params) {
		t.Fatalf("expected %d posted params, got %d", len(params), len(postedParams))
	}
	for id, value := range params {
		if got := postedParams[id]; got != value {
			t.Fatalf("param %d mismatch: got %q, want %q", id, got, value)
		}
	}
	if pollCount.Load() != 1 {
		t.Fatalf("expected exactly 1 polling request, got %d", pollCount.Load())
	}
}

func strPtr(s string) *string { return &s }
