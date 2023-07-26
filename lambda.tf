resource "aws_iam_role" "lambda_role" {
  name = local.function_name

  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "lambda_vpc_execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole"
  role       = aws_iam_role.lambda_role.name
}

data "aws_iam_policy_document" "lambda_s3" {
  statement {
    sid       = "FindBucket"
    effect    = "Allow"
    resources = ["arn:aws:s3:::*"] # tfsec:ignore:aws-iam-no-policy-wildcards Need to find the bucket.

    actions = [
      "s3:GetBucketLocation",
      "s3:ListAllMyBuckets",
    ]
  }

  statement {
    sid    = "ReadBucket"
    effect = "Allow"

    # tfsec:ignore:aws-iam-no-policy-wildcards Need to read all objects in the bucket.
    resources = [
      local.s3_bucket.arn,
      "${local.s3_bucket.arn}/*",
    ]

    actions = [
      "s3:GetObject"
    ]
  }
}

resource "aws_iam_policy" "lambda_s3" {
  policy = data.aws_iam_policy_document.lambda_s3.json
  name   = "${local.function_name}-s3"
}

resource "aws_iam_role_policy_attachment" "lambda_s3" {
  policy_arn = aws_iam_policy.lambda_s3.arn
  role       = aws_iam_role.lambda_role.name
}

resource "aws_iam_role_policy_attachment" "xray" {
  count      = var.enable_tracing ? 1 : 0
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
  role       = aws_iam_role.lambda_role.name
}

resource "null_resource" "build_lambda" {
  provisioner "local-exec" {
    working_dir = "lambda"
    command     = "GOOS=linux GOARCH=amd64 go build -o opa_lambda ."
  }

  triggers = {
    always_run = timestamp()
  }
}

data "archive_file" "lambda_zip" {
  depends_on  = [null_resource.build_lambda]
  type        = "zip"
  source_file = "${path.module}/lambda/opa_lambda"
  output_path = "lambda.zip"
}

# tfsec:ignore:aws-lambda-enable-tracing Tracing is optional.
resource "aws_lambda_function" "this" {
  function_name = local.function_name
  handler       = "opa_lambda"
  role          = aws_iam_role.lambda_role.arn
  runtime       = "go1.x"

  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256


  environment {
    variables = {
      S3_BUCKET = local.s3_bucket.id
    }
  }

  dynamic "vpc_config" {
    for_each = length(var.security_group_ids) > 0 && length(var.subnet_ids) > 0 ? [1] : []

    content {
      security_group_ids = var.security_group_ids
      subnet_ids         = var.subnet_ids
    }
  }

  dynamic "tracing_config" {
    for_each = var.enable_tracing ? [1] : []
    content {
      mode = "Active"
    }
  }

  tags = local.tags
}

# tfsec:ignore:aws-cloudwatch-log-group-customer-key Nothing sensistive in the logs so AWS managed key is appropriate.
resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/${local.function_name}"
  retention_in_days = 7
}
