resource "resourcely_context_question" "project_code" {
  label  = "project_code"
  prompt = "What is the project code that this infrastructure should be  associated with?"

  qtype         = "QTYPE_TEXT"
  answer_format = "ANSWER_REGEX"
  regex_pattern = "^[A-Z0-9]{6}$"

  blueprint_categories = [
    "BLUEPRINT_COMPUTE",
    "BLUEPRINT_DATABASE",
    "BLUEPRINT_NETWORKING",
    "BLUEPRINT_SERVERLESS_COMPUTE"
  ]

  priority = 2
}
