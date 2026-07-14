package service

type DomainData struct {
	Domain       string   `json:"domain"`
	Registrar    string   `json:"registrar"`
	RegisteredAt string   `json:"registered_at"`
	ExpiredAt    string   `json:"expired_at"`
	LastUpdated  string   `json:"last_updated"`
	Status       []string `json:"status"`
	Nameservers  []string `json:"nameservers"`
}

type DomainService struct{}

func NewDomainService() *DomainService {
	return &DomainService{}
}

func (s *DomainService) Extract(domain string) (*DomainData, error) {
	// TODO: implement RDAP API call and response mapping
	return &DomainData{
		Domain:       domain,
		Registrar:    "",
		RegisteredAt: "",
		ExpiredAt:    "",
		LastUpdated:  "",
		Status:       []string{},
		Nameservers:  []string{},
	}, nil
}
