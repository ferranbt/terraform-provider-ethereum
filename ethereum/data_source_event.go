package ethereum

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func datasourceEvent() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceEventRead,
		Schema: map[string]*schema.Schema{
			"event": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Required: true,
			},
			"artifact": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logs": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func datasourceEventRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	hash := d.Get("hash").(string)

	artifact, err := resolveContract(d.Get("artifact").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	eventName := d.Get("event").(string)
	event, ok := artifact.Abi.Events[eventName]
	if !ok {
		return diag.FromErr(fmt.Errorf("event '%s' not found", eventName))
	}

	client := m.(*client)
	receipt, err := client.httpClient.Eth().GetTransactionReceipt(ethgo.HexToHash(hash))
	if err != nil {
		return diag.FromErr(err)
	}

	var matchLog *ethgo.Log
	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}
		if event.ID() != log.Topics[0] {
			continue
		}
		matchLog = log
	}

	if matchLog == nil {
		return diag.FromErr(fmt.Errorf("no logs match"))
	}

	result, err := event.ParseLog(matchLog)
	if err != nil {
		return diag.FromErr(err)
	}

	res := map[string]string{}
	for k, v := range result {
		res[k] = fmt.Sprint(v)
	}

	d.Set("logs", res)
	d.Set("address", matchLog.Address.String())
	d.SetId(hash)

	return nil
}
