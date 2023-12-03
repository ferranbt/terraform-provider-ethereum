---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ethereum_transaction Resource - terraform-provider-ethereum"
subcategory: ""
description: |-
  The `transaction` resource sends an arbitrary Ethereum transaction.
---

# ethereum_transaction (Resource)

The `transaction` resource sends an arbitrary Ethereum transaction.

## Example usage

```
resource "ethereum_transaction" "example" {
  signer = data.ethereum_eoa.account.signer
  to = "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
  value = 100
}
```

<!-- schema generated by tfplugindocs -->

## Schema

### Required

- `signer` (String): Wallet that makes the transaction.

### Optional

- `input` (String): Input bytes for the transaction.
- `to` (String): Target of the transaction. If null, it is a code deployment.
- `value` (Number): Transfer value for the transaction.

### Read-Only

- `gas_used` (Number): Gas used by the transaction.
- `hash` (String): Hash of the transaction.
- `id` (String) The ID of this resource.