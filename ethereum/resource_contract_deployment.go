package ethereum

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

func ContractDeploymentResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"artifact_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"artifact_contract": {
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
			"hash": {
				Type:     schema.TypeString,
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
	key, err := wallet.NewWalletFromMnemonic("test test test test test test test test test test test junk")
	if err != nil {
		return diag.FromErr(err)
	}

	path := d.Get("artifact_path").(string)
	contract := d.Get("artifact_contract").(string)

	var fullPath string
	err = filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, contract+".json") {
				fullPath = path
			}
			return nil
		})
	if err != nil {
		panic(err)
	}

	rawData, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}

	var artifact *artifact
	if err := json.Unmarshal(rawData, &artifact); err != nil {
		panic(err)
	}

	code, err := hex.DecodeString(artifact.Bytecode[2:])
	if err != nil {
		panic(err)
	}

	if cons := artifact.Abi.Constructor; cons != nil {
		input, ok := d.GetOk("input")
		if !ok {
			return diag.FromErr(fmt.Errorf("no input set but required"))
		}
		xxx, err := cons.Inputs.Encode(input)
		if err != nil {
			panic(err)
		}
		code = append(code, xxx...)
	}

	txn := &ethgo.Transaction{
		Input: code,
	}

	client := meta.(*jsonrpc.Client)
	hash, receipt, err := sendTransaction(client, key, txn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash.String())
	d.Set("hash", hash.String())
	d.Set("gas_used", receipt.GasUsed)
	d.Set("contract_address", receipt.ContractAddress.String())

	return nil
}

type artifact struct {
	Abi      *abi.ABI `json:"abi"`
	Bytecode string   `json:"bytecode"`
}

func resourceContractDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceContractDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
