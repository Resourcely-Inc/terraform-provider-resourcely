resource "resourcely_global_value" "departments" {
  name        = "Departments"
  key         = "departments"
  description = "All departments within MyCompany"

  type = "PRESET_VALUE_TEXT"
  options = [
    {
      key   = "product",
      label = "Product",
      value = "product",
    },
    {
      key   = "it",
      label = "IT",
      value = "IT",
    },
    {
      key   = "marketing",
      label = "Marketing",
      value = "market",
    },
    {
      key   = "engingeering",
      label = "Engineering",
      value = "eng",
    },
  ]
}
