# Ethereum Terraform Provider

The [Ethereum provider](https://github.com/ferranbt/terraform-provider-ethereum) allows [Terraform](https://terraform.io) to manage [Ethereum](https://ethereum.org/en/) resources.

## Examples

Create the provider:

```hcl
provider "ethereum" {}
```

Send a transaction:

```hcl
data "ethereum_eoa" "account" {
	mnemonic = "test test test test test test test test test test test junk"
}

resource "ethereum_transaction" "example" {
  signer = data.ethereum_eoa.account.signer
  to = "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
  value = 100
}
```

Deploy a contract:

```hcl
resource "ethereum_contract_deployment" "contract" {
  signer = data.ethereum_eoa.account.signer

  artifact_path     = "../package/artifacts"
  artifact_contract = "Proxy"

  input = [
    "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
  ]
}
```

Query the blockchain:

```hcl
data "ethereum_block" "latest" {
  tag = "latest"
}
```

Query a dns name:

```hcl
data "ethereum_ens" "resolve" {
  name = "arachnid.eth"
}
```
