package ethereum

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

type client struct {
	httpClient *jsonrpc.Client
	nonceLock  sync.Mutex
}

func newClient(host string) (*client, error) {
	if host == "" {
		host = defaultHost
	}
	httpClient, err := jsonrpc.NewClient(host)
	if err != nil {
		return nil, err
	}

	clt := &client{
		httpClient: httpClient,
	}
	return clt, nil
}

func (c *client) Http() *jsonrpc.Eth {
	return c.httpClient.Eth()
}

type transaction struct {
	To       *ethgo.Address
	Input    []byte
	Value    *big.Int
	Signer   []byte
	GasLimit uint64
}

func (c *client) sendTransaction(txn *transaction) (ethgo.Hash, *ethgo.Receipt, error) {
	if txn.Signer == nil {
		return ethgo.Hash{}, nil, fmt.Errorf("signer not found")
	}
	key, err := wallet.NewWalletFromPrivKey(txn.Signer)
	if err != nil {
		return ethgo.Hash{}, nil, err
	}

	chainID, err := c.httpClient.Eth().ChainID()
	if err != nil {
		return ethgo.Hash{}, nil, err
	}

	gasPrice, err := c.httpClient.Eth().GasPrice()
	if err != nil {
		return ethgo.Hash{}, nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	if txn.GasLimit == 0 {
		txn.GasLimit, err = c.httpClient.Eth().EstimateGas(&ethgo.CallMsg{From: key.Address(), To: txn.To, Data: txn.Input, Value: txn.Value})
		if err != nil {
			return ethgo.Hash{}, nil, fmt.Errorf("gas estimation failed: %v", err)
		}
	}

	c.nonceLock.Lock()
	nonce, err := c.httpClient.Eth().GetNonce(key.Address(), ethgo.Pending)
	if err != nil {
		defer c.nonceLock.Unlock()
		return ethgo.Hash{}, nil, err
	}

	ethTxn := &ethgo.Transaction{
		Input:    txn.Input,
		To:       txn.To,
		Value:    txn.Value,
		Gas:      txn.GasLimit,
		GasPrice: gasPrice,
		Nonce:    nonce,
	}

	signer := wallet.NewEIP155Signer(chainID.Uint64())
	ethTxn, err = signer.SignTx(ethTxn, key)
	if err != nil {
		defer c.nonceLock.Unlock()
		return ethgo.Hash{}, nil, err
	}

	raw, _ := ethTxn.MarshalRLPTo(nil)
	hash, err := c.httpClient.Eth().SendRawTransaction(raw)
	if err != nil {
		defer c.nonceLock.Unlock()
		return ethgo.Hash{}, nil, err
	}

	c.nonceLock.Unlock()

	tt := time.NewTimer(15 * time.Second)
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			receipt, _ := c.httpClient.Eth().GetTransactionReceipt(hash)
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
