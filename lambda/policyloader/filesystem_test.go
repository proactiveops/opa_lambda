// policyloader/filesystem_test.go
package policyloader_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"opa_lambda/policyloader"
)

func TestFilesystemLoadPolicy(t *testing.T) {
	ctx := context.TODO()
	loader := &policyloader.FilesystemPolicyLoader{}

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	testPolicy := "package test\n\nallow = true\n"

	policyPath := cwd + "/policies"
	os.Mkdir(policyPath, 0700)
	policyFilePath := policyPath + "/test.rego"

	os.WriteFile(policyFilePath, []byte(testPolicy), 0600)

	policy, err := loader.LoadPolicy(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, testPolicy, policy)

	defer os.RemoveAll(policyPath)
}

func TestFilesystemLoadPolicyInvalidName(t *testing.T) {
	ctx := context.TODO()
	loader := &policyloader.FilesystemPolicyLoader{}

	_, err := loader.LoadPolicy(ctx, "invalid/name")
	assert.Error(t, err)
}

func TestFilesystemLoadPolicyNotFound(t *testing.T) {
	ctx := context.TODO()
	loader := &policyloader.FilesystemPolicyLoader{}

	_, err := loader.LoadPolicy(ctx, "not-found")
	assert.Error(t, err)
}
