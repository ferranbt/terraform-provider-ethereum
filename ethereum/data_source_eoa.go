package ethereum

import (
	"context"
	"encoding/hex"
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
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"privkey"},
			},
			"privkey": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"mnemonic"},
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
	var (
		key *wallet.Key
		err error
	)

	if mnemonic, ok := d.GetOk("mnemonic"); ok {
		key, err = wallet.NewWalletFromMnemonic(mnemonic.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if privKeyStr, ok := d.GetOk("privkey"); ok {
		privKey, err := hex.DecodeString(privKeyStr.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		key, err = wallet.NewWalletFromPrivKey(privKey)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.FromErr(fmt.Errorf("no wallet key provided"))
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
