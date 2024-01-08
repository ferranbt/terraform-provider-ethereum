package ethereum

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
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
