// policyloader/util.go
package policyloader

import (
	"strings"
)

// KeyToFilename converts a policy key name to a filename.
func KeyToFilename(key string) (string, error) {
	if strings.Contains(key, "/") {
		return "", &InvalidKeyNameError{Key: key}
	}

	filename := strings.ReplaceAll(key, ".", "/")
	return filename + ".rego", nil
}
