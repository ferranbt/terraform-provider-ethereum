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
		Description: "Filter transactions from a block range.",
		Schema: map[string]*schema.Schema{
			"start_block": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The block number to start the filter from. ",
			},
			"limit_blocks": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of blocks to filter. ",
			},
			"from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The address to filter transactions from.",
			},
			"to": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The address to filter transactions to.",
			},
			"is_transfer": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to filter only transfer transactions.",
			},
			"hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hash of the transaction that matches the filter",
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
