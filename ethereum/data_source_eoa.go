package ethereum

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo/wallet"
)

func datasourceEoa() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceEoaRead,
		Schema: map[string]*schema.Schema{
			"mnemonic": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func datasourceEoaRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mnemonic := d.Get("mnemonic")

	key, err := wallet.NewWalletFromMnemonic(mnemonic.(string))
	if err != nil {
		return diag.FromErr(err)
	}
	fmt.Println(key)

	return nil
}
