package resources_core

import (
	"context"
	"testing"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

// ── Gate «писать ли id при ошибке create» (isCreated) ─────────────────────────

func TestShouldKeepIDOnCreateError(t *testing.T) {
	trueVal := true
	falseVal := false
	cases := []struct {
		name     string
		state    *core.InstanceStateResponse
		stateErr error
		expected bool
	}{
		{
			name:     "инстанс не создан (полный откат) — id НЕ сохраняем",
			state:    &core.InstanceStateResponse{IsCreated: &falseVal, ExplainedStatus: "not created"},
			expected: false,
		},
		{
			name:     "инстанс создан — id сохраняем",
			state:    &core.InstanceStateResponse{IsCreated: &trueVal},
			expected: true,
		},
		{
			name:     "legacy-эндпоинт без isCreated — поведение как раньше (сохраняем)",
			state:    &core.InstanceStateResponse{IsCreated: nil},
			expected: true,
		},
		{
			name:     "ошибка чтения состояния (IAM 408/таймаут) — как раньше (сохраняем)",
			state:    nil,
			stateErr: context.DeadlineExceeded,
			expected: true,
		},
		{
			name:     "пустое состояние без ошибки — как раньше (сохраняем)",
			state:    nil,
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := shouldKeepIDOnCreateError(c.state, c.stateErr)
			if got != c.expected {
				t.Errorf("shouldKeepIDOnCreateError(): expected %v, got %v", c.expected, got)
			}
		})
	}
}

func TestIsStatusSuspended(t *testing.T) {
	cases := []struct {
		status   string
		expected bool
	}{
		{"suspended", true},
		{"suspend", true},
		{"running", false},
		{"deleted", false},
		{"", false},
	}
	for _, c := range cases {
		result := isStatusSuspended(c.status)
		if result != c.expected {
			t.Errorf("isStatusSuspended(%q): expected %v, got %v", c.status, c.expected, result)
		}
	}
}

func TestIsStatusNonAdoptable(t *testing.T) {
	cases := []struct {
		status   string
		expected bool
	}{
		{"", true},
		{"not created", true},
		{"pending", true},
		{"creating", true},
		{"failed", true},
		{"error", true},
		{"running", false},
		{"suspended", false},
	}
	for _, c := range cases {
		result := isStatusNonAdoptable(c.status)
		if result != c.expected {
			t.Errorf("isStatusNonAdoptable(%q): expected %v, got %v", c.status, c.expected, result)
		}
	}
}

func TestDeleteBehaviorDefault(t *testing.T) {
	// Empty destroyBehavior defaults to "state_only" which returns nil immediately
	// without making any API calls, so nil client is safe here.
	err := DeleteResource(context.Background(), nil, "test-uid", "")
	if err != nil {
		t.Errorf("expected nil error for default (state_only) destroy behavior, got: %v", err)
	}
}

// ── P1.6: Табличные тесты classifyStatus ────────────────────────────────────

func TestClassifyStatus(t *testing.T) {
	cases := []struct {
		name     string
		inst     *core.InstanceStateResponse
		expected InstanceState
	}{
		{
			name:     "NOT_CREATED",
			inst:     &core.InstanceStateResponse{ExplainedStatus: "not created"},
			expected: StateNotCreated,
		},
		{
			name: "CREATING",
			inst: &core.InstanceStateResponse{
				ExplainedStatus:       "creating",
				OperationIsInProgress: true,
			},
			expected: StateCreating,
		},
		{
			name: "RUNNING",
			inst: &core.InstanceStateResponse{
				ExplainedStatus: "running",
			},
			expected: StateRunning,
		},
		{
			name: "RUNNING_PENDING",
			inst: &core.InstanceStateResponse{
				ExplainedStatus:    "running",
				OperationIsPending: true,
			},
			expected: StateRunningPending,
		},
		{
			name: "MODIFYING",
			inst: &core.InstanceStateResponse{
				ExplainedStatus:       "modifying",
				OperationIsInProgress: true,
			},
			expected: StateModifying,
		},
		{
			name: "SUSPENDED",
			inst: &core.InstanceStateResponse{
				ExplainedStatus: "suspended",
			},
			expected: StateSuspended,
		},
		{
			name: "SUSPENDING",
			inst: &core.InstanceStateResponse{
				ExplainedStatus:       "suspending",
				OperationIsInProgress: true,
			},
			expected: StateSuspending,
		},
		{
			name: "DELETED",
			inst: &core.InstanceStateResponse{
				ExplainedStatus: "deleted",
				IsDeleted:       true,
			},
			expected: StateDeleted,
		},
		{
			name: "DELETING",
			inst: &core.InstanceStateResponse{
				ExplainedStatus:       "deleting",
				OperationIsInProgress: true,
			},
			expected: StateDeleting,
		},
		{
			name: "CREATION_FAILED",
			inst: &core.InstanceStateResponse{
				ExplainedStatus: "creation failed",
			},
			expected: StateCreationFailed,
		},
		{
			name: "UNKNOWN (empty)",
			inst: &core.InstanceStateResponse{
				ExplainedStatus: "",
			},
			expected: StateUnknown,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := classifyStatus(c.inst)
			if result != c.expected {
				t.Errorf("classifyStatus(%q): expected %d, got %d",
					c.inst.ExplainedStatus, c.expected, result)
			}
		})
	}
}

func TestIsAdoptable(t *testing.T) {
	adoptable := []InstanceState{StateRunning, StateRunningPending, StateSuspended}
	notAdoptable := []InstanceState{StateNotCreated, StateCreating, StateCreationFailed,
		StateModifying, StateModificationFailed, StateSuspending, StateSuspendFailed,
		StateResuming, StateResumeFailed, StateDeleting, StateDeletionFailed,
		StateDeleted, StateOrphaned, StateUnknown}

	for _, s := range adoptable {
		if !isAdoptable(s) {
			t.Errorf("state %d should be adoptable", s)
		}
	}
	for _, s := range notAdoptable {
		if isAdoptable(s) {
			t.Errorf("state %d should NOT be adoptable", s)
		}
	}
}

func TestIsTerminalError(t *testing.T) {
	errors := []InstanceState{StateCreationFailed, StateModificationFailed,
		StateSuspendFailed, StateResumeFailed, StateDeletionFailed,
		StateOrphaned, StateDeleted}
	nonErrors := []InstanceState{StateNotCreated, StateCreating, StateRunning,
		StateRunningPending, StateModifying, StateSuspending, StateSuspended,
		StateResuming, StateDeleting, StateUnknown}

	for _, s := range errors {
		if !isTerminalError(s) {
			t.Errorf("state %d should be terminal error", s)
		}
	}
	for _, s := range nonErrors {
		if isTerminalError(s) {
			t.Errorf("state %d should NOT be terminal error", s)
		}
	}
}
