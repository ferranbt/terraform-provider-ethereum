package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func datasourceCall() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceCallRead,
		Description: "Call a contract method.",
		Schema: map[string]*schema.Schema{
			"to": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The address of the contract to call.",
			},
			"artifact": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The artifact of the contract to call.",
			},
			"method": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the method in the contract to call. ",
			},
			"input": {
				Type:        schema.TypeList,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The inputs of the contract method to call.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The outputs of the contract call.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func datasourceCallRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	addr := ethgo.HexToAddress(d.Get("to").(string))
	callMsg := &ethgo.CallMsg{
		To: &addr,
	}

	artifact, err := resolveContract(d.Get("artifact").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var inputs interface{}
	if rawInputs, ok := d.GetOk("input"); ok {
		inputs, err = decodeInputs(rawInputs)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to decode inputs: %v", err))
		}
	} else {
		inputs = []interface{}{}
	}

	methodName := d.Get("method").(string)
	method, ok := artifact.Abi.Methods[methodName]
	if !ok {
		return diag.FromErr(fmt.Errorf("method '%s' not found", methodName))
	}

	buf, err := method.Encode(inputs)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to abi encode: %v", err))
	}
	callMsg.Data = buf

	res, err := m.(*client).httpClient.Eth().Call(callMsg, ethgo.Latest)
	if err != nil {
		return diag.FromErr(err)
	}

	resBuf, err := hex.DecodeString(strings.TrimPrefix(res, "0x"))
	if err != nil {
		return diag.FromErr(err)
	}
	output, err := method.Outputs.Decode(resBuf)
	if err != nil {
		return diag.FromErr(err)
	}

	outputMap, ok := output.(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf("incorrect output?"))
	}

	resMap := map[string]string{}
	for k, v := range outputMap {
		resMap[k] = fmt.Sprint(v)
	}

	d.Set("output", resMap)
	d.SetId(addr.String())

	return nil
}
