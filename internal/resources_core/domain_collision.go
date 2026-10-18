package resources_core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/repinoid/terraform-provider-nubes/internal/core"
)

type domainOwnerInstance struct {
	InstanceUID string
	DisplayName string
	ServiceID   int
}

type domainOwnersSnapshot struct {
	ownersByDomain map[string]domainOwnerInstance

	expiresAt time.Time
}

var domainOwnersCache = struct {
	mu    sync.Mutex
	items map[string]domainOwnersSnapshot
}{
	items: make(map[string]domainOwnersSnapshot),
}

func findDomainOwnerInstance(ctx context.Context, client *core.UniversalClient, domain string, allowedServiceIDs []int) (*domainOwnerInstance, error) {
	targetDomain := strings.TrimSpace(domain)
	if client == nil || targetDomain == "" {
		return nil, nil
	}
	allowedSet := normalizeServiceIDSet(allowedServiceIDs)
	cacheKey := buildDomainCacheKey(client, allowedSet)

	owners, err := getOrBuildDomainOwnersSnapshot(ctx, client, cacheKey, allowedSet)
	if err != nil {
		return nil, err
	}

	owner, ok := owners[strings.ToLower(targetDomain)]
	if !ok {
		return nil, nil
	}
	ownerCopy := owner
	return &ownerCopy, nil
}

func getOrBuildDomainOwnersSnapshot(ctx context.Context, client *core.UniversalClient, cacheKey string, allowedSet map[int]struct{}) (map[string]domainOwnerInstance, error) {
	now := time.Now()

	domainOwnersCache.mu.Lock()
	if cached, ok := domainOwnersCache.items[cacheKey]; ok && now.Before(cached.expiresAt) {
		ownersCopy := make(map[string]domainOwnerInstance, len(cached.ownersByDomain))
		for domainKey, owner := range cached.ownersByDomain {
			ownersCopy[domainKey] = owner
		}
		domainOwnersCache.mu.Unlock()
		return ownersCopy, nil
	}
	domainOwnersCache.mu.Unlock()

	built, err := collectDomainOwners(ctx, client, allowedSet)
	if err != nil {
		return nil, err
	}

	domainOwnersCache.mu.Lock()
	domainOwnersCache.items[cacheKey] = domainOwnersSnapshot{
		ownersByDomain: built,
		expiresAt:      now.Add(30 * time.Second),
	}
	domainOwnersCache.mu.Unlock()

	ownersCopy := make(map[string]domainOwnerInstance, len(built))
	for domainKey, owner := range built {
		ownersCopy[domainKey] = owner
	}
	return ownersCopy, nil
}

func collectDomainOwners(ctx context.Context, client *core.UniversalClient, allowedSet map[int]struct{}) (map[string]domainOwnerInstance, error) {
	owners := map[string]domainOwnerInstance{}

	page := 1
	for {
		url := fmt.Sprintf("%s/instances?page=%d&size=100", client.ApiEndpoint, page)
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0")
		if client.ApiToken != "" {
			req.Header.Set("Authorization", "Bearer "+client.ApiToken)
		}

		resp, err := client.HttpClient.Do(req)
		if err != nil {
			return nil, err
		}

		var parsed struct {
			Results []struct {
				InstanceUID string `json:"instanceUid"`
				DisplayName string `json:"displayName"`
				ServiceID   int    `json:"serviceId"`
			} `json:"results"`
		}

		decodeErr := json.NewDecoder(resp.Body).Decode(&parsed)
		_ = resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("HTTP статус %d", resp.StatusCode)
		}
		if decodeErr != nil {
			return nil, decodeErr
		}
		if len(parsed.Results) == 0 {
			break
		}

		for _, item := range parsed.Results {
			if len(allowedSet) > 0 {
				if _, ok := allowedSet[item.ServiceID]; !ok {
					continue
				}
			}

			params, err := client.GetInstanceStateParams(ctx, item.InstanceUID)
			if err != nil {
				continue
			}
			for key, value := range params {
				if !strings.EqualFold(strings.TrimSpace(key), "domain") {
					continue
				}
				normalizedDomain := strings.ToLower(strings.TrimSpace(value))
				if normalizedDomain == "" {
					continue
				}
				if _, exists := owners[normalizedDomain]; !exists {
					owners[normalizedDomain] = domainOwnerInstance{
						InstanceUID: item.InstanceUID,
						DisplayName: item.DisplayName,
						ServiceID:   item.ServiceID,
					}
				}
			}
		}

		page++
		if page > 100 {
			break
		}
	}

	return owners, nil
}

func normalizeServiceIDSet(serviceIDs []int) map[int]struct{} {
	set := make(map[int]struct{})
	for _, serviceID := range serviceIDs {
		if serviceID <= 0 {
			continue
		}
		set[serviceID] = struct{}{}
	}
	return set
}

func buildDomainCacheKey(client *core.UniversalClient, allowedSet map[int]struct{}) string {
	ids := make([]int, 0, len(allowedSet))
	for serviceID := range allowedSet {
		ids = append(ids, serviceID)
	}
	sort.Ints(ids)
	parts := make([]string, 0, len(ids))
	for _, serviceID := range ids {
		parts = append(parts, fmt.Sprintf("%d", serviceID))
	}
	return client.ApiEndpoint + "|" + strings.Join(parts, ",")
}
