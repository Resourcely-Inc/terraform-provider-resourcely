---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "resourcely_global_value Resource - terraform-provider-resourcely"
subcategory: ""
description: |-
  A Resourcely GlobalValue
---

# resourcely_global_value (Resource)

A Resourcely GlobalValue



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) An immutable identifier used to reference this global value in blueprints or guardrails.

Must start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.
- `name` (String) A short display name
- `options` (Attributes List) The list of value options for this global value (see [below for nested schema](#nestedatt--options))
- `type` (String) The type of options in the global value. Can be one of `PRESET_VALUE_TEXT`, `PRESET_VALUE_NUMBER`, `PRESET_VALUE_LIST`, `PRESET_VALUE_OBJECT`

### Optional

- `description` (String) A longer description
- `is_deprecated` (Boolean) True if the global value should not be used in new blueprints or guardrails

### Read-Only

- `id` (String) UUID for this version of the global value
- `series_id` (String) UUID for the global value
- `version` (Number) Version of the global value

<a id="nestedatt--options"></a>
### Nested Schema for `options`

Required:

- `key` (String) An immutable identifier for ths option.

Must start with a lowercase letter in `a-z` and include only characters in `a-z0-9_`.
- `label` (String) A unique short display name
- `value` (String) A JSON encoding of the option's value. This value must match the declared type of the global value.

Example: `value = jsonencode("a")`

Example: `value = jsonencode(["a", "b"])`

Optional:

- `description` (String) A longer description