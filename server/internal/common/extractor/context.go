package extractor

import (
	"fmt"

	"github.com/alsey89/people-matter/internal/common/API"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func ExtractTenantIdentifierFromContext(c echo.Context) (*string, error) {
	tenantIdentifier, ok := c.Get(API.ContextTenantID).(string)
	if !ok {
		return nil, fmt.Errorf("ExtractTenantIdentifierFromContext: %s", "error extracting tenant identifier from context")
	}
	if tenantIdentifier == "" {
		return nil, fmt.Errorf("ExtractTenantIdentifierFromContext: %s", "no tenant identifier in context")
	}

	return &tenantIdentifier, nil
}

func ExtractTokenAndClaimsFromContext(c echo.Context) (*jwt.Token, jwt.MapClaims, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, nil, fmt.Errorf("ExtractTokenFromContext: %s", "error extracting token from context")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("ExtractTokenFromContext: %s", "error extracting claims from token")
	}

	return user, claims, nil
}
