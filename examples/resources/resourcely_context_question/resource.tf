resource "resourcely_context_question" "data_classification" {
  label  = "classification"
  prompt = "What kind of data is stored by this infrastructure?"

  qtype = "QTYPE_MULTI_SELECT"
  answer_choices = [
    { label = "financial" },
    { label = "pii" },
    { label = "proprietary" },
    { label = "public" },
  ]

  blueprint_categories = ["BLUEPRINT_BLOB_STORAGE", "BLUEPRINT_DATABASE"]

  priority = 1
}
