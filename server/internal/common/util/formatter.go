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

// Converts a string to a uint
func StringToUINT(s string) (*uint, error) {
	var output uint
	_, err := fmt.Sscanf(s, "%d", &output)
	if err != nil {
		return nil, fmt.Errorf("StringToUint: %s", err)
	}
	return &output, nil
}

// Safely converts a float64 to uint, ensuring it's a non-negative whole number.
func Float64ToUINT(floatNum float64) (*uint, error) {
	var positiveWholeUINT uint

	if floatNum < 0 {
		return nil, fmt.Errorf("Float64ToUINT: %f is negative", floatNum)
	}

	//casting float to int will truncate the decimal part
	intNum := int(floatNum)

	//if they are not equal, then the float is not a whole number
	if floatNum != float64(intNum) {
		return nil, fmt.Errorf("Float64ToUINT: %f is not a whole number", floatNum)
	}

	positiveWholeUINT = uint(intNum)

	return &positiveWholeUINT, nil
}

// InsertSubdomainToRawURL inserts a subdomain to a raw URL.
// Example: InsertSubdomainToRawURL("https://example.com", "subdomain") -> "https://subdomain.example.com"
// Example: InsertSubdomainToRawURL("https://example.com:8080", "subdomain") -> "https://subdomain.example.com:8080"
func InsertSubdomainToRawURL(rawURL, subdomain string) (*string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	scheme := parsedURL.Scheme
	hostname := parsedURL.Hostname()
	port := parsedURL.Port()
	path := parsedURL.Path
	query := parsedURL.RawQuery

	if subdomain != "" {
		hostname = subdomain + "." + hostname
	}
	if port != "" {
		hostname = hostname + ":" + port
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, hostname)

	// Append path and query if present
	if path != "" {
		baseURL += path
	}
	if query != "" {
		baseURL += "?" + query
	}

	return &baseURL, nil
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
