package resources_core

import "testing"

// FIX(uuid-case): регрессия today's bug — adopt приостановленного инстанса Штурвала
// падал с "required params mismatch ... startupConfiguration", потому что UUID
// внутри JSON сравнивался с учётом регистра (plan — lowercase из config,
// actual — UPPERCASE от API). Теперь UUID внутри JSON эквивалентны.
func TestParamsEquivalentUUIDCaseInsideJSON(t *testing.T) {
	plan := `{"nsxtUid":"2c37fed1-e8f8-4a84-8434-7851c7c8b5d6","vdcUid":"d0937335-276b-475b-baa4-d8e6d16bad51","clusterName":"shturval-dev-00","exVip":true}`
	actual := `{"clusterName":"shturval-dev-00","exVip":true,"nsxtUid":"2C37FED1-E8F8-4A84-8434-7851C7C8B5D6","vdcUid":"d0937335-276b-475b-baa4-d8e6d16bad51"}`

	if !paramsEquivalent(plan, actual) {
		t.Fatalf("paramsEquivalent() = false, ожидалось true: регистр UUID внутри JSON не должен влиять")
	}

	// Реально другой nsxtUid — по-прежнему не эквивалентно.
	other := `{"clusterName":"shturval-dev-00","exVip":true,"nsxtUid":"86a01033-adac-424f-8296-aa10a001d6ed","vdcUid":"d0937335-276b-475b-baa4-d8e6d16bad51"}`
	if paramsEquivalent(plan, other) {
		t.Fatalf("paramsEquivalent() = true для другого nsxtUid — неверно")
	}
}
