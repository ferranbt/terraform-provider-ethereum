package ethereum

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func datasourceBlock() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceBlockRead,
		Schema: map[string]*schema.Schema{
			"number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func datasourceBlockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client)

	var block *ethgo.Block
	var err error

	if hashStr, ok := d.GetOk("hash"); ok {
		// resolve the block by its hash
		block, err = client.Http().GetBlockByHash(ethgo.HexToHash(hashStr.(string)), true)

	} else if tagStr, ok := d.GetOk("tag"); ok {
		// resolve the block with a tag ('latest', 'finalized')
		var tag ethgo.BlockNumber
		if tagStr == "latest" {
			tag = ethgo.Latest
		}
		block, err = client.Http().GetBlockByNumber(tag, true)

	} else if numVal, ok := d.GetOk("number"); ok {
		// resolve the block by number
		block, err = client.Http().GetBlockByNumber(ethgo.BlockNumber(numVal.(int)), true)

	}

	if err != nil {
		return diag.FromErr(err)
	}
	if block == nil {
		return diag.FromErr(fmt.Errorf("block not found"))
	}

	// fill-in all the values from the result, including the
	// optional ones from the input
	d.SetId(block.Hash.String())
	d.Set("hash", block.Hash.String())
	d.Set("number", block.Number)
	d.Set("timestamp", block.Timestamp)

	return nil
}
