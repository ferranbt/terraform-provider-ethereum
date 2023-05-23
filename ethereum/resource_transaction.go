package ethereum

import (
	"context"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

func TransactionResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"to": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"input": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},
			"hash": {
				Type:     schema.TypeString,
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
	key, err := wallet.NewWalletFromMnemonic("test test test test test test test test test test test junk")
	if err != nil {
		return diag.FromErr(err)
	}

	txn := &ethgo.Transaction{}

	if val, ok := d.GetOk("to"); ok {
		addr := ethgo.HexToAddress(val.(string))
		txn.To = &addr
	}
	if val, ok := d.GetOk("value"); ok {
		txn.Value = new(big.Int).SetInt64(int64(val.(float64)))
	}
	if val, ok := d.GetOk("input"); ok {
		buf, err := hex.DecodeString(val.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		txn.Input = buf
	}

	client := meta.(*jsonrpc.Client)
	hash, receipt, err := sendTransaction(client, key, txn)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hash.String())
	d.Set("hash", hash.String())
	d.Set("gas_used", receipt.GasUsed)

	return nil
}

func resourceTransactionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceTransactionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func sendTransaction(client *jsonrpc.Client, key *wallet.Key, txn *ethgo.Transaction) (ethgo.Hash, *ethgo.Receipt, error) {
	chainID, err := client.Eth().ChainID()
	if err != nil {
		return ethgo.Hash{}, nil, err
	}

	gasPrice, err := client.Eth().GasPrice()
	if err != nil {
		return ethgo.Hash{}, nil, err
	}
	txn.GasPrice = gasPrice

	gasLimit, err := client.Eth().EstimateGas(&ethgo.CallMsg{To: txn.To, Data: txn.Input, Value: txn.Value})
	if err != nil {
		return ethgo.Hash{}, nil, err
	}
	txn.Gas = gasLimit

	nonce, err := client.Eth().GetNonce(key.Address(), ethgo.Latest)
	if err != nil {
		return ethgo.Hash{}, nil, err
	}
	txn.Nonce = nonce

	signer := wallet.NewEIP155Signer(chainID.Uint64())
	txn, err = signer.SignTx(txn, key)
	if err != nil {
		return ethgo.Hash{}, nil, err
	}

	raw, _ := txn.MarshalRLPTo(nil)
	hash, err := client.Eth().SendRawTransaction(raw)
	if err != nil {
		return ethgo.Hash{}, nil, err
	}

	tt := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			receipt, _ := client.Eth().GetTransactionReceipt(hash)
			if receipt != nil {

				if receipt.Status != 1 {
					panic("not success")
				}

				return hash, receipt, nil
			}
		case <-tt.C:
			panic("not found")
		}
	}
}
