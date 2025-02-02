package util

import (
	"fmt"
	"net/url"
	"strings"
)

// Returns scope.key in string format
func GetConfigPath(scope string, key string) string {
	return fmt.Sprintf("%s.%s", scope, key)
}

func PathToFullURL(path string, subdomain string, domain string) (*string, error) {
	if path == "" || subdomain == "" || domain == "" {
		return nil, fmt.Errorf("PathToFullURL: missing one or more required parameters. path: %s, subdomain: %s, domain: %s", path, subdomain, domain)
	}

	// if path is not a valid relative path, return an error
	if path[0] != '/' {
		return nil, fmt.Errorf("PathToFullURL: path expects a relative path. path: %s", path)
	}

	// localhost is http
	if strings.Contains(domain, "localhost") {
		fullURL := fmt.Sprintf("http://%s.%s%s", subdomain, domain, path)
		return &fullURL, nil
	}

	fullURL := fmt.Sprintf("https://%s.%s%s", subdomain, domain, path)
	return &fullURL, nil
}

// Encodes a string to be used as a URL parameter
func EncodeQueryParam(param string) string {
	return url.QueryEscape(param)
}

// Ecodes a string to be used as a path parameter
func EncodePathParam(param string) string {
	return url.PathEscape(param)
}
