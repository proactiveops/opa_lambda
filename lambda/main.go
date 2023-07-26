/**
* OPA Lambda
*
* MIT License. See LICENSE file for details.
 */
package main

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"opa_lambda/policyevaluator"
	"opa_lambda/policyloader"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

// A LambdaRequest is the event used to invoke the Lambda function.
type LambdaEvent struct {
	PolicyName string           `json:"policy"`  // The name of the OPA policy to check.
	Payload    *json.RawMessage `json:"payload"` // The payload to evaluate the policy against.
}

type LambdaResponse struct {
	Output interface{} `json:"output"` // The output of the policy evaluation.
	Error  error       `json:"error"`  // The error, if any, that occurred during policy evaluation.
}

// Handle requests for policy evaluation when running on AWS Lambda.
func handleLambda(ctx context.Context, req LambdaEvent) (output LambdaResponse, err error) {
	log.SetFormatter(&log.JSONFormatter{})

	log.Infof("Evaluating policy: %s", req.PolicyName)

	pl, err := policyloader.NewPolicyLoader(ctx)
	if err != nil {
		log.Error(err)
		return LambdaResponse{Error: err}, err
	}

	pe := policyevaluator.NewPolicyEvaluator(pl)

	var result *policyevaluator.EvaluationResult
	if result, err = pe.EvaluatePolicy(ctx, req.PolicyName, *req.Payload); err != nil {
		log.Error(err)
		return LambdaResponse{Error: err}, err
	}

	output = LambdaResponse{Output: result.Value}
	return
}

// Handle requests when testing locally.
func handleLocal() {
	log.SetFormatter(&log.TextFormatter{})

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Unable to read input from stdin")
	}

	ctx := context.Background()
	pl, err := policyloader.NewPolicyLoader(ctx)
	if err != nil {
		log.Fatal(err)
	}

	pe := policyevaluator.NewPolicyEvaluator(pl)

	result, err := pe.EvaluatePolicy(ctx, os.Args[1], input)
	if err != nil {
		log.Fatal(err)
	}

	output, err := json.Marshal(result.Value)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(output))
}

func main() {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		// Lambda Environment
		lambda.Start(handleLambda)
	} else {
		// Local development
		handleLocal()
	}
}
