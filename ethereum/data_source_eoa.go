package ethereum

import (
	"context"
	"encoding/hex"

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
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signer": {
				Type:     schema.TypeString,
				Computed: true,
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

	priv, err := key.MarshallPrivateKey()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(key.Address().String())
	d.Set("address", key.Address().String())
	d.Set("signer", hex.EncodeToString(priv))
	return nil
}
