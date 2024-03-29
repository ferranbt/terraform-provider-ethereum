package ethereum

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const defaultHost = "http://localhost:8545"

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultHost,
				Description: "The host of the Ethereum node. Defaults to '" + defaultHost + "'.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ethereum_eoa":                datasourceEoa(),
			"ethereum_block":              datasourceBlock(),
			"ethereum_ens":                datasourceENS(),
			"ethereum_event":              datasourceEvent(),
			"ethereum_call":               datasourceCall(),
			"ethereum_gas_price":          datasourceGetGasPrice(),
			"ethereum_transaction":        datasourceTransaction(),
			"ethereum_filter_transaction": datasourceFilterTransaction(),
			"ethereum_contract_code":      datasourceContractCode(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ethereum_transaction":         TransactionResource(),
			"ethereum_contract_deployment": ContractDeploymentResource(),
			"ethereum_eoa":                 EOAResource(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		client, err := newClient(d.Get("host").(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return client, nil
	}

	return provider
}
