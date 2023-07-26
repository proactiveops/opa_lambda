package policyloader

import (
	"context"
	"os"
)

// PolicyLoader loads policies.
type PolicyLoader interface {
	LoadPolicy(ctx context.Context, key string) (string, error)
}

// NewPolicyLoader creates a new PolicyLoader.
func NewPolicyLoader(ctx context.Context) (PolicyLoader, error) {
	var loader PolicyLoader
	var err error

	if bucketName := os.Getenv("S3_BUCKET"); bucketName != "" {
		loader, err = NewS3PolicyLoader(bucketName)
	} else {
		loader = &FilesystemPolicyLoader{}
	}

	return loader, err
}
