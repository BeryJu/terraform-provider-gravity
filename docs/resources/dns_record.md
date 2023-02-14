---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gravity_dns_record Resource - terraform-provider-gravity"
subcategory: ""
description: |-
  
---

# gravity_dns_record (Resource)



## Example Usage

```terraform
resource "gravity_dns_zone" "example" {
  # Make sure zone ends with a trailing slash
  zone          = "my-domain.com."
  authoritative = true
  handlers = [
    {
      type = "memory",
    },
    {
      type = "etcd",
    },
  ]
}

resource "gravity_dns_record" "record" {
  zone     = gravity_dns_zone.example.zone
  hostname = "foo"
  uid      = "0"
  data     = "1.1.1.1"
  type     = "A"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `data` (String)
- `hostname` (String)
- `type` (String)
- `uid` (String)
- `zone` (String)

### Optional

- `mx_preference` (Number)
- `srv_port` (Number)
- `srv_priority` (Number)
- `srv_weight` (Number)

### Read-Only

- `fqdn` (String) Generated.
- `id` (String) The ID of this resource.

