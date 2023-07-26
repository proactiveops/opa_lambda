variable "enable_tracing" {
  description = "Enable AWS X-Ray tracing."
  type        = bool
  default     = false
}

variable "environment" {
  description = "The environment (e.g., 'dev', 'prod') the function will be deployed into."
  type        = string
}

variable "function_name" {
  description = "The name of the function. The value of var.environment will be appended to this name."
  type        = string
  default     = "opa-lambda"
}

variable "s3_bucket" {
  description = "The name of an existing S3 bucket used for storing rego files. If omitted a new S3 bucket will be created."
  type        = string
  default     = ""
}

variable "security_group_ids" {
  description = "The security group IDs for the Lambda function. Skip if you won't want to run the lambda within a VPC."
  type        = list(string)
  default     = []
}

variable "subnet_ids" {
  description = "The subnet IDs for the Lambda function. Skip if you won't want to run the lambda within a VPC."
  type        = list(string)
  default     = []
}

variable "tags" {
  description = "The tags to apply to the resources."
  type        = map(string)
  default     = {}
}

locals {
  create_bucket = var.s3_bucket == ""
  function_name = "${var.function_name}-${var.environment}"

  s3_bucket = local.create_bucket ? aws_s3_bucket.this[0] : data.aws_s3_bucket.this[0]

  tags = merge(var.tags, {
    Environment = var.environment
  })
}
