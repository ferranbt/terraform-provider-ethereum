package ethereum

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/wallet"
)

func TestUtils_DecodeArtifact(t *testing.T) {
	_, err := resolveContract("./fixtures:foundry")
	require.NoError(t, err)

	_, err = resolveContract("./fixtures:hardhat")
	require.NoError(t, err)
}

func TestUtils_DecodeInputs(t *testing.T) {
	cases := []struct {
		in  []interface{}
		out []interface{}
	}{
		{
			in: []interface{}{
				"1",
			},
			out: []interface{}{
				"1",
			},
		},
		{
			// one nested level
			in: []interface{}{
				`{"a": 1, "b": 2}`,
				"3",
			},
			out: []interface{}{
				map[string]interface{}{
					"a": float64(1),
					"b": float64(2),
				},
				"3",
			},
		},
		{
			// two nested levels
			in: []interface{}{
				`{"a": 1, "b": "{\"c\": 1}"}`,
			},
			out: []interface{}{
				map[string]interface{}{
					"a": float64(1),
					"b": map[string]interface{}{
						"c": float64(1),
					},
				},
			},
		},
	}

	for _, c := range cases {
		out, err := decodeInputs(c.in)
		require.NoError(t, err)
		require.Equal(t, c.out, out)
	}
}

func TestParseEther(t *testing.T) {
	cases := []struct {
		str string
		res *big.Int
	}{
		{
			"1",
			big.NewInt(1),
		},
		{
			"0x1",
			big.NewInt(1),
		},
		{
			"1 ether",
			big.NewInt(1000000000000000000),
		},
		{
			"0.4 ether",
			big.NewInt(400000000000000022),
		},
		{
			"1 gwei",
			big.NewInt(1000000000),
		},
	}

	for _, c := range cases {
		val, err := parseEtherValue(c.str)
		require.NoError(t, err)
		require.Equal(t, val, c.res)
	}
}

func TestIsOwner(t *testing.T) {

	buf, _ := hex.DecodeString("6fd0b749794e67a63154bdbf406b18ebb3b372111472402ed5032a2651704a71")
	key, err := wallet.NewWalletFromPrivKey(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(key.Address())

	clt, err := newClient("http://localhost:8449")
	if err != nil {
		panic(err)
	}

	/*
		x, err := abi.NewABIFromList([]string{
			"function isChainOwner(address addr) external view returns (bool)",
		})
		if err != nil {
			panic(err)
		}

		owner := ethgo.HexToAddress("0xBee947Aec820389c1560C87BD96e2723BaD05b61")
		inputs, err := x.Methods["isChainOwner"].Encode([]interface{}{key.Address()})
		if err != nil {
			panic(err)
		}

		to := ethgo.HexToAddress("0x0000000000000000000000000000000000000070")

		res, err := clt.httpClient.Eth().Call(&ethgo.CallMsg{
			From: owner,
			To:   &to,
			Data: inputs,
		}, ethgo.Latest)
		if err != nil {
			panic(err)
		}

		bufHex, _ := hex.DecodeString(res[2:])

		resMap, err := x.Methods["isChainOwner"].Decode(bufHex)
		if err != nil {
			panic(err)
		}

		fmt.Println(resMap)
	*/

	fmt.Println(clt.httpClient.Eth().GetBalance(key.Address(), ethgo.Latest))
}
