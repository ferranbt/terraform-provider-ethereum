package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEOA_mnemonic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					mnemonic = "test test test test test test test test test test test junk"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ethereum_eoa.account", "signer"),
					resource.TestCheckResourceAttr(
						"data.ethereum_eoa.account", "address", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				),
			},
		},
	})
}

func TestAccEOA_privKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_eoa" "account" {
					privkey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ethereum_eoa.account", "signer"),
					resource.TestCheckResourceAttr(
						"data.ethereum_eoa.account", "address", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				),
			},
		},
	})
}
