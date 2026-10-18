package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ===== Поиск инстансов для refSvc (по displayName/UUID, с учётом сервиса) =====

func (c *UniversalClient) findInstanceUidByDisplayNameRefSvc(ctx context.Context, serviceId int, displayName string) (string, error) {
	// Собираем ВСЕ совпадения, затем выбираем лучшее (running > suspended > остальные).
	// Deleted инстансы пропускаются — они не должны участвовать в resolve.
	type candidate struct {
		uid    string
		status string
	}
	var candidates []candidate

	page := 1
	for {
		path := fmt.Sprintf("/instances?page=%d&size=100", page)
		respBody, _, err := c.doRequest(ctx, "GET", path, nil)
		if err != nil {
			return "", err
		}

		var res struct {
			Results []struct {
				InstanceUid     string `json:"instanceUid"`
				DisplayName     string `json:"displayName"`
				ServiceId       int    `json:"serviceId"`
				IsDeleted       bool   `json:"isDeleted"`
				ExplainedStatus string `json:"explainedStatus"`
			} `json:"results"`
		}
		if err := json.Unmarshal(respBody, &res); err != nil {
			return "", err
		}

		if len(res.Results) == 0 {
			break
		}

		for _, item := range res.Results {
			if item.ServiceId != serviceId || !strings.EqualFold(item.DisplayName, displayName) {
				continue
			}
			// Пропускаем deleted инстансы
			if item.IsDeleted {
				continue
			}
			status := strings.ToLower(strings.TrimSpace(item.ExplainedStatus))
			if status == "deleted" {
				continue
			}
			candidates = append(candidates, candidate{uid: item.InstanceUid, status: status})
		}

		page++
		if page > 100 {
			break
		}
	}

	if len(candidates) == 0 {
		return "", nil
	}
	if len(candidates) == 1 {
		return candidates[0].uid, nil
	}
	// Несколько non-deleted совпадений — предпочитаем running
	for _, c := range candidates {
		if strings.Contains(c.status, "running") || strings.Contains(c.status, "active") {
			return c.uid, nil
		}
	}
	// Нет running — предпочитаем suspended
	for _, c := range candidates {
		if strings.Contains(c.status, "suspend") {
			return c.uid, nil
		}
	}
	// Вернуть первый (лучше чем ничего)
	return candidates[0].uid, nil
}

// findInstanceDisplayNameByUidRefSvc resolves display_name by instance UID and service ID.
func (c *UniversalClient) findInstanceDisplayNameByUidRefSvc(ctx context.Context, serviceId int, instanceUid string) (string, error) {
	page := 1
	for {
		path := fmt.Sprintf("/instances?page=%d&size=100", page)
		respBody, _, err := c.doRequest(ctx, "GET", path, nil)
		if err != nil {
			return "", err
		}

		var res struct {
			Results []struct {
				InstanceUid string `json:"instanceUid"`
				DisplayName string `json:"displayName"`
				ServiceId   int    `json:"serviceId"`
			} `json:"results"`
		}
		if err := json.Unmarshal(respBody, &res); err != nil {
			return "", err
		}

		if len(res.Results) == 0 {
			break
		}

		for _, item := range res.Results {
			if item.ServiceId == serviceId && strings.EqualFold(item.InstanceUid, instanceUid) {
				return item.DisplayName, nil
			}
		}

		page++
		if page > 100 {
			break
		}
	}

	return "", nil
}

// getInstanceDisplayNameByUidRefSvc fetches display_name by instance UID and checks service ID.
func (c *UniversalClient) getInstanceDisplayNameByUidRefSvc(ctx context.Context, serviceId int, instanceUid string) (string, error) {
	if strings.TrimSpace(instanceUid) == "" {
		return "", nil
	}

	respBody, _, err := c.doRequest(ctx, "GET", fmt.Sprintf("/instances/%s", instanceUid), nil)
	if err != nil {
		return "", err
	}

	var res struct {
		Instance struct {
			InstanceUid string `json:"instanceUid"`
			DisplayName string `json:"displayName"`
			ServiceId   int    `json:"serviceId"`
		} `json:"instance"`
	}
	if err := json.Unmarshal(respBody, &res); err != nil {
		return "", err
	}
	if strings.TrimSpace(res.Instance.InstanceUid) == "" {
		return "", nil
	}
	if serviceId > 0 && res.Instance.ServiceId != serviceId {
		return "", nil
	}

	return res.Instance.DisplayName, nil
}
