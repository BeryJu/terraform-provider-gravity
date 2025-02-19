---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gravity_user Resource - terraform-provider-gravity"
subcategory: ""
description: |-
  
---

# gravity_user (Resource)



## Example Usage

```terraform
# Simple user (read-only)

resource "gravity_user" "read-only" {
  username = "read-only"
  permissions = jsonencode([
    {
      path    = "/*",
      methods = ["get", "head"]
    }
  ])
}

# Admin user

resource "gravity_user" "admin" {
  username = "admin"
  permissions = jsonencode([
    {
      path    = "/*",
      methods = ["*"]
    }
  ])
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `username` (String)

### Optional

- `permissions` (String) Defaults to `[]`.

### Read-Only

- `id` (String) The ID of this resource.
