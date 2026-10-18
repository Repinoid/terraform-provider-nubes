package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// FindInstanceByDisplayName finds an instance by display_name for a given serviceId.
// Если найдено больше одного non-deleted инстанса — возвращает ошибку.
func (c *UniversalClient) FindInstanceByDisplayName(ctx context.Context, serviceId int, displayName string) (*InstanceStateResponse, error) {
	debug := func(format string, args ...interface{}) {
		if !isHTTPDebugEnabled() {
			return
		}
		f, err := os.OpenFile(debugLogPath("nubes_find_debug.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return
		}
		defer f.Close()
		fmt.Fprintf(f, format+"\n", args...)
	}
	var found []InstanceStateResponse

	// Быстрый поиск через search API
	searchQuery := url.Values{}
	searchQuery.Set("fields", "instanceUid,displayName,serviceId")
	searchQuery.Set("page", "1")
	searchQuery.Set("pageSize", "100")
	searchQuery.Set("search", displayName)
	searchQuery.Set("isAuxiliary", "false")
	searchQuery.Set("isDeleted", "false")
	if serviceId > 0 {
		searchQuery.Set("serviceId", fmt.Sprintf("%d", serviceId))
	}
	respBody, _, err := c.doRequest(ctx, "GET", "/instances?"+searchQuery.Encode(), nil)
	if err != nil {
		debug("[FIND-DEBUG] search doRequest err: %v (serviceId=%d, name=%q)", err, serviceId, displayName)
	} else {
		var res struct {
			Results []struct {
				InstanceUid string `json:"instanceUid"`
				DisplayName string `json:"displayName"`
				ServiceId   int    `json:"serviceId"`
			} `json:"results"`
		}
		if err := json.Unmarshal(respBody, &res); err != nil {
			debug("[FIND-DEBUG] search json.Unmarshal err: %v (serviceId=%d, name=%q)", err, serviceId, displayName)
		} else {
			debug("[FIND-DEBUG] search returned %d results (serviceId=%d, name=%q)", len(res.Results), serviceId, displayName)
			for _, item := range res.Results {
				matchSvc := item.ServiceId == serviceId
				matchName := strings.EqualFold(item.DisplayName, displayName)
				debug("[FIND-DEBUG] search item uid=%s name=%q svcId=%d matchSvc=%v matchName=%v",
					item.InstanceUid, item.DisplayName, item.ServiceId, matchSvc, matchName)
				if matchSvc && matchName {
					state, err := c.GetInstanceStateRaw(ctx, item.InstanceUid)
					if err != nil {
						debug("[FIND-DEBUG] GetInstanceStateRaw err for uid=%s: %v", item.InstanceUid, err)
						continue
					}
					if state == nil {
						debug("[FIND-DEBUG] GetInstanceStateRaw returned nil for uid=%s", item.InstanceUid)
						continue
					}
					if isInstanceDeleted(state) {
						debug("[FIND-DEBUG] GetInstanceStateRaw returned deleted for uid=%s (status=%q)", item.InstanceUid, state.ExplainedStatus)
						continue
					}
					found = append(found, *state)
				}
			}
		}
	}

	// Если быстрый поиск не дал результатов — пагинированный fallback
	if len(found) == 0 {
		debug("[FIND-DEBUG] search path failed, entering fallback (serviceId=%d, name=%q)", serviceId, displayName)
		const maxFallbackPages = 100
		hitPageLimit := false
		page := 1
		for {
			reqURL := c.buildURL(fmt.Sprintf("/instances?page=%d&size=100", page))
			req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("User-Agent", userAgent)
			if c.ApiToken != "" {
				req.Header.Set("Authorization", "Bearer "+c.ApiToken)
			}

			resp, err := c.HttpClient.Do(req)
			if err != nil {
				debug("[FIND-DEBUG] fallback page=%d HttpClient.Do err: %v", page, err)
				return nil, err
			}

			// Читаем тело сразу (не defer в цикле — иначе копятся открытые тела),
			// чтобы при не-200 отдать юзеру читаемый текст, а не голый статус.
			respBody, readErr := io.ReadAll(resp.Body)
			resp.Body.Close()
			if readErr != nil {
				debug("[FIND-DEBUG] fallback page=%d read err: %v", page, readErr)
				return nil, readErr
			}

			if resp.StatusCode != 200 {
				debug("[FIND-DEBUG] fallback page=%d HTTP %d", page, resp.StatusCode)
				return nil, formatAPIError(resp.StatusCode, respBody)
			}

			var res struct {
				Results []struct {
					InstanceUid string `json:"instanceUid"`
					DisplayName string `json:"displayName"`
					ServiceId   int    `json:"serviceId"`
				} `json:"results"`
			}
			if err := json.Unmarshal(respBody, &res); err != nil {
				debug("[FIND-DEBUG] fallback page=%d json decode err: %v", page, err)
				return nil, err
			}

			if len(res.Results) == 0 {
				break
			}

			debug("[FIND-DEBUG] fallback page=%d has %d results", page, len(res.Results))
			for _, item := range res.Results {
				matchSvc := item.ServiceId == serviceId
				matchName := strings.EqualFold(item.DisplayName, displayName)
				if matchSvc || matchName {
					debug("[FIND-DEBUG] fallback item uid=%s name=%q svcId=%d matchSvc=%v matchName=%v",
						item.InstanceUid, item.DisplayName, item.ServiceId, matchSvc, matchName)
				}
				if matchSvc && matchName {
					state, err := c.GetInstanceStateRaw(ctx, item.InstanceUid)
					if err != nil {
						debug("[FIND-DEBUG] fallback GetInstanceStateRaw err for uid=%s: %v", item.InstanceUid, err)
						continue
					}
					if state == nil {
						debug("[FIND-DEBUG] fallback GetInstanceStateRaw returned nil for uid=%s", item.InstanceUid)
						continue
					}
					if isInstanceDeleted(state) {
						debug("[FIND-DEBUG] fallback GetInstanceStateRaw returned deleted for uid=%s (status=%q)", item.InstanceUid, state.ExplainedStatus)
						continue
					}
					found = append(found, *state)
				}
			}

			page++
			if page > maxFallbackPages {
				hitPageLimit = true
				break
			}
		}

		if hitPageLimit && len(found) == 0 {
			return nil, fmt.Errorf("поиск инстанса по display_name достиг лимита пагинации (%d страниц); сузьте фильтр или повторите с более точными параметрами", maxFallbackPages)
		}
	}

	if len(found) == 0 {
		return nil, nil
	}
	if len(found) == 1 {
		return &found[0], nil
	}

	// Множественные non-deleted инстансы — ошибка
	details := make([]string, 0, len(found))
	for _, inst := range found {
		details = append(details, fmt.Sprintf("  - %s (статус: %s)", inst.InstanceUid, inst.ExplainedStatus))
	}
	return nil, fmt.Errorf(
		"обнаружено %d инстансов с именем '%s' (serviceId=%d):\n%s\nНевозможно определить какой adopt-ить. Удалите лишние через ЛК (Управление облаком) или укажите конкретный UUID через 'terraform import'",
		len(found), displayName, serviceId, strings.Join(details, "\n"),
	)
}
