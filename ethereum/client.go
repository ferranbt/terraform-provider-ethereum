package ethereum

import (
	"context"
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

	from := key.Address()

	chainID, err := c.httpClient.Eth().ChainID()
	if err != nil {
		return ethgo.Hash{}, nil, err
	}

	gasPrice, err := c.httpClient.Eth().GasPrice()
	if err != nil {
		return ethgo.Hash{}, nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	if txn.GasLimit == 0 {
		msg := &ethgo.CallMsg{From: from, To: txn.To, Data: txn.Input, GasPrice: gasPrice, Value: txn.Value}
		txn.GasLimit, err = c.httpClient.Eth().EstimateGas(msg)
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

func (c *client) filterTransactions(ctx context.Context, input filterTransactionInput) (ethgo.Hash, error) {
	mngr := &transactionFilter{
		input: input,
		clt:   c.httpClient.Eth(),
	}
	return mngr.run(ctx)
}

type filterTransactionInput struct {
	From        *ethgo.Address
	To          *ethgo.Address
	IsTransfer  *bool
	StartBlock  uint64
	LimitBlocks *uint64
}

type transactionFilterClient interface {
	GetBlockByNumber(i ethgo.BlockNumber, full bool) (*ethgo.Block, error)
}

type transactionFilter struct {
	input      filterTransactionInput
	clt        transactionFilterClient
	waitPeriod time.Duration
}

func (t *transactionFilter) run(ctx context.Context) (ethgo.Hash, error) {
	if t.waitPeriod == 0 {
		t.waitPeriod = 5 * time.Second
	}

	// get the latest block
	latest, err := t.clt.GetBlockByNumber(ethgo.Latest, false)
	if err != nil {
		return ethgo.Hash{}, err
	}

	startBlock := t.input.StartBlock
	initialBlock := startBlock

	endBlock := latest.Number
	if startBlock > endBlock {
		return ethgo.Hash{}, fmt.Errorf("start block is greater than the latest block")
	}

	for {
		// sync the batch of blocks
		hash, err := t.syncBatch(ctx, initialBlock, startBlock, endBlock)
		if err != nil {
			return ethgo.Hash{}, err
		}
		if hash != nil {
			return *hash, nil
		}

		// validate if the context is still valid or the execution
		// was stopped
		select {
		case <-ctx.Done():
			return ethgo.Hash{}, ctx.Err()
		default:
		}

		startBlock = endBlock

		// wait until the chain has advanced and sync again
		for {
			latest, err := t.clt.GetBlockByNumber(ethgo.Latest, false)
			if err != nil {
				return ethgo.Hash{}, err
			}
			if latest.Number > endBlock {
				endBlock = latest.Number
				break
			}

			// sleep for n seconds and try again or exit if the context
			// was canceled
			select {
			case <-time.After(t.waitPeriod):
			case <-ctx.Done():
				return ethgo.Hash{}, ctx.Err()
			}
		}
	}
}

func validateTxn(txn *ethgo.Transaction, input filterTransactionInput) bool {
	if input.From != nil && txn.From != *input.From {
		return false
	}
	if input.To != nil {
		if txn.To == nil {
			return false
		}
		if *input.To != *txn.To {
			return false
		}
	}
	if input.IsTransfer != nil {
		isTransfer := *input.IsTransfer
		if isTransfer && txn.Value == nil {
			return false
		}
		if !isTransfer && txn.Value != nil {
			return false
		}
	}
	return true
}

func (t *transactionFilter) syncBatch(ctx context.Context, initialBlock, ini, end uint64) (*ethgo.Hash, error) {
	for i := ini; i < end; i++ {
		// validate if we are too far away from the initial block
		if t.input.LimitBlocks != nil {
			limitBlocks := *t.input.LimitBlocks
			if i > initialBlock+limitBlocks {
				return nil, fmt.Errorf("limit blocks exceeded")
			}
		}

		// query the block and validate the transactions
		// exit on the first transaction that matches the input
		num := ethgo.BlockNumber(i)
		block, err := t.clt.GetBlockByNumber(num, true)
		if err != nil {
			return nil, err
		}

		for _, txn := range block.Transactions {
			if validateTxn(txn, t.input) {
				return &txn.Hash, nil
			}
		}

		// validate if the context is still valid or the execution
		// was stopped
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	return nil, nil
}
