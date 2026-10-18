package resources_core

import (
	"context"
	"sort"
	"strings"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

type RequiredParamDiff struct {
	Desired string
	Actual  string
}

func RequiredParamsMismatch(ctx context.Context, client *core.UniversalClient, serviceID int, desiredParamsByID map[int]string, instanceUID string) (map[string]RequiredParamDiff, error) {
	if client == nil || serviceID <= 0 || strings.TrimSpace(instanceUID) == "" {
		return nil, nil
	}
	if len(desiredParamsByID) == 0 {
		return nil, nil
	}

	desiredParamsByID = client.ResolveS3UidInAllMapFixed(desiredParamsByID, nil)

	requiredIDs, err := getServiceRequiredParamIDs(serviceID)
	if err != nil {
		return nil, err
	}
	if len(requiredIDs) == 0 {
		return nil, nil
	}

	codeMap, err := getServiceParamCodeMap(serviceID)
	if err != nil {
		return nil, err
	}
	refMap, err := getServiceParamRefMap(serviceID)
	if err != nil {
		return nil, err
	}

	desiredByCode := make(map[string]string)
	for id, val := range desiredParamsByID {
		if !requiredIDs[id] {
			continue
		}
		code := codeMap[id]
		if strings.TrimSpace(code) == "" {
			continue
		}
		if refSvcId, ok := refMap[id]; ok && refSvcId > 0 {
			resolved, err := client.ResolveRefSvcParamValue(ctx, refSvcId, val)
			if err != nil {
				return nil, err
			}
			if strings.TrimSpace(resolved) != "" {
				val = resolved
			}
		}
		norm := normalizeCompareValue(val)
		if norm == "" {
			continue
		}
		desiredByCode[code] = norm
	}
	if len(desiredByCode) == 0 {
		return nil, nil
	}

	actualParams, err := client.GetInstanceStateParams(ctx, instanceUID)
	if err != nil {
		return nil, err
	}

	mismatches := make(map[string]RequiredParamDiff)
	for code, desired := range desiredByCode {
		actual := normalizeCompareValue(actualParams[code])
		if paramsEquivalent(desired, actual) {
			continue
		}
		mismatches[code] = RequiredParamDiff{Desired: desired, Actual: actual}
	}
	if len(mismatches) == 0 {
		return nil, nil
	}
	return mismatches, nil
}

func paramsEquivalent(desired string, actual string) bool {
	if desired == actual {
		return true
	}
	if JSONStringsEquivalent(desired, actual) {
		return true
	}
	return false
}

func FormatRequiredParamMismatchMessage(mismatches map[string]RequiredParamDiff) string {
	if len(mismatches) == 0 {
		return ""
	}
	keys := make([]string, 0, len(mismatches))
	for code := range mismatches {
		keys = append(keys, code)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, code := range keys {
		diff := mismatches[code]
		desired := diff.Desired
		actual := diff.Actual
		if strings.TrimSpace(desired) == "" {
			desired = "<empty>"
		}
		if strings.TrimSpace(actual) == "" {
			actual = "<empty>"
		}
		parts = append(parts, code+" (plan="+desired+", actual="+actual+")")
	}
	return strings.Join(parts, "; ")
}
