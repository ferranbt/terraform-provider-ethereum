package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccENS_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_ens" "resolve" {
					name = "arachnid.eth"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ethereum_ens.resolve", "address", "0xb8c2C29ee19D8307cb7255e1Cd9CbDE883A267d5"),
				),
			},
		},
	})
}
