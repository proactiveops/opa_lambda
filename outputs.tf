output "lambda_function_name" {
  description = "The name of the Lambda function."
  value       = aws_lambda_function.this.function_name
}

output "lambda_function_arn" {
  description = "The Amazon Resource Name (ARN) of the Lambda function."
  value       = aws_lambda_function.this.arn
}

output "lambda_role_arn" {
  description = "The Amazon Resource Name (ARN) of the IAM role used by the Lambda function."
  value       = aws_iam_role.lambda_role.arn
}

output "log_group_name" {
  description = "The name of the CloudWatch log group for the Lambda function."
  value       = aws_cloudwatch_log_group.this.name
}

output "s3_bucket_name" {
  description = "The name of the S3 bucket used to store the policy files."
  value       = local.s3_bucket.id
}
