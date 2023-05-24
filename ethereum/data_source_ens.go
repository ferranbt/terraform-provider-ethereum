package ethereum

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo/ens"
)

func datasourceENS() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceENSRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceENSRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ens, err := ens.NewENS(ens.WithClient(m.(*client).httpClient))
	if err != nil {
		return diag.FromErr(err)
	}

	addr, err := ens.Resolve(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(addr.String())
	d.Set("address", addr.String())
	return nil
}
