package extractor

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// Extract a path a parameter and return it as a specific type
func ExtractFromPathParamAs(c echo.Context, key string, output interface{}) (interface{}, error) {

	pathParam := c.Param(key)
	if pathParam == "" {
		return nil, fmt.Errorf("extractFromPathParamAs: path param %s is empty", key)
	}

	// Convert string to specific type using fmt.Sscanf
	_, err := fmt.Sscanf(pathParam, "%v", output)
	if err != nil {
		return nil, fmt.Errorf("extractFromPathParamAs: failed to parse %s as %T: %w", pathParam, output, err)
	}

	return output, nil
}

// Extract a query parameter and return it as a specific type
func ExtractFromQueryParamAs(c echo.Context, key string, output interface{}) (interface{}, error) {

	queryParam := c.QueryParam(key)
	if queryParam == "" {
		return nil, fmt.Errorf("extractFromQueryParamAs: query param %s is empty", key)
	}

	// Convert string to specific type using fmt.Sscanf
	_, err := fmt.Sscanf(queryParam, "%v", output)
	if err != nil {
		return nil, fmt.Errorf("extractFromQueryParamAs: failed to parse %s as %T: %w", queryParam, output, err)
	}

	return output, nil
}
