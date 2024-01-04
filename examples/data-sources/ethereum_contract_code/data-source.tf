data "contract_code" "deposit_contract" {
  address = "0x00000000219ab540356cBB839Cbe05303d7705Fa"

  // validate that the contract exists
  lifecycle {
    postcondition {
      condition     = self.code != ""
      error_message = "Deposit contract does not exist"
    }
  }
}
