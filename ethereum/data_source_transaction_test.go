package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatasource_Transaction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}

				resource "ethereum_eoa" "target" {}

				resource "ethereum_transaction" "update" {
					signer = data.ethereum_eoa.account.signer
					to = resource.ethereum_eoa.target.address
					value = "1 gwei"
				}

				data "ethereum_transaction" "res" {
					hash = ethereum_transaction.update.hash
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ethereum_transaction.res", "from"),
					resource.TestCheckResourceAttrSet(
						"data.ethereum_transaction.res", "to"),
					resource.TestCheckResourceAttrSet(
						"data.ethereum_transaction.res", "value"),
					resource.TestCheckResourceAttrSet(
						"data.ethereum_transaction.res", "gas"),
					resource.TestCheckResourceAttrSet(
						"data.ethereum_transaction.res", "gas_price"),
					resource.TestCheckResourceAttrSet(
						"data.ethereum_transaction.res", "nonce"),
				),
			},
		},
	})
}
