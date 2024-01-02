package ethereum

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/wallet"
)

var defTestSigner = []byte{}

func init() {
	key, _ := wallet.NewWalletFromMnemonic("test test test test test test test test test test test junk")
	defTestSigner, _ = key.MarshallPrivateKey()
}

func TestClient_SendTransaction_Simple(t *testing.T) {
	// Simple transaction
	testAccPreCheck(t)

	clt, _ := newClient("")

	acct, _ := wallet.GenerateKey()
	target := acct.Address()

	sendBalance := big.NewInt(100000)

	txn := &transaction{
		To:     &target,
		Value:  sendBalance,
		Signer: defTestSigner,
	}

	_, receipt, err := clt.sendTransaction(txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
	require.NoError(t, err)

	balance, err := clt.Http().GetBalance(target, ethgo.Latest)
	require.NoError(t, err)
	require.Equal(t, sendBalance, balance)
}

func TestClient_SendTransaction_Concurrent(t *testing.T) {
	// the client should handle the concurrent nonce
	testAccPreCheck(t)

	clt, _ := newClient("")

	acct, _ := wallet.GenerateKey()
	target := acct.Address()

	sendBalance := big.NewInt(100000)

	var wg sync.WaitGroup

	num := 3
	for i := 0; i < num; i++ {
		wg.Add(1)

		go func() {
			txn := &transaction{
				To:     &target,
				Value:  sendBalance,
				Signer: defTestSigner,
			}

			_, receipt, err := clt.sendTransaction(txn)
			require.NoError(t, err)
			require.Equal(t, receipt.Status, uint64(1))
			require.NoError(t, err)

			defer wg.Done()
		}()
	}

	wg.Wait()

	balance, err := clt.Http().GetBalance(target, ethgo.Latest)
	require.NoError(t, err)
	require.Equal(t, balance, new(big.Int).Mul(sendBalance, new(big.Int).SetInt64(int64(num))))
}

func TestTransactionFilter_RunSimple(t *testing.T) {
	blockNums := uint64(10)

	mock := &mockTransactionFilterClient{
		ch: make(chan uint64),
	}
	mock.move(blockNums)

	mngr := &transactionFilter{
		input: filterTransactionInput{
			StartBlock: 1,
		},
		clt:        mock,
		waitPeriod: 1 * time.Second,
	}

	go mngr.run(context.Background())

	// wait the initial batch of 'blockNums'
	count := uint64(1)
	for ; count <= blockNums; count++ {
		select {
		case num := <-mock.ch:
			require.Equal(t, num, count)
		case <-time.After(2 * time.Second):
			t.Fatal("timeout")
		}
	}

	// move the chain and wait for the result
	mock.move(5)

	for ; count < blockNums+5; count++ {
		select {
		case num := <-mock.ch:
			require.Equal(t, num, count)
		case <-time.After(2 * time.Second):
			t.Fatal("timeout")
		}
	}
}

func uintPtr(n uint64) *uint64 {
	return &n
}

func TestTransactionFilter_LimitBatch(t *testing.T) {
	mock := &mockTransactionFilterClient{}
	mock.move(1000)

	mngr := &transactionFilter{
		input: filterTransactionInput{
			StartBlock:  10,
			LimitBlocks: uintPtr(100),
		},
		clt: mock,
	}
	mngr.run(context.Background())

	// the latest queried block should be 110 (start + limit)
	require.Equal(t, uint64(110), mock.latestQueried)
}

func TestTransactionFilter_LimitWatch(t *testing.T) {
	mock := &mockTransactionFilterClient{}
	mock.move(20)

	mngr := &transactionFilter{
		clt:        mock,
		waitPeriod: 1 * time.Second,
		input: filterTransactionInput{
			StartBlock:  10,
			LimitBlocks: uintPtr(20),
		},
	}

	doneCh := make(chan struct{})
	go func() {
		mngr.run(context.Background())
		close(doneCh)
	}()

	// mine enough blocks for the process to reach the limit
	// and wait for it to finish
	mock.move(20)

	<-doneCh
	require.Equal(t, uint64(30), mock.latestQueried)
}

type mockTransactionFilterClient struct {
	lock          sync.Mutex
	blocks        uint64
	ch            chan uint64
	latestQueried uint64
}

func (m *mockTransactionFilterClient) move(n uint64) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.blocks += n
}

func (m *mockTransactionFilterClient) GetBlockByNumber(i ethgo.BlockNumber, full bool) (*ethgo.Block, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if i == ethgo.Latest {
		return &ethgo.Block{
			Number: m.blocks,
		}, nil
	}

	num := uint64(i)
	if num > m.blocks {
		return nil, fmt.Errorf("block not found")
	}

	m.latestQueried = num
	if m.ch != nil {
		m.ch <- num
	}

	return &ethgo.Block{
		Number: num,
	}, nil
}

func TestTransactionFilter_ValidateInput(t *testing.T) {
	txn1 := &ethgo.Transaction{
		From:  ethgo.Address{0x1},
		To:    &ethgo.Address{0x2},
		Value: big.NewInt(100),
	}

	trueVal := true
	falseVal := false

	var cases = []struct {
		txn   *ethgo.Transaction
		input filterTransactionInput
		valid bool
	}{
		{
			// to is invalid
			txn: txn1,
			input: filterTransactionInput{
				From: &ethgo.Address{0x1},
			},
			valid: true,
		},
		{
			// from is invalid
			txn: txn1,
			input: filterTransactionInput{
				From: &ethgo.Address{0x2},
			},
			valid: false,
		},
		{
			// to is invalid
			txn: txn1,
			input: filterTransactionInput{
				To: &ethgo.Address{0x1},
			},
			valid: false,
		},
		{
			// to is not set
			txn: &ethgo.Transaction{},
			input: filterTransactionInput{
				To: &ethgo.Address{0x1},
			},
			valid: false,
		},
		{
			// to is valid
			txn: txn1,
			input: filterTransactionInput{
				To: &ethgo.Address{0x2},
			},
			valid: true,
		},
		{
			// value is set
			txn: txn1,
			input: filterTransactionInput{
				IsTransfer: &trueVal,
			},
			valid: true,
		},
		{
			// value is not-set
			txn: &ethgo.Transaction{},
			input: filterTransactionInput{
				IsTransfer: &trueVal,
			},
			valid: false,
		},
		{
			// value should not be set
			txn: &ethgo.Transaction{},
			input: filterTransactionInput{
				IsTransfer: &falseVal,
			},
			valid: true,
		},
	}

	for _, c := range cases {
		require.Equal(t, c.valid, validateTxn(c.txn, c.input))
	}
}

func TestTransactionFilterXXX(t *testing.T) {
	clt, _ := jsonrpc.NewClient("http://localhost:8449")

	trueX := true

	mngr := &transactionFilter{
		input: filterTransactionInput{
			StartBlock:  0,
			LimitBlocks: uintPtr(10),
			IsTransfer:  &trueX,
		},
		clt:        clt.Eth(),
		waitPeriod: 1 * time.Second,
	}
	mngr.run(context.Background())
}
