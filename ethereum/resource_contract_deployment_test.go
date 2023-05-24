package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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

					artifact_path     = "../testcases/out"
					artifact_contract = "Hello"
				  
					input = [
					  "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5"
					]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ethereum_contract_deployment.deploy", "hash"),
					resource.TestCheckResourceAttrSet(
						"ethereum_contract_deployment.deploy", "contract_address"),
				),
			},
		},
	})
}
