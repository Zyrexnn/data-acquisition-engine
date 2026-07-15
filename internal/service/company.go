package service

import (
	"fmt"
	"strings"
	"sync"
)

type CompanyData struct {
	Website  *WebsiteData  `json:"website"`
	Domain   *DomainData   `json:"domain"`
	Location *LocationData `json:"location"`
}

type CompanyService struct {
	websiteSvc  *WebsiteService
	domainSvc   *DomainService
	locationSvc *LocationService
}

func NewCompanyService() *CompanyService {
	return &CompanyService{
		websiteSvc:  NewWebsiteService(),
		domainSvc:   NewDomainService(),
		locationSvc: NewLocationService(),
	}
}

func CleanDomain(input string) string {
	input = strings.TrimSpace(input)
	input = strings.TrimPrefix(input, "https://")
	input = strings.TrimPrefix(input, "http://")
	input = strings.TrimPrefix(input, "www.")
	input = strings.TrimRight(input, "/")
	return input
}

func extractDomainName(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

func (s *CompanyService) GetInfo(domain string) (*CompanyData, error) {
	var (
		websiteWg sync.WaitGroup
		locationWg sync.WaitGroup
		mu         sync.Mutex
		website    *WebsiteData
		domainD    *DomainData
		location   *LocationData
	)

	websiteWg.Add(2)

	go func() {
		defer websiteWg.Done()
		data, err := s.websiteSvc.Extract("https://" + domain)
		if err == nil {
			mu.Lock()
			website = data
			mu.Unlock()
		}
	}()

	go func() {
		defer websiteWg.Done()
		data, err := s.domainSvc.Extract(domain)
		if err == nil {
			mu.Lock()
			domainD = data
			mu.Unlock()
		}
	}()

	websiteWg.Wait()

	if website != nil && website.Title != "" {
		locationWg.Add(1)
		go func() {
			defer locationWg.Done()
			data, err := s.locationSvc.Find(website.Title)
			if err == nil && data != nil {
				mu.Lock()
				location = data
				mu.Unlock()
			}
		}()
		locationWg.Wait()
	}

	if location == nil {
		domainName := extractDomainName(domain)
		if domainName != "" {
			data, err := s.locationSvc.Find(domainName)
			if err == nil {
				location = data
			}
		}
	}

	if website == nil && domainD == nil && location == nil {
		return nil, fmt.Errorf("all services failed to return data")
	}

	return &CompanyData{
		Website:  website,
		Domain:   domainD,
		Location: location,
	}, nil
}
