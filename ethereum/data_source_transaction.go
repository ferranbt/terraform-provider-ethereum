package ethereum

import (
	"context"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func datasourceTransaction() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceTransactionRead,
		Schema: map[string]*schema.Schema{
			"hash": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"to": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gas": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gas_price": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nonce": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"input": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceTransactionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hash := d.Get("hash").(string)

	txn, err := m.(*client).httpClient.Eth().GetTransactionByHash(ethgo.HexToHash(hash))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash)
	d.Set("hash", hash)
	d.Set("from", txn.From.String())
	d.Set("to", txn.To.String())
	d.Set("value", txn.Value.String())
	d.Set("gas", int(txn.Gas))
	d.Set("gas_price", int(txn.GasPrice))
	d.Set("nonce", int(txn.Nonce))
	d.Set("input", hex.EncodeToString(txn.Input))

	return nil
}
