package ethereum

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

func datasourceBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceBlockRead,
		Schema: map[string]*schema.Schema{
			"number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*jsonrpc.Client)

	var block *ethgo.Block
	var err error

	if val, ok := d.GetOk("number"); ok {
		block, err = client.Eth().GetBlockByNumber(ethgo.BlockNumber(val.(int)), true)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("hash", block.Hash.String())

	return nil
}
