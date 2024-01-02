package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGetCode_DeployedContract(t *testing.T) {
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

				data "ethereum_contract_code" "deploy" {
					addr = ethereum_contract_deployment.deploy.contract_address
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ethereum_contract_code.deploy", "code"),
				),
			},
		},
	})
}

func TestAccGetCode_Empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_contract_code" "deploy" {
					addr = "0x0101010101010101010101010101010101010101"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ethereum_contract_code.deploy", "code", ""),
				),
			},
		},
	})
}
