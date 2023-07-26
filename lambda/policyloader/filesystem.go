// policyloader/filesystem.go
package policyloader

import (
	"context"
	"os"
)

// FilesystemPolicyLoader loads policies from the filesystem.
type FilesystemPolicyLoader struct{}

// LoadPolicy loads a policy from the filesystem.
func (p *FilesystemPolicyLoader) LoadPolicy(ctx context.Context, key string) (string, error) {
	filename, err := KeyToFilename("policies." + key)
	if err != nil {
		return "", err
	}

	rawBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", &FileNotFoundError{Key: key}
	}

	return string(rawBytes), nil
}
