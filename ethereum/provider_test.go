package ethereum

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"ethereum": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	// check that the devnet endpoint is available and running
	if _, err := http.Get(defaultHost); err != nil {
		t.Fatal("devnet not reachable")
	}
}
