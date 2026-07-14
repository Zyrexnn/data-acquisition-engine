package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type LocationData struct {
	DisplayName string                 `json:"display_name"`
	Latitude    string                 `json:"latitude"`
	Longitude   string                 `json:"longitude"`
	Importance  float64                `json:"importance"`
	OSMType     string                 `json:"osm_type"`
	Address     map[string]interface{} `json:"address"`
}

type LocationService struct {
	client *http.Client
}

func NewLocationService() *LocationService {
	return &LocationService{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type nominatimResult struct {
	DisplayName string                 `json:"display_name"`
	Lat         string                 `json:"lat"`
	Lon         string                 `json:"lon"`
	Importance  float64                `json:"importance"`
	OSMType     string                 `json:"osm_type"`
	Address     map[string]interface{} `json:"address"`
}

func (s *LocationService) Find(query string) (*LocationData, error) {
	if query == "" {
		return nil, fmt.Errorf("query is empty")
	}

	encoded := url.QueryEscape(query)
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=jsonv2&addressdetails=1", encoded)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "DataAcquisitionEngine/1.0 (contact@syntaxteras.com)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach Nominatim API: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Nominatim response: %w", err)
	}

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("Nominatim API returned 403 Forbidden")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Nominatim API returned status %d", resp.StatusCode)
	}

	var results []nominatimResult
	if err := json.Unmarshal(bodyBytes, &results); err != nil {
		return nil, fmt.Errorf("failed to parse Nominatim JSON: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no location found for query: %s", query)
	}

	r := results[0]
	return &LocationData{
		DisplayName: r.DisplayName,
		Latitude:    r.Lat,
		Longitude:   r.Lon,
		Importance:  r.Importance,
		OSMType:     r.OSMType,
		Address:     r.Address,
	}, nil
}
