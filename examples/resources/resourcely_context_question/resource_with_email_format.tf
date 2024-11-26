resource "resourcely_context_question" "owner" {
  label  = "owner"
  prompt = "What is the email address of the team or person responsible for this infrastructure?"

  qtype         = "QTYPE_TEXT"
  answer_format = "ANSWER_EMAIL"

  blueprint_categories = [
    "BLUEPRINT_COMPUTE",
    "BLUEPRINT_DATABASE",
    "BLUEPRINT_NETWORKING",
    "BLUEPRINT_SERVERLESS_COMPUTE"
  ]

  priority = 0
}
