data "ethereum_eoa" "account" {
  mnemonic = "test test test test test test test test test test test junk"
}

resource "ethereum_contract_deployment" "deploy" {
  signer = data.ethereum_eoa.account.signer

  artifact = "../testcases/out:Hello"

  input = [
    "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
  ]
}
