package service

type WebsiteData struct {
	URL        string            `json:"url"`
	Title      string            `json:"title"`
	Description string           `json:"description"`
	Canonical  string            `json:"canonical"`
	Favicon    string            `json:"favicon"`
	Emails     []string          `json:"emails"`
	Phones     []string          `json:"phones"`
	SocialMedia []string         `json:"social_media"`
	OpenGraph  map[string]string `json:"open_graph"`
}

type WebsiteService struct{}

func NewWebsiteService() *WebsiteService {
	return &WebsiteService{}
}

func (s *WebsiteService) Extract(url string) (*WebsiteData, error) {
	// TODO: implement HTML parsing, meta extraction, regex for emails/phones/social
	return &WebsiteData{
		URL:         url,
		Title:       "",
		Description: "",
		Canonical:   "",
		Favicon:     "",
		Emails:      []string{},
		Phones:      []string{},
		SocialMedia: []string{},
		OpenGraph:   map[string]string{},
	}, nil
}
