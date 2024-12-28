package util

import (
	"fmt"
	"net/url"
	"strings"
)

// GetConfigPath returns scope.key in string format
// e.g. GetConfigPath("global", "client_base_url") -> "global.client_base_url"
func GetConfigPath(scope string, key string) string {
	return fmt.Sprintf("%s.%s", scope, key)
}

// ExtractDomain extracts the domain from a URL.
// Strips the "www." prefix and ports if present.
// e.g. https://www.google.com -> google.com |
// e.g. https://www.google.com/search?q=golang -> google.com |
// e.g. http://localhost:8080 -> localhost |
func ExtractDomain(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	hostname := parsedURL.Hostname()

	hostname = strings.TrimPrefix(hostname, "www.")

	return hostname, nil
}

// ExtractProtocol extracts the protocol from a URL.
// e.g. https://www.google.com -> https |
// e.g. http://localhost:8080 -> http |
func ExtractProtocol(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	protocol := parsedURL.Scheme

	return protocol, nil
}
