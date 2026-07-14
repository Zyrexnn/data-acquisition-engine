package service

import "sync"

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

func (s *CompanyService) GetInfo(domain string) *CompanyData {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		website  *WebsiteData
		domainD  *DomainData
		location *LocationData
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		data, err := s.websiteSvc.Extract("https://" + domain)
		if err == nil {
			mu.Lock()
			website = data
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		data, err := s.domainSvc.Extract(domain)
		if err == nil {
			mu.Lock()
			domainD = data
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		data, err := s.locationSvc.Find(domain)
		if err == nil {
			mu.Lock()
			location = data
			mu.Unlock()
		}
	}()

	wg.Wait()

	return &CompanyData{
		Website:  website,
		Domain:   domainD,
		Location: location,
	}
}
