package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatasource_FilterTransaction(t *testing.T) {
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

				resource "ethereum_transaction" "transfer" {
					signer = data.ethereum_eoa.account.signer
					to = resource.ethereum_contract_deployment.deploy.contract_address
					value = "1 gwei"
				}

				data "ethereum_filter_transaction" "filter" {
					start_block = resource.ethereum_contract_deployment.deploy.block_num
					to = resource.ethereum_transaction.transfer.to
					is_transfer = true
					limit_blocks = 10
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ethereum_filter_transaction.filter", "hash"),
				),
			},
		},
	})
}
