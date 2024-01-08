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
		Description: "Send a transaction.",
		Schema: map[string]*schema.Schema{
			"to": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The address of the contract to call.",
			},
			"artifact": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The ABI artifact of the contract to call.",
				RequiredWith: []string{
					"method",
				},
				ConflictsWith: []string{
					"function",
				},
			},
			"method": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the method in the contract to call.",
				RequiredWith: []string{
					"artifact",
				},
				ConflictsWith: []string{
					"function",
				},
			},
			"function": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "The typed function to call.",
				ConflictsWith: []string{
					"artifact",
					"method",
				},
			},
			"input": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The inputs of the contract method to call.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gas_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The gas limit of the transaction. This is the maximum amount of gas that can be used to execute the transaction. ",
			},
			"raw_input": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The raw input of the transaction. Alternative to artifact, method and input.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The value of the transaction. This is the amount of wei transferred from the sender to the receiver. ",
			},
			"signer": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The signer of the transaction. This is the private key of the wallet. ",
			},
			"hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hash of the transaction.",
			},
			"block_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The block number at which the transaction is included.",
			},
			"gas_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of gas used to execute the transaction.",
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
