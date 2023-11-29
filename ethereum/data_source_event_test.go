package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatasource_Events(t *testing.T) {
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
					artifact = "../testcases/out:WithEvents"
				}

				resource "ethereum_transaction" "update" {
					signer = data.ethereum_eoa.account.signer
					to = resource.ethereum_contract_deployment.deploy.contract_address

					artifact = "../testcases/out:WithEvents"
					method = "applyFunc"
				}

				data "ethereum_event" "res" {
					hash = ethereum_transaction.update.hash
					artifact = "../testcases/out:WithEvents"
					event = "One"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ethereum_event.res", "logs.num", "1"),
					resource.TestCheckResourceAttrSet(
						"data.ethereum_event.res", "address"),
				),
			},
		},
	})
}
