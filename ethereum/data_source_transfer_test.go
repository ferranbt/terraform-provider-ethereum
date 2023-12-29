package ethereum

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatasource_Transfer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "ethereum_transfer" "xx" {
					start_block = "1000"
					end_block = "latest"
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
