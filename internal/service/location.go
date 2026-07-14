package service

type LocationData struct {
	DisplayName string                 `json:"display_name"`
	Latitude    string                 `json:"latitude"`
	Longitude   string                 `json:"longitude"`
	Importance  float64                `json:"importance"`
	OSMType     string                 `json:"osm_type"`
	Address     map[string]interface{} `json:"address"`
}

type LocationService struct{}

func NewLocationService() *LocationService {
	return &LocationService{}
}

func (s *LocationService) Find(query string) (*LocationData, error) {
	// TODO: implement Nominatim API call with User-Agent header
	return &LocationData{
		DisplayName: "",
		Latitude:    "",
		Longitude:   "",
		Importance:  0.0,
		OSMType:     "",
		Address:     map[string]interface{}{},
	}, nil
}
