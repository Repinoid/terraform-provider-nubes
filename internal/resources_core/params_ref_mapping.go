package resources_core

import (
	"context"
	"io/fs"
	"strings"
	"sync"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
	resourcesyaml "github.com/repinoid/terraform-provider-nubes/resources_yaml"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"gopkg.in/yaml.v3"
)

type yamlParamRef struct {
	ID       int `yaml:"id"`
	RefSvcId int `yaml:"ref_svc_id"`
}

type yamlServiceRef struct {
	ServiceID  int `yaml:"service_id"`
	Operations []struct {
		Kind   string         `yaml:"kind"`
		Action string         `yaml:"action"`
		Params []yamlParamRef `yaml:"params"`
	} `yaml:"operations"`
}

var (
	serviceParamRefMap     map[int]map[int]int
	serviceParamRefMapOnce sync.Once
	serviceParamRefMapErr  error
)

func getServiceParamRefMap(serviceID int) (map[int]int, error) {
	serviceParamRefMapOnce.Do(loadServiceParamRefMaps)
	if serviceParamRefMapErr != nil {
		return nil, serviceParamRefMapErr
	}
	if serviceParamRefMap == nil {
		return nil, nil
	}
	return serviceParamRefMap[serviceID], nil
}

func loadServiceParamRefMaps() {
	serviceParamRefMap = make(map[int]map[int]int)

	entries, err := fs.ReadDir(resourcesyaml.Files, ".")
	if err != nil {
		serviceParamRefMapErr = err
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) == 0 || name[len(name)-5:] != ".yaml" {
			continue
		}

		b, err := resourcesyaml.Files.ReadFile(name)
		if err != nil {
			serviceParamRefMapErr = err
			return
		}

		var svc yamlServiceRef
		if err := yaml.Unmarshal(b, &svc); err != nil {
			serviceParamRefMapErr = err
			return
		}
		if svc.ServiceID == 0 {
			continue
		}

		idToRef := make(map[int]int)
		for _, op := range svc.Operations {
			if op.Kind != "instance" || op.Action != "create" {
				continue
			}
			for _, p := range op.Params {
				if p.ID == 0 || p.RefSvcId == 0 {
					continue
				}
				idToRef[p.ID] = p.RefSvcId
			}
			break
		}
		serviceParamRefMap[svc.ServiceID] = idToRef
	}
}

func getServiceParamRefCodeMap(serviceID int) (map[string]int, error) {
	codeMap, err := getServiceParamCodeMap(serviceID)
	if err != nil {
		return nil, err
	}
	refMap, err := getServiceParamRefMap(serviceID)
	if err != nil {
		return nil, err
	}
	if len(codeMap) == 0 || len(refMap) == 0 {
		return nil, nil
	}

	codeToRef := make(map[string]int)
	for id, refSvcId := range refMap {
		code, ok := codeMap[id]
		if !ok || strings.TrimSpace(code) == "" {
			continue
		}
		if refSvcId <= 0 {
			continue
		}
		codeToRef[code] = refSvcId
	}
	return codeToRef, nil
}

// ResolveRefSvcParamDisplayNames maps UUIDs in state_params back to display_name values.
func ResolveRefSvcParamDisplayNames(ctx context.Context, client *core.UniversalClient, serviceID int, params map[string]string) (map[string]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if client == nil || len(params) == 0 {
		return params, diags
	}

	refCodeMap, err := getServiceParamRefCodeMap(serviceID)
	if err != nil {
		return params, diags
	}
	if len(refCodeMap) == 0 {
		return params, diags
	}

	resolved := make(map[string]string, len(params))
	for code, value := range params {
		newValue := value
		if refSvcId, ok := refCodeMap[code]; ok && refSvcId > 0 {
			mapped, err := client.ResolveRefSvcParamDisplayName(ctx, refSvcId, value)
			if err == nil && strings.TrimSpace(mapped) != "" {
				newValue = mapped
			}
		}
		resolved[code] = newValue
	}

	return resolved, diags
}
