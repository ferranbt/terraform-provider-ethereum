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
		Description: "Create a new EOA wallet.",
		Schema: map[string]*schema.Schema{
			"mnemonic": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"privkey"},
				Description:   "The mnemonic of the wallet to use.",
			},
			"privkey": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"mnemonic"},
				Description:   "The private key of the wallet to use.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the wallet.",
			},
			"signer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The signer of the wallet. This is the private key of the wallet.",
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
