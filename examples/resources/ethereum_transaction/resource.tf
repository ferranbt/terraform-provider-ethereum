data "ethereum_eoa" "account" {
  mnemonic = "test test test test test test test test test test test junk"
}

resource "ethereum_eoa" "target" {}

// Send 1 gwei to the target account
resource "ethereum_transaction" "update" {
  signer = data.ethereum_eoa.account.signer
  to     = resource.ethereum_eoa.target.address
  value  = "1 gwei"
}

// Contract call using the ABI artifact and method
resource "ethereum_transaction" "update" {
  signer = data.ethereum_eoa.account.signer
  to     = "0x..."

  artifact = "../testcases/out:Inputs"
  method   = "applyFunc"

  input = [
    "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe6",
    "2",
    "0xaa84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
    jsonencode({
      "number" = "3"
    })
  ]
}
