package resources_core

import (
	"io/fs"
	"sync"

	resourcesyaml "github.com/repinoid/terraform-provider-nubes/resources_yaml"

	"gopkg.in/yaml.v3"
)

type yamlRequiredParam struct {
	ID       int  `yaml:"id"`
	Required bool `yaml:"required"`
}

type yamlRequiredService struct {
	ServiceID  int `yaml:"service_id"`
	Operations []struct {
		Kind   string              `yaml:"kind"`
		Action string              `yaml:"action"`
		Params []yamlRequiredParam `yaml:"params"`
	} `yaml:"operations"`
}

var (
	serviceRequiredParamMap     map[int]map[int]bool
	serviceRequiredParamMapOnce sync.Once
	serviceRequiredParamMapErr  error
)

func getServiceRequiredParamIDs(serviceID int) (map[int]bool, error) {
	serviceRequiredParamMapOnce.Do(loadServiceRequiredParamMaps)
	if serviceRequiredParamMapErr != nil {
		return nil, serviceRequiredParamMapErr
	}
	if serviceRequiredParamMap == nil {
		return nil, nil
	}
	return serviceRequiredParamMap[serviceID], nil
}

func loadServiceRequiredParamMaps() {
	serviceRequiredParamMap = make(map[int]map[int]bool)

	entries, err := fs.ReadDir(resourcesyaml.Files, ".")
	if err != nil {
		serviceRequiredParamMapErr = err
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
			serviceRequiredParamMapErr = err
			return
		}

		var svc yamlRequiredService
		if err := yaml.Unmarshal(b, &svc); err != nil {
			serviceRequiredParamMapErr = err
			return
		}
		if svc.ServiceID == 0 {
			continue
		}

		requiredIDs := make(map[int]bool)
		for _, op := range svc.Operations {
			if op.Kind != "instance" || op.Action != "create" {
				continue
			}
			for _, p := range op.Params {
				if p.ID == 0 || !p.Required {
					continue
				}
				requiredIDs[p.ID] = true
			}
			break
		}
		serviceRequiredParamMap[svc.ServiceID] = requiredIDs
	}
}
