package ethereum

// getGasPrice

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func datasourceContractCode() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceContractCodeRead,
		Schema: map[string]*schema.Schema{
			"addr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func datasourceContractCodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	addrStr := d.Get("addr").(string)

	var addr ethgo.Address
	if err := addr.UnmarshalText([]byte(addrStr)); err != nil {
		return diag.FromErr(err)
	}

	code, err := m.(*client).httpClient.Eth().GetCode(addr, ethgo.Latest)
	if err != nil {
		return diag.FromErr(err)
	}
	code = strings.TrimPrefix(code, "0x")

	d.SetId(addrStr)
	d.Set("code", code)
	return nil
}
