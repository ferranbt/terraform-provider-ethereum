package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlock_basic(t *testing.T) {
	blockIsComplete := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(
			"data.ethereum_block.block", "hash"),
		resource.TestCheckResourceAttrSet(
			"data.ethereum_block.block", "number"),
		resource.TestCheckResourceAttrSet(
			"data.ethereum_block.block", "timestamp"),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_block" "block" {
					tag = "latest"
				}
				`,
				Check: blockIsComplete,
			},
			{
				Config: `
				data "ethereum_block" "block" {
					number = 1
				}
				`,
				Check: blockIsComplete,
			},
		},
	})
}
