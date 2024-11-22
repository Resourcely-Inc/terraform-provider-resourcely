resource "resourcely_guardrail" "s3_bucket_naming_convention" {
  name        = "S3 Bucket Naming Convention"
  description = "Ensures that all S3 Buckets comply with our standardized naming convention, promoting consistency and ease of identification across our AWS environments."

  cloud_provider = "PROVIDER_AMAZON"
  category       = "GUARDRAIL_BEST_PRACTICES"

  content = <<-EOT
    GUARDRAIL "S3 Bucket Naming Convention"
      WHEN aws_s3_bucket
        REQUIRE bucket STARTS WITH "mycompany-"
  EOT
}
