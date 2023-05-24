package ethereum

import (
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo"
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
