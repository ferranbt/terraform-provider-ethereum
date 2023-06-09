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
	"github.com/umbracle/ethgo/abi"
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
			"signer": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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
	signer, err := hex.DecodeString(d.Get("signer").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	txn := &transaction{
		Signer: signer,
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
		return diag.FromErr(err)
	}

	rawData, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return diag.FromErr(err)
	}

	artifact, err := decodeArtifact(rawData)
	if err != nil {
		return diag.FromErr(err)
	}

	code, err := hex.DecodeString(artifact.Bytecode.Object[2:])
	if err != nil {
		return diag.FromErr(err)
	}

	if cons := artifact.Abi.Constructor; cons != nil {
		input, ok := d.GetOk("input")
		if !ok {
			return diag.FromErr(fmt.Errorf("no input set but required"))
		}

		// convert any json input into a map struct
		inputList := []interface{}{}
		for _, val := range input.([]interface{}) {
			valStr := val.(string)

			if strings.HasPrefix(valStr, "{") {
				var inputMap map[string]interface{}
				if err := json.Unmarshal([]byte(valStr), &inputMap); err != nil {
					return diag.FromErr(err)
				}
				inputList = append(inputList, inputMap)
			} else {
				inputList = append(inputList, valStr)
			}
		}

		xxx, err := cons.Inputs.Encode(inputList)
		if err != nil {
			return diag.FromErr(err)
		}
		code = append(code, xxx...)
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

	return nil
}

type artifact struct {
	Abi      *abi.ABI `json:"abi"`
	Bytecode bytecode `json:"bytecode"`
}

type bytecode struct {
	Object string `json:"object"`
}

type artifactHardhat struct {
	Abi      *abi.ABI `json:"abi"`
	Bytecode string
}

func decodeArtifact(data []byte) (*artifact, error) {
	// first try to decode with the foundry artifact format
	var fArtifact *artifact
	if err := json.Unmarshal(data, &fArtifact); err == nil {
		return fArtifact, nil
	}

	// try to decode with hardhat artifact format
	var hArtifact artifactHardhat
	if err := json.Unmarshal(data, &hArtifact); err == nil {
		return &artifact{Abi: hArtifact.Abi, Bytecode: bytecode{Object: hArtifact.Bytecode}}, nil
	}

	return nil, fmt.Errorf("unknown artifact format")
}

func resourceContractDeploymentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceContractDeploymentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
