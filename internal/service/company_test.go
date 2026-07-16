package service

import "testing"

func TestCleanDomain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"full url with https and www", "https://www.paper.id", "paper.id"},
		{"full url with http", "http://example.com", "example.com"},
		{"full url with https", "https://example.com", "example.com"},
		{"domain with trailing slash", "example.com/", "example.com"},
		{"domain with www prefix", "www.example.com", "example.com"},
		{"plain domain", "example.com", "example.com"},
		{"domain with subdomain", "sub.example.com", "sub.example.com"},
		{"domain with whitespace", "  example.com  ", "example.com"},
		{"id domain with www", "www.paper.id", "paper.id"},
		{"complex url", "https://www.sub.example.co.id/path", "sub.example.co.id/path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanDomain(tt.input)
			if result != tt.expected {
				t.Errorf("CleanDomain(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractDomainName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple domain", "paper.id", "paper"},
		{"com domain", "example.com", "example"},
		{"co id domain", "tokopedia.co.id", "tokopedia"},
		{"subdomain", "sub.example.com", "sub"},
		{"www prefix", "www.example.com", "www"},
		{"no tld", "example", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDomainName(tt.input)
			if result != tt.expected {
				t.Errorf("extractDomainName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
