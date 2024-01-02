package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
)

func ContractDeploymentResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"artifact": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"input": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"signer": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gas_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"contract_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CreateContext: resourceContractDeploymentCreate,
		ReadContext:   resourceContractDeploymentRead,
		DeleteContext: resourceContractDeploymentDelete,
	}
}

func resourceContractDeploymentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	signer, err := hex.DecodeString(d.Get("signer").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	txn := &transaction{
		Signer: signer,
	}

	artifact, err := resolveContract(d.Get("artifact").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	code, err := hex.DecodeString(artifact.Bytecode.Object[2:])
	if err != nil {
		return diag.FromErr(err)
	}

	if cons := artifact.Abi.Constructor; cons != nil {
		var inputs interface{}
		if rawInputs, ok := d.GetOk("input"); ok {
			inputs, err = decodeInputs(rawInputs)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to decode inputs: %v", err))
			}
		} else {
			inputs = []interface{}{}
		}

		inputsBytecode, err := cons.Inputs.Encode(inputs)
		if err != nil {
			return diag.FromErr(err)
		}
		code = append(code, inputsBytecode...)
	}

	txn.Input = code

	client := meta.(*client)
	hash, receipt, err := client.sendTransaction(txn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash.String())
	d.Set("hash", hash.String())
	d.Set("gas_used", receipt.GasUsed)
	d.Set("contract_address", receipt.ContractAddress.String())
	d.Set("block_num", int(receipt.BlockNumber))

	return nil
}

func resourceContractDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client)

	hash := d.Id()
	receipt, err := client.Http().GetTransactionReceipt(ethgo.HexToHash(hash))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash)
	d.Set("hash", hash)
	d.Set("gas_used", receipt.GasUsed)
	d.Set("contract_address", receipt.ContractAddress.String())

	return nil
}

func resourceContractDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
