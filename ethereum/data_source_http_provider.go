package ethereum

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceHTTPProvider() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceHTTPProviderProviderRead,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  defaultHost,
			},
			"chain_id": {
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func datasourceHTTPProviderProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	url := d.Get("url").(string)

	clt, err := newClient(url)
	if err != nil {
		return diag.FromErr(err)
	}

	chainId, err := clt.Http().ChainID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("chain_id", chainId.Int64())
	d.SetId(url)
	return nil
}
