---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ethereum_eoa Resource - terraform-provider-ethereum"
subcategory: ""
description: |-
  Create a new EOA wallet.
---

# ethereum_eoa (Resource)

Create a new EOA wallet.

## Example Usage

```terraform
resource "ethereum_eoa" "account" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `address` (String) The address of the wallet.
- `id` (String) The ID of this resource.
- `signer` (String) The signer of the wallet. This is the private key of the wallet.
