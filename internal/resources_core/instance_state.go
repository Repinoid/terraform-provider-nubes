package resources_core

import (
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

// InstanceState — enum всех возможных состояний инстанса (16).
// Используется вместо strings.Contains для устойчивости к изменению текста API.
type InstanceState int

const (
	StateUnknown            InstanceState = iota
	StateNotCreated                       // 1
	StateCreating                         // 2
	StateCreationFailed                   // 3
	StateRunning                          // 4
	StateRunningPending                   // 5
	StateModifying                        // 6
	StateModificationFailed               // 7
	StateSuspending                       // 8
	StateSuspendFailed                    // 9
	StateSuspended                        // 10
	StateResuming                         // 11
	StateResumeFailed                     // 12
	StateDeleting                         // 13
	StateDeletionFailed                   // 14
	StateDeleted                          // 15
	StateOrphaned                         // 16
)

// classifyStatus определяет состояние инстанса по данным API.
// Использует комбинацию булевых флагов и explainedStatus из API.
// Централизованная классификация — единственное место где матчится текст статуса.
func classifyStatus(inst *core.InstanceStateResponse) InstanceState {
	id := inst.IsDeleted
	ip := inst.OperationIsInProgress
	op := inst.OperationIsPending
	status := strings.ToLower(strings.TrimSpace(inst.ExplainedStatus))

	// 15: DELETED — явный флаг isDeleted
	if id {
		return StateDeleted
	}

	// 10: SUSPENDED — статус suspended + нет операции
	if strings.Contains(status, "suspended") && !ip && !op {
		return StateSuspended
	}

	// 8: SUSPENDING — статус suspended + операция в процессе
	if strings.Contains(status, "suspend") && ip {
		return StateSuspending
	}

	// 11: RESUMING — suspended + операция в процессе (resume)
	if strings.Contains(status, "suspended") && ip {
		return StateResuming
	}

	// 1: NOT_CREATED
	if strings.Contains(status, "not created") {
		return StateNotCreated
	}

	// 2: CREATING
	if strings.Contains(status, "creating") && ip {
		return StateCreating
	}

	// 4: RUNNING
	if strings.Contains(status, "running") && !ip && !op {
		return StateRunning
	}

	// 5: RUNNING_PENDING
	if strings.Contains(status, "running") && op {
		return StateRunningPending
	}

	// 6: MODIFYING
	if strings.Contains(status, "modif") && ip {
		return StateModifying
	}

	// 13: DELETING
	if strings.Contains(status, "delet") && ip {
		return StateDeleting
	}

	// Ошибки — по тексту статуса
	if strings.Contains(status, "failed") || strings.Contains(status, "error") {
		switch {
		case strings.Contains(status, "creation") || strings.Contains(status, "creat"):
			return StateCreationFailed
		case strings.Contains(status, "modif"):
			return StateModificationFailed
		case strings.Contains(status, "suspend"):
			return StateSuspendFailed
		case strings.Contains(status, "resume"):
			return StateResumeFailed
		case strings.Contains(status, "delet"):
			return StateDeletionFailed
		default:
			return StateCreationFailed // fallback
		}
	}

	return StateUnknown
}

// isAdoptable проверяет, можно ли adopt-ить инстанс в данном состоянии.
func isAdoptable(state InstanceState) bool {
	switch state {
	case StateRunning, StateRunningPending, StateSuspended:
		return true
	default:
		return false
	}
}

// isTerminalError проверяет, является ли состояние терминальной ошибкой.
func isTerminalError(state InstanceState) bool {
	switch state {
	case StateCreationFailed, StateModificationFailed, StateSuspendFailed,
		StateResumeFailed, StateDeletionFailed, StateOrphaned, StateDeleted:
		return true
	default:
		return false
	}
}
