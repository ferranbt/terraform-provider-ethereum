package ethereum

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func TransactionResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"to": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"input": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"signer": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gas_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		CreateContext: resourceTransactionCreate,
		ReadContext:   resourceTransactionRead,
		DeleteContext: resourceTransactionDelete,
	}
}

func resourceTransactionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	signer, err := hex.DecodeString(d.Get("signer").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	txn := &transaction{
		Signer: signer,
	}

	if val, ok := d.GetOk("to"); ok {
		addr := ethgo.HexToAddress(val.(string))
		txn.To = &addr
	}
	if val, ok := d.GetOk("value"); ok {
		txn.Value = new(big.Int).SetInt64(int64(val.(float64)))
	}
	if val, ok := d.GetOk("input"); ok {
		buf, err := hex.DecodeString(val.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		txn.Input = buf
	}

	client := meta.(*client)
	hash, receipt, err := client.sendTransaction(txn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash.String())
	d.Set("hash", hash.String())
	d.Set("gas_used", receipt.GasUsed)

	return nil
}

func resourceTransactionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTransactionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
