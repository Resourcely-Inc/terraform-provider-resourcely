---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "resourcely_blueprint Resource - terraform-provider-resourcely"
subcategory: ""
description: |-
  A Resourcely Blueprint
---

# resourcely_blueprint (Resource)

A Resourcely Blueprint



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_provider` (String)
- `content` (String)
- `name` (String)

### Optional

- `categories` (Set of String)
- `description` (String)
- `guidance` (String)
- `labels` (Set of String)
- `scope` (String)

### Read-Only

- `id` (String) UUID for this version.
- `series_id` (String) UUID for the blueprint
- `version` (Number) Specific version of the blueprint