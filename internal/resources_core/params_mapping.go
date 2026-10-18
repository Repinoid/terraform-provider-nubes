package resources_core

import (
	"io/fs"
	"sync"

	resourcesyaml "github.com/repinoid/terraform-provider-nubes/resources_yaml"

	"gopkg.in/yaml.v3"
)

type yamlParam struct {
	ID   int    `yaml:"id"`
	Code string `yaml:"code"`
}

type yamlService struct {
	ServiceID  int `yaml:"service_id"`
	Operations []struct {
		Kind   string      `yaml:"kind"`
		Action string      `yaml:"action"`
		Params []yamlParam `yaml:"params"`
	} `yaml:"operations"`
}

var (
	serviceParamMap     map[int]map[int]string
	serviceParamMapOnce sync.Once
	serviceParamMapErr  error
)

func getServiceParamCodeMap(serviceID int) (map[int]string, error) {
	serviceParamMapOnce.Do(loadServiceParamMaps)
	if serviceParamMapErr != nil {
		return nil, serviceParamMapErr
	}
	if serviceParamMap == nil {
		return nil, nil
	}
	return serviceParamMap[serviceID], nil
}

func loadServiceParamMaps() {
	serviceParamMap = make(map[int]map[int]string)

	entries, err := fs.ReadDir(resourcesyaml.Files, ".")
	if err != nil {
		serviceParamMapErr = err
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
			serviceParamMapErr = err
			return
		}

		var svc yamlService
		if err := yaml.Unmarshal(b, &svc); err != nil {
			serviceParamMapErr = err
			return
		}
		if svc.ServiceID == 0 {
			continue
		}

		idToCode := make(map[int]string)
		for _, op := range svc.Operations {
			if op.Kind != "instance" || op.Action != "create" {
				continue
			}
			for _, p := range op.Params {
				if p.ID == 0 || p.Code == "" {
					continue
				}
				idToCode[p.ID] = p.Code
			}
			break
		}
		serviceParamMap[svc.ServiceID] = idToCode
	}
}
