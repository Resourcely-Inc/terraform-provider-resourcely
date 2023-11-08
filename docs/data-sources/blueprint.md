---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "resourcely_blueprint Data Source - terraform-provider-resourcely"
subcategory: ""
description: |-
  A resourcely blueprint
---

# resourcely_blueprint (Data Source)

A resourcely blueprint



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `series_id` (String) UUID for the blueprint

### Read-Only

- `categories` (Set of String)
- `cloud_provider` (String)
- `content` (String)
- `description` (String)
- `guidance` (String)
- `id` (String) UUID for this version.
- `labels` (Set of String)
- `name` (String)
- `scope` (String)
- `version` (Number) Specific version of the blueprint