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
		Description: "Get a transaction by hash.",
		Schema: map[string]*schema.Schema{
			"hash": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The hash of the transaction to get.",
			},
			"from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the sender of the transaction. This is calculated from the signature of the transaction.",
			},
			"to": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the receiver of the transaction. This is empty if the transaction is a contract creation transaction.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the transaction. This is the amount of wei transferred from the sender to the receiver.",
			},
			"gas": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The gas limit of the transaction. This is the maximum amount of gas that can be used to execute the transaction. ",
			},
			"gas_price": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The gas price of the transaction. This is the amount of wei that the sender is willing to pay for each unit of gas. ",
			},
			"nonce": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The nonce of the transaction.",
			},
			"input": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The input of the transaction.",
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
