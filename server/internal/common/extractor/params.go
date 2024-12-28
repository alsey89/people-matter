package extractor

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func ExtractIDFromPathParamAsUINT(c echo.Context, key string) (*uint, error) {

	pathParam := c.Param(key)
	if pathParam == "" {
		return nil, fmt.Errorf("extractIDAsUINTFromPathParam: path param %s is empty", key)
	}

	// Convert string to uint using fmt.Sscanf
	var output uint

	_, err := fmt.Sscanf(pathParam, "%d", &output)
	if err != nil {
		return nil, fmt.Errorf("extractIDAsUINTFromPathParam: failed to parse %s as uint: %w", pathParam, err)
	}

	return &output, nil
}

func ExtractBoolFromQueryParam(c echo.Context, key string, defaultValue bool) (*bool, error) {

	queryParam := c.QueryParam(key)
	if queryParam == "" {
		return &defaultValue, nil
	}

	// Convert string to bool using fmt.Sscanf
	var output bool

	_, err := fmt.Sscanf(queryParam, "%t", &output)
	if err != nil {
		return nil, fmt.Errorf("extractBoolFromQueryParam: failed to parse %s as bool: %w", queryParam, err)
	}

	return &output, nil
}
