package policyloader_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"opa_lambda/policyloader"
)

func TestKeyToFilename(t *testing.T) {
	tests := []struct {
		key      string
		filename string
		err      error
	}{
		{
			key:      "policy",
			filename: "policy.rego",
			err:      nil,
		},
		{
			key:      "policy.rego",
			filename: "policy/rego.rego",
			err:      nil,
		},
		{
			key:      "policy.json",
			filename: "policy/json.rego",
			err:      nil,
		},
		{
			key:      "policy.json.rego",
			filename: "policy/json/rego.rego",
			err:      nil,
		},
		{
			key:      "policy/with/slashes",
			filename: "",
			err:      &policyloader.InvalidKeyNameError{Key: "policy/with/slashes"},
		},
	}

	for _, test := range tests {
		filename, err := policyloader.KeyToFilename(test.key)
		assert.Equal(t, test.filename, filename)
		assert.Equal(t, test.err, err)
	}
}
