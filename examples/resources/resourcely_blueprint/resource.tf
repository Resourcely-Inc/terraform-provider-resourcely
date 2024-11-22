resource "resourcely_blueprint" "private_s3_bucket" {
  name        = "Private S3 Bucket"
  description = "Creates a private, S3 bucket with configuration versioning."

  cloud_provider = "PROVIDER_AMAZON"
  categories     = ["BLUEPRINT_BLOB_STORAGE"]

  is_published = true

  content = <<-EOT
    ---
    constants:
      __name: "{{ bucket }}_{{ __guid }}"
    ---
    resource "aws_s3_" "{{ __name }}" {
      bucket = {{ bucket }}
    }

    resource "aws_s3_bucket_public_access_block" "{{ __name }}" {
      bucket = aws_s3_bucket.{{ __name }}.id

      block_public_acls       = true
      block_public_policy     = true
      ignore_public_acls      = true
      restrict_public_buckets = true
    }

    resource "aws_s3_bucket_ownership_controls" "{{ __name }}" {
      bucket = aws_s3_bucket.{{ __name }}.id

      rule {
        object_ownership = "BucketOwnerEnforced"
      }
    }
  EOT
}
