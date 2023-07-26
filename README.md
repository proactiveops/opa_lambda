# Lambda Function for Open Policy Agent

This Lambda function evaluates JSON objects against [Open Policy Agent's (OPA)](https://openpolicyagent.org/) [Rego policy files](https://www.openpolicyagent.org/docs/latest/policy-language/). The application uses a serverless architecture. The Rego files are stored in an S3 Bucket and cached in memory to improve performance.

The function is designed for high throughput, low latency policy evaluation.

This is an alternative to using [OPA's REST API server](https://www.openpolicyagent.org/docs/latest/rest-api/). If you have sustained, consistent policy evaluation traffic the REST API is worth considering. On the other hand if your traffic is uneven or sporadic, this lambda function is likely to be more cost effective and efficient.

To learn more about working with this Lambda function, check out the [Proactive Ops blog post about evaluating rego files with AWS Lambda](https://proactiveops.io/archive/serverless-policy-as-code).

## Setup

The lambda function and associated resources are packaged in this repository as a Terraform module. This is the easiest way to deploy the resources into your AWS account.

### Terraform

To deploy the lambda using Terraform add the following snippet to your existing project:

```hcl
# opa-lambda.tf

module "opa_lambda" {
  source = "git@github.com:proactiveops/opa-lambda.git?ref=USE-LATEST-HASH-HERE"

  enable_tracing = true # This is needed if you want XRay tracing enabled. It is off by default.

  environment   = "dev"                    # This is only mandatory variable. You must specify the name of the environment for this instance.
  function_name = "my-opa-fn-or-something" # The name of the lambda function. The default value is "opa". We append the value of var.environment to this value for naming the function.

  s3_bucket = aws_s3_bucket.opa_policies.arn # Only needed if you have an existing bucket. If omitted a bucket will be created.

  security_group_ids = [ # Optional security groups for the lambda execution. Subnets must also be supplied. Use encouraged but not required.
    "sg-d15ab1ed",
    "sg-deadend2",
  ]

  subnet_ids = [ # Subnets used by the lambda function. Should be a private subnet. Security group IDs also required. Optional but encouraged.
    "subnet-20ldcafef00d2111",
    "subnet-deadd00dcafef00d"
  ]

  # AWS Tags to apply to all resources provisioned by module.
  tags = var.tags
}
```

## Example Invocation

To invoke the function use the following payload:

```json
{
  "policy": "module.example",
  "payload": {
    "..." : "..."
  }
}
```

## Local Development

The function has a local development mode which helps users test their policy files locally before deploying them.

To invoke the function locally make sure your policy files are in the path `policies` and then run:

```sh
cd lambda
cat example.json | go run main.go <policy_name>
```

### S3 Integration

Set the `S3_BUCKET` environment variable to the name of your S3 bucket if you want to use an S3 bucket to store policies. If not provided, the application will load policies from the local filesystem. For this to work, the AWS credentials must be configured in your environment or in the `~/.aws/credentials` file.

To run the application locally, execute the following command:

```sh
cd lambda
cat example.json | go run main.go <policy_name>
```

### Building

To build the application, run:

```sh
cd lambda
go build -o opa_lambda main.go
```

### Testing

To run the tests, execute:

``` sh
go test ./...
```

## Contributing

Please submit issues and pull requests for any bug reports, feature requests, or other contributions.