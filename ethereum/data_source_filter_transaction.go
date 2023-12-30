package ethereum

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func datasourceFilterTransaction() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceFilterTransactionRead,
		Schema: map[string]*schema.Schema{
			"start_block": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"limit_blocks": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_transfer": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceFilterTransactionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var input filterTransactionInput

	startBlock := d.Get("start_block").(int)
	if startBlock < 0 {
		startBlock = 0
	}
	input.StartBlock = uint64(startBlock)

	if v, ok := d.GetOk("from"); ok {
		from := ethgo.HexToAddress(v.(string))
		input.From = &from
	}
	if v, ok := d.GetOk("to"); ok {
		to := ethgo.HexToAddress(v.(string))
		input.To = &to
	}
	if v, ok := d.GetOk("is_transfer"); ok {
		isTransfer := v.(bool)
		input.IsTransfer = &isTransfer
	}
	if v, ok := d.GetOk("limit_blocks"); ok {
		limitBlocks := uint64(v.(int))
		input.LimitBlocks = &limitBlocks
	}

	client := m.(*client)
	hash, err := client.filterTransactions(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash.String())
	d.Set("hash", hash.String())

	return nil
}
