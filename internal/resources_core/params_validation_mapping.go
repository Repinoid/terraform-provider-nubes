package resources_core

import (
	"io/fs"
	"sync"

	resourcesyaml "github.com/repinoid/terraform-provider-nubes/resources_yaml"

	"gopkg.in/yaml.v3"
)

type yamlParamValidation struct {
	ID        int      `yaml:"id"`
	Code      string   `yaml:"code"`
	RefSvcId  int      `yaml:"ref_svc_id"`
	Func      string   `yaml:"func"`
	ValueList []string `yaml:"value_list"`
}

type yamlServiceValidation struct {
	ServiceID  int `yaml:"service_id"`
	Operations []struct {
		Kind   string                `yaml:"kind"`
		Action string                `yaml:"action"`
		Params []yamlParamValidation `yaml:"params"`
	} `yaml:"operations"`
}

type paramValidationSpec struct {
	Code      string
	RefSvcId  int
	Func      string
	ValueList []string
}

var (
	serviceParamValidationMap     map[int]map[int]paramValidationSpec
	serviceParamValidationMapOnce sync.Once
	serviceParamValidationMapErr  error
)

func getServiceParamValidationSpecMap(serviceID int) (map[int]paramValidationSpec, error) {
	serviceParamValidationMapOnce.Do(loadServiceParamValidationMaps)
	if serviceParamValidationMapErr != nil {
		return nil, serviceParamValidationMapErr
	}
	if serviceParamValidationMap == nil {
		return nil, nil
	}
	return serviceParamValidationMap[serviceID], nil
}

func loadServiceParamValidationMaps() {
	serviceParamValidationMap = make(map[int]map[int]paramValidationSpec)

	entries, err := fs.ReadDir(resourcesyaml.Files, ".")
	if err != nil {
		serviceParamValidationMapErr = err
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
			serviceParamValidationMapErr = err
			return
		}

		var svc yamlServiceValidation
		if err := yaml.Unmarshal(b, &svc); err != nil {
			serviceParamValidationMapErr = err
			return
		}
		if svc.ServiceID == 0 {
			continue
		}

		specByID := make(map[int]paramValidationSpec)
		for _, op := range svc.Operations {
			if op.Kind != "instance" || op.Action != "create" {
				continue
			}
			for _, p := range op.Params {
				if p.ID == 0 {
					continue
				}
				cp := make([]string, 0, len(p.ValueList))
				cp = append(cp, p.ValueList...)
				specByID[p.ID] = paramValidationSpec{
					Code:      p.Code,
					RefSvcId:  p.RefSvcId,
					Func:      p.Func,
					ValueList: cp,
				}
			}
			break
		}

		serviceParamValidationMap[svc.ServiceID] = specByID
	}
}
