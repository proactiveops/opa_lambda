// policyevaluator/policyevaluator.go
package policyevaluator

import (
	"context"
	"encoding/json"

	"opa_lambda/policyloader"

	"github.com/open-policy-agent/opa/rego"
)

// EvaluationResult is the result of evaluating a policy.
type EvaluationResult struct {
	Value interface{} `json:"result"` // The OPA result
}

// PolicyEvaluator evaluates policies.
type PolicyEvaluator struct {
	loader policyloader.PolicyLoader
}

// NewPolicyEvaluator creates a new PolicyEvaluator.
func NewPolicyEvaluator(loader policyloader.PolicyLoader) *PolicyEvaluator {
	return &PolicyEvaluator{loader: loader}
}

// EvaluatePolicy evaluates a policy.
func (pe *PolicyEvaluator) EvaluatePolicy(ctx context.Context, policyName string, raw []byte) (*EvaluationResult, error) {
	var input interface{}
	if err := json.Unmarshal(raw, &input); err != nil {
		return nil, err
	}

	module, err := pe.loader.LoadPolicy(ctx, policyName)
	if err != nil {
		return nil, err
	}

	query, err := rego.New(
		rego.Query("data."+policyName),
		rego.Module(policyName+".rego", module),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}

	result, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return &EvaluationResult{Value: result}, nil
	}

	return &EvaluationResult{Value: result[0].Expressions[0].Value}, nil
}
