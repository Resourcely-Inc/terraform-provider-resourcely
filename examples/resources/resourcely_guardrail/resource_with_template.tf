resource "resourcely_guardrail" "s3_bucket_naming_convention_from_template" {
  name        = "S3 Bucket Naming Convention"
  description = "Ensures that all S3 Buckets comply with our standardized naming convention, promoting consistency and ease of identification across our AWS environments."

  cloud_provider = "PROVIDER_AMAZON"
  category       = "GUARDRAIL_BEST_PRACTICES"

  guardrail_template_series_id = "4909a93c-b248-4e5a-bff5-7cc101702351"
  guardrail_template_inputs = jsonencode({
    prefix   = "mycompany-"
    approver = "default"
  })
}
