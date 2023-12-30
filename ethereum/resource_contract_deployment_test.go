package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func checkContractDeployed() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			"ethereum_contract_deployment.deploy", "hash"),
		resource.TestCheckResourceAttrSet(
			"ethereum_contract_deployment.deploy", "contract_address"),
		resource.TestCheckResourceAttrSet(
			"ethereum_contract_deployment.deploy", "gas_used"),
		resource.TestCheckResourceAttrSet(
			"ethereum_contract_deployment.deploy", "block_num"),
	)
}

func TestAccContractDeployment_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "ethereum_eoa" "account" {
						mnemonic = "test test test test test test test test test test test junk"
					}

					resource "ethereum_contract_deployment" "deploy" {
						signer = data.ethereum_eoa.account.signer

						artifact     = "../testcases/out:Hello"

						input = [
						  "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
						]
					}
					`,
				Check: checkContractDeployed(),
			},
		},
	})
}

func TestAccContractDeployment_Inputs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					data "ethereum_eoa" "account" {
						mnemonic = "test test test test test test test test test test test junk"
					}

					resource "ethereum_contract_deployment" "deploy" {
						signer = data.ethereum_eoa.account.signer

						artifact = "../testcases/out:Inputs"

						input = [
							"0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5",
							"1",
							"0xcc84c3b12f6ae46a791f06a0297bb2d9e60d1d4e0f7c0aff2f5be06cea9189d4",
							jsonencode({
							  "number" = "1"
							})
						]
					}
					`,
				Check: checkContractDeployed(),
			},
		},
	})
}
