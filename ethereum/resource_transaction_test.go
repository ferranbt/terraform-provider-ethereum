package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTransaction_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}

				resource "ethereum_transaction" "send_value" {
					signer = data.ethereum_eoa.account.signer
					to = "0x74B73aC4158B64004F8379966052b215E2A5fc77"
					value = 100
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ethereum_transaction.send_value", "hash"),
				),
			},
		},
	})
}
