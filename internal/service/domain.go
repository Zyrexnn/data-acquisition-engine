package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type DomainData struct {
	Domain       string   `json:"domain"`
	Registrar    string   `json:"registrar"`
	RegisteredAt string   `json:"registered_at"`
	ExpiredAt    string   `json:"expired_at"`
	LastUpdated  string   `json:"last_updated"`
	Status       []string `json:"status"`
	Nameservers  []string `json:"nameservers"`
}

type DomainService struct {
	client *http.Client
}

func NewDomainService() *DomainService {
	return &DomainService{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type rdapResponse struct {
	Entities []struct {
		Roles      []string `json:"roles"`
		VcardArray []interface{} `json:"vcardArray"`
	} `json:"entities"`
	Events []struct {
		EventAction string `json:"eventAction"`
		EventDate   string `json:"eventDate"`
	} `json:"events"`
	Status []struct {
		Status string `json:"status"`
	} `json:"status"`
	Nameservers []struct {
		LDHName string `json:"ldhName"`
	} `json:"nameservers"`
}

type rdapStatus []string

func (s *rdapStatus) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	switch v := raw.(type) {
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		*s = result
	}
	return nil
}

type rdapResponseAlt struct {
	Entities []struct {
		Roles      []string `json:"roles"`
		VcardArray []interface{} `json:"vcardArray"`
	} `json:"entities"`
	Events []struct {
		EventAction string `json:"eventAction"`
		EventDate   string `json:"eventDate"`
	} `json:"events"`
	Status      rdapStatus `json:"status"`
	Nameservers []struct {
		LDHName string `json:"ldhName"`
	} `json:"nameservers"`
}

func (s *DomainService) Extract(domain string) (*DomainData, error) {
	cleaned := CleanDomain(domain)
	if cleaned == "" {
		return nil, fmt.Errorf("domain is empty after cleaning")
	}

	url := fmt.Sprintf("https://rdap.org/domain/%s", cleaned)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to reach RDAP API: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read RDAP response: %w", err)
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("domain not found: %s", cleaned)
	}
	if resp.StatusCode == 429 {
		return nil, fmt.Errorf("RDAP API rate limited")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("RDAP API returned status %d", resp.StatusCode)
	}

	var rdap rdapResponseAlt
	if err := json.Unmarshal(bodyBytes, &rdap); err != nil {
		return nil, fmt.Errorf("failed to parse RDAP JSON: %w", err)
	}

	data := &DomainData{
		Domain:      cleaned,
		Status:      []string(rdap.Status),
		Nameservers: []string{},
	}

	for _, entity := range rdap.Entities {
		for _, role := range entity.Roles {
			if role == "registrar" {
				data.Registrar = extractVcardFn(entity.VcardArray)
				break
			}
		}
	}

	for _, event := range rdap.Events {
		t, err := time.Parse(time.RFC3339, event.EventDate)
		if err != nil {
			continue
		}
		formatted := t.Format("2006-01-02 15:04:05")
		switch event.EventAction {
		case "registration":
			data.RegisteredAt = formatted
		case "expiration":
			data.ExpiredAt = formatted
		case "last changed":
			data.LastUpdated = formatted
		}
	}

	for _, ns := range rdap.Nameservers {
		if ns.LDHName != "" {
			data.Nameservers = append(data.Nameservers, strings.ToLower(ns.LDHName))
		}
	}

	return data, nil
}

func extractVcardFn(vcard []interface{}) string {
	if len(vcard) < 2 {
		return ""
	}
	arr, ok := vcard[1].([]interface{})
	if !ok {
		return ""
	}
	for _, entry := range arr {
		parts, ok := entry.([]interface{})
		if !ok || len(parts) < 4 {
			continue
		}
		if parts[0] == "fn" {
			if fn, ok := parts[3].(string); ok {
				return fn
			}
		}
	}
	return ""
}
