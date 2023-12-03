package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCall_basic(t *testing.T) {
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
					artifact = "../testcases/out:Call"
				}

				data "ethereum_call" "multicall" {
					artifact = "../testcases/out:Call"
					method = "multipleOutput"
					to = resource.ethereum_contract_deployment.deploy.contract_address
				}

				data "ethereum_call" "singlecall" {
					artifact = "../testcases/out:Call"
					method = "multipleOutput"
					to = resource.ethereum_contract_deployment.deploy.contract_address
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ethereum_call.multicall", "output.0", "1"),
					resource.TestCheckResourceAttr(
						"data.ethereum_call.multicall", "output.1", "2"),
					resource.TestCheckResourceAttr(
						"data.ethereum_call.singlecall", "output.0", "1"),
				),
			},
		},
	})
}
