package ethereum

import (
	"context"

	"github.com/umbracle/ethgo/jsonrpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http://localhost:8545",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ethereum_eoa":   datasourceEoa(),
			"ethereum_block": datasourceBlock(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ethereum_transaction":         TransactionResource(),
			"ethereum_contract_deployment": ContractDeploymentResource(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		client, err := jsonrpc.NewClient(d.Get("host").(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return client, nil
	}

	return provider
}
