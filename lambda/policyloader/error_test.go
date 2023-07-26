// policyloader/filesystem_test.go
package policyloader_test

import (
	"opa_lambda/policyloader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorFileNotFoundError(t *testing.T) {
	err := &policyloader.FileNotFoundError{Key: "test"}
	assert.Equal(t, "unable to locate policy file: test", err.Error())
}

func TestErrorInvalidKeyNameError(t *testing.T) {
	err := &policyloader.InvalidKeyNameError{Key: "t/e/s/t"}
	assert.Equal(t, "policy key name contains slash: t/e/s/t", err.Error())
}
