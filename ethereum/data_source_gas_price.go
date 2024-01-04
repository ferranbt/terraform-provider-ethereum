package ethereum

// getGasPrice

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceGetGasPrice() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceGetGasPriceRead,
		Description: "Get the gas price at the current block.",
		Schema: map[string]*schema.Schema{
			"gas_price": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The gas price in wei at the current block",
			},
		},
	}
}

func datasourceGetGasPriceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	gasPrice, err := m.(*client).httpClient.Eth().GasPrice()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", gasPrice))
	d.Set("gas_price", int(gasPrice))

	return nil
}
