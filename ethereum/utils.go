package ethereum

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/umbracle/ethgo/abi"
)

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

	return nil, fmt.Errorf("unknown artifact format: %s", string(data))
}

// resolveContract resolves a contract abi specification
// from a 'fullPath' reference that includeds both the path
// and the contract name as fullPath:name.
func resolveContract(fullPath string) (*artifact, error) {
	parts := strings.Split(fullPath, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("there are no two parts")
	}

	relPath, contractName := parts[0], parts[1]
	var contractPath string

	err := filepath.Walk(relPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, contractName+".json") {
				contractPath = path
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	if contractPath == "" {
		return nil, fmt.Errorf("contract not found")
	}
	data, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, err
	}

	artifact, err := decodeArtifact(data)
	if err != nil {
		return nil, err
	}
	return artifact, nil
}

func decodeInputs(input interface{}) (interface{}, error) {
	var err error

	switch obj := input.(type) {
	case string:
		// check if the string is a nested json
		if strings.HasPrefix(obj, "{") {
			var inputMap map[string]interface{}
			if err := json.Unmarshal([]byte(obj), &inputMap); err != nil {
				return nil, err
			}
			return decodeInputs(inputMap)
		}

	case map[string]interface{}:
		// clean each of the elements in case there is another json
		newObj := map[string]interface{}{}
		for k, v := range obj {
			if k == "chainConfig" {
				// TODO: fix this case
				newObj[k] = v
				continue
			}
			if newObj[k], err = decodeInputs(v); err != nil {
				return nil, err
			}
		}
		return newObj, nil

	case []interface{}:
		newSlice := make([]interface{}, len(obj))
		for indx, val := range obj {
			if newSlice[indx], err = decodeInputs(val); err != nil {
				return nil, err
			}
		}
		return newSlice, nil
	}

	return input, nil
}
func convert(val float64, decimals int64) *big.Int {
	// Convert the float64 to a big.Float
	v := new(big.Float).SetFloat64(val)

<<<<<<< HEAD
=======
func convert(val float64, decimals int64) *big.Int {
	// Convert the float64 to a big.Float
	v := new(big.Float).SetFloat64(val)

>>>>>>> main
	// Calculate the multiplier (10^decimals)
	exp := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil))

	// Multiply the value by the multiplier
	result := new(big.Float).Mul(v, exp)

	// Convert the result to a *big.Int
	intResult, _ := result.Int(nil)
	return intResult
}

var unitSuffixesFn = map[string]func(i float64) *big.Int{
	" gwei": func(i float64) *big.Int {
		return convert(i, 9)
	},
	" ether": func(i float64) *big.Int {
		return convert(i, 18)
	},
}

func parseEtherValue(i string) (*big.Int, error) {
	for p, fn := range unitSuffixesFn {
		if strings.HasSuffix(i, p) {
			// try to decode the value as float and apply the function
			num, err := strconv.ParseFloat(strings.TrimSuffix(i, p), 64)
			if err != nil {
				return nil, err
			}
			if num < 0 {
				return nil, fmt.Errorf("cannot be lower than zero")
			}
			return fn(num), nil
		}
	}

	var ok bool

	// try to decode directly as big.Int
	num := new(big.Int)
	if strings.HasPrefix(i, "0x") {
		num, ok = num.SetString(strings.TrimPrefix(i, "0x"), 16)
		if !ok {
			return nil, fmt.Errorf("failed to decode hex number")
		}
	} else {
		num, ok = num.SetString(i, 10)
		if !ok {
			return nil, fmt.Errorf("failed to decode number")
		}
	}
	return num, nil
}
