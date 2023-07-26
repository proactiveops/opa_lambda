data "aws_s3_bucket" "this" {
  count  = local.create_bucket ? 0 : 1
  bucket = var.s3_bucket
}

# tfsec:ignore:aws-s3-enable-bucket-logging We support bring your own bucket (BYOB) if access logging is needed.
resource "aws_s3_bucket" "this" {
  count  = local.create_bucket ? 1 : 0
  bucket = "${data.aws_caller_identity.current.account_id}-${data.aws_region.current.name}-${local.function_name}"
  tags   = local.tags
}

# tfsec:ignore:aws-s3-encryption-customer-key We support bring your own bucket (BYOB) if KMS encryption is needed.
resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  count  = local.create_bucket ? 1 : 0
  bucket = aws_s3_bucket.this[0].id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_versioning" "this" {
  count  = local.create_bucket ? 1 : 0
  bucket = aws_s3_bucket.this[0].id

  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_ownership_controls" "this" {
  count  = local.create_bucket ? 1 : 0
  bucket = aws_s3_bucket.this[0].id

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "this" {
  count  = local.create_bucket ? 1 : 0
  bucket = aws_s3_bucket.this[0].id

  rule {
    id = "purgeOldVersions"

    status = "Enabled"
    noncurrent_version_expiration {
      noncurrent_days = 30
    }
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  count  = local.create_bucket ? 1 : 0
  bucket = aws_s3_bucket.this[0].id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

data "aws_iam_policy_document" "bucket_policy" {
  count = local.create_bucket ? 1 : 0

  statement {
    sid = "AllowIAMAdmin"
    principals {
      type        = "AWS"
      identifiers = [data.aws_caller_identity.current.id]
    }

    actions = [
      "s3:*",
    ]

    resources = [
      aws_s3_bucket.this[0].arn,
      "${aws_s3_bucket.this[0].arn}/*",
    ]
  }

  statement {
    sid = "AllowLambdaRead"
    principals {
      type = "AWS"
      identifiers = [
        aws_iam_role.lambda_role.arn
      ]
    }

    actions = [
      "s3:GetObject",
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.this[0].arn,
      "${aws_s3_bucket.this[0].arn}/*",
    ]
  }
}

resource "aws_s3_bucket_policy" "this" {
  count  = local.create_bucket ? 1 : 0
  bucket = aws_s3_bucket.this[0].id
  policy = data.aws_iam_policy_document.bucket_policy[0].json
}
