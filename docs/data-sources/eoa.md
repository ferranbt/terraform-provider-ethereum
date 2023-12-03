---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ethereum_eoa Data Source - terraform-provider-ethereum"
subcategory: ""
description: |-
  The `eoa` data source declares a wallet by its mnemonic address.
---

# ethereum_eoa (Data Source)

The `eoa` data source declares a wallet by its mnemonic address.

## Example usage

```
data "ethereum_eoa" "account" {
	mnemonic = "test test test test test test test test test test test junk"
}
```

## Schema

### Optional

- `mnemonic` (String): Mnemonic of the private key.

### Read-Only

- `address` (String)
- `id` (String) The ID of this resource.
- `signer` (String)