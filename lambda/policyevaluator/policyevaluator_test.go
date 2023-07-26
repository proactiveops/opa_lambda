package policyevaluator

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const emptyRegoPolicy = `package empty`

const malformedRegoPolicy = `package bad

default allow = garbage-value

allow {
    input.action = "break"`

const exampleRegoPolicy = `package valid

default allow = false

allow {
    input.user == "alice"
    input.action == "read"
}`

type mockPolicyLoader struct{}

func (m *mockPolicyLoader) LoadPolicy(ctx context.Context, policyID string) (string, error) {
	if policyID == "valid" {
		return exampleRegoPolicy, nil
	}
	if policyID == "malformed" {
		return malformedRegoPolicy, nil
	}
	if policyID == "empty" {
		return emptyRegoPolicy, nil
	}
	return "", errors.New("policy not found")
}

func TestPolicyEvaluator(t *testing.T) {
	mockLoader := &mockPolicyLoader{}
	eval := NewPolicyEvaluator(mockLoader)

	payload := json.RawMessage(`{"user": "alice", "action": "read"}`)
	result, err := eval.EvaluatePolicy(context.Background(), "valid", payload)
	assert.NoError(t, err)
	value, ok := result.Value.(map[string]interface{})
	assert.True(t, ok)
	assert.True(t, value["allow"].(bool))

	payload = json.RawMessage(`{"user": "bob", "action": "write"}`)
	result, err = eval.EvaluatePolicy(context.Background(), "valid", payload)
	assert.NoError(t, err)
	value, ok = result.Value.(map[string]interface{})
	assert.True(t, ok)
	assert.False(t, value["allow"].(bool))
}

func TestPolicyEvaluator_BrokenJson(t *testing.T) {
	mockLoader := &mockPolicyLoader{}
	eval := NewPolicyEvaluator(mockLoader)

	payload := json.RawMessage(`{"this is invalid json"}`)
	_, err := eval.EvaluatePolicy(context.Background(), "valid", payload)
	assert.Error(t, err)
}

func TestPolicyEvaluator_MissingPolicy(t *testing.T) {
	mockLoader := &mockPolicyLoader{}
	eval := NewPolicyEvaluator(mockLoader)

	payload := json.RawMessage(`{"this": "is", "valid": "json"}`)
	_, err := eval.EvaluatePolicy(context.Background(), "missing", payload)
	assert.Error(t, err)
}

func TestPolicyEvaluator_BadPolicy(t *testing.T) {
	mockLoader := &mockPolicyLoader{}
	eval := NewPolicyEvaluator(mockLoader)

	payload := json.RawMessage(`{}`)
	_, err := eval.EvaluatePolicy(context.Background(), "malformed", payload)
	assert.Error(t, err)
}

func TestPolicyEvaluator_EmptyPayload(t *testing.T) {
	mockLoader := &mockPolicyLoader{}
	eval := NewPolicyEvaluator(mockLoader)

	payload := json.RawMessage(`{}`)
	_, err := eval.EvaluatePolicy(context.Background(), "valid", payload)
	assert.NoError(t, err)
	// assert.Equal(t, false, result.Value.(map[string]interface{})[0].(map[string]interface{})["allow"].(bool))
}

func TestPolicyEvaluator_EmptyPolicy(t *testing.T) {
	mockLoader := &mockPolicyLoader{}
	eval := NewPolicyEvaluator(mockLoader)

	payload := json.RawMessage(`{}`)
	result, err := eval.EvaluatePolicy(context.Background(), "empty", payload)
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}(map[string]interface{}{}), result.Value.(map[string]interface{}))
}
