package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
)

func TransactionResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"artifact": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"method",
				},
				ConflictsWith: []string{
					"function",
				},
			},
			"method": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"artifact",
				},
				ConflictsWith: []string{
					"function",
				},
			},
			"function": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"artifact",
					"method",
				},
			},
			"input": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gas_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"raw_input": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
		},
		CreateContext: resourceTransactionCreate,
		ReadContext:   resourceTransactionRead,
		DeleteContext: resourceTransactionDelete,
	}
}

func resourceTransactionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	signer, err := hex.DecodeString(d.Get("signer").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	txn := &transaction{
		Signer: signer,
	}

	if val, ok := d.GetOk("to"); ok {
		addr := ethgo.HexToAddress(val.(string))
		txn.To = &addr
	}
	if val, ok := d.GetOk("value"); ok {
		txn.Value, err = parseEtherValue(val.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse transfer value '%s': %v", val.(string), err))
		}
	}
	if val, ok := d.GetOk("raw_input"); ok {
		buf, err := hex.DecodeString(val.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		txn.Input = buf
	}
	if val, ok := d.GetOk("gas_limit"); ok {
		gasLimit := val.(int)
		if gasLimit < 0 {
			return diag.FromErr(fmt.Errorf("gas limit cannot be less than 0 but %d found", gasLimit))
		}
		txn.GasLimit = uint64(gasLimit)
	}

	var method *abi.Method

	if val, ok := d.GetOk("artifact"); ok {
		artifact, err := resolveContract(val.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		methodName := d.Get("method").(string)
		method, ok = artifact.Abi.Methods[methodName]
		if !ok {
			return diag.FromErr(fmt.Errorf("method '%s' not found", methodName))
		}
	}
	if val, ok := d.GetOk("function"); ok {
		if method, err = abi.NewMethod(val.(string)); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse function '%s': %v", val.(string), err))
		}
	}

	if method != nil {
		var inputs interface{}
		if rawInputs, ok := d.GetOk("input"); ok {
			inputs, err = decodeInputs(rawInputs)
			if err != nil {
				return diag.FromErr(fmt.Errorf("failed to decode inputs: %v", err))
			}
		} else {
			inputs = []interface{}{}
		}

		buf, err := method.Encode(inputs)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to abi encode: %v", err))
		}
		txn.Input = buf
	}

	client := meta.(*client)
	hash, receipt, err := client.sendTransaction(txn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash.String())
	d.Set("hash", hash.String())
	d.Set("gas_used", receipt.GasUsed)
	d.Set("block_num", int(receipt.BlockNumber))

	return nil
}

func resourceTransactionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client)

	hash := d.Id()
	receipt, err := client.Http().GetTransactionReceipt(ethgo.HexToHash(hash))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash)
	d.Set("hash", hash)
	d.Set("gas_used", receipt.GasUsed)
	d.Set("block_num", int(receipt.BlockNumber))

	return nil
}

func resourceTransactionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
