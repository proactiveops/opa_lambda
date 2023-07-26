package policyloader_test

import (
	"context"
	"opa_lambda/policyloader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPolicyLoader_Filesystem(t *testing.T) {
	loader, err := policyloader.NewPolicyLoader(context.TODO())
	assert.NoError(t, err)
	assert.IsType(t, &policyloader.FilesystemPolicyLoader{}, loader)
}

func TestNewPolicyLoader_S3(t *testing.T) {
	t.Setenv("S3_BUCKET", "test")

	loader, err := policyloader.NewPolicyLoader(context.TODO())
	assert.NoError(t, err)
	assert.IsType(t, &policyloader.S3PolicyLoader{}, loader)
}
