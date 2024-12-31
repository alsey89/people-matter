package util

import (
	"fmt"
)

// GetConfigPath returns scope.key in string format
// e.g. GetConfigPath("global", "client_base_url") -> "global.client_base_url"
func GetConfigPath(scope string, key string) string {
	return fmt.Sprintf("%s.%s", scope, key)
}
