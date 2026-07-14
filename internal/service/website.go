package service

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type WebsiteData struct {
	URL         string            `json:"url"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Canonical   string            `json:"canonical"`
	Favicon     string            `json:"favicon"`
	Emails      []string          `json:"emails"`
	Phones      []string          `json:"phones"`
	SocialMedia []string          `json:"social_media"`
	OpenGraph   map[string]string `json:"open_graph"`
}

type WebsiteService struct {
	client *http.Client
}

func NewWebsiteService() *WebsiteService {
	return &WebsiteService{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *WebsiteService) Extract(url string) (*WebsiteData, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("target returned status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	rawHTML := string(bodyBytes)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	data := &WebsiteData{
		URL:       url,
		OpenGraph: make(map[string]string),
	}

	data.Title = strings.TrimSpace(doc.Find("title").First().Text())

	doc.Find("meta").Each(func(i int, sel *goquery.Selection) {
		if name, _ := sel.Attr("name"); strings.EqualFold(name, "description") {
			data.Description, _ = sel.Attr("content")
		}
	})

	if href, exists := doc.Find("link[rel='canonical']").First().Attr("href"); exists {
		data.Canonical = href
	}

	faviconSelectors := []string{
		"link[rel='shortcut icon']",
		"link[rel='icon']",
		"link[rel='apple-touch-icon']",
	}
	for _, sel := range faviconSelectors {
		if href, exists := doc.Find(sel).First().Attr("href"); exists {
			data.Favicon = href
			break
		}
	}

	doc.Find("meta").Each(func(i int, sel *goquery.Selection) {
		property, _ := sel.Attr("property")
		content, _ := sel.Attr("content")
		switch property {
		case "og:title":
			data.OpenGraph["title"] = content
		case "og:description":
			data.OpenGraph["description"] = content
		case "og:image":
			data.OpenGraph["image"] = content
		}
	})

	data.Emails = extractEmails(rawHTML)
	data.Phones = extractPhones(rawHTML)
	data.SocialMedia = extractSocialMedia(rawHTML)

	return data, nil
}

func extractEmails(html string) []string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	matches := re.FindAllString(html, -1)
	return uniqueStrings(matches)
}

func extractPhones(html string) []string {
	re := regexp.MustCompile(`(\+62[\s\-]?[\d\s\-]{8,15}|0[8][\d\s\-]{8,13})`)
	matches := re.FindAllString(html, -1)
	cleaned := make([]string, 0, len(matches))
	for _, m := range matches {
		phone := strings.NewReplacer(" ", "", "-", "").Replace(m)
		cleaned = append(cleaned, phone)
	}
	return uniqueStrings(cleaned)
}

func extractSocialMedia(html string) []string {
	re := regexp.MustCompile(`https?://(?:www\.)?(?:linkedin\.com/in|instagram\.com|twitter\.com|x\.com|facebook\.com|fb\.com)/[a-zA-Z0-9._\-]+`)
	matches := re.FindAllString(html, -1)
	cleaned := make([]string, 0, len(matches))
	for _, m := range matches {
		m = strings.TrimRight(m, "/")
		cleaned = append(cleaned, m)
	}
	return uniqueStrings(cleaned)
}

func uniqueStrings(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	result := make([]string, 0, len(input))
	for _, s := range input {
		if _, exists := seen[s]; !exists {
			seen[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
