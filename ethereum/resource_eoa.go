package ethereum

import (
	"context"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo/wallet"
)

func EOAResource() *schema.Resource {
	return &schema.Resource{
		Description: "Create a new EOA wallet.",
		Schema: map[string]*schema.Schema{
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
		CreateContext: resourceEOACreate,
		ReadContext:   resourceEOARead,
		DeleteContext: resourceEOADelete,
	}
}

func resourceEOACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	key, err := wallet.GenerateKey()
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

func resourceEOARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceEOADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
