package main

import (
	_ "embed"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//go:embed erc20.abi.json
var erc20ABIJSON string

func BalanceOfCalldata(address string) []byte {
	fnKeccak256 := crypto.Keccak256([]byte("balanceOf(address)"))
	fnFourBytes := slices.Clone(fnKeccak256[:4])

	addrBytes := common.HexToAddress(address).Bytes()
	addrBytes = common.LeftPadBytes(addrBytes, 32)

	calldata := append(fnFourBytes, addrBytes...)

	return calldata
}

func VerifyBalanceOfReturnData(returnData []byte) (*big.Int, error) {
	if len(returnData) != 32 {
		return nil, fmt.Errorf("return data error: expected 32 bytes, got %d bytes, address is not a token contract", len(returnData))
	}
	return new(big.Int).SetBytes(returnData), nil
}

func ParseERC20ABI() (abi.ABI, error) {
	reader := strings.NewReader(erc20ABIJSON)

	parsed, err := abi.JSON(reader)
	if err != nil {
		return abi.ABI{}, fmt.Errorf("parse erc20 abi error: %w", err)
	}

	return parsed, nil
}
func PackBalanceOf(parsed abi.ABI, address string) ([]byte, error) {
	isAddr := common.IsHexAddress(address)
	if !isAddr {
		return nil, fmt.Errorf("address %s: invalid address", address)
	}

	addr := common.HexToAddress(address)

	packed, err := parsed.Pack("balanceOf", addr)
	if err != nil {
		return nil, fmt.Errorf("pack balance of error: %w", err)
	}

	return packed, nil
}

func UnpackBalanceOf(parsed abi.ABI, returnData []byte) (*big.Int, error) {
	if len(returnData) == 0 {
		return nil, fmt.Errorf("unpack balance of error: expected 32 bytes, got %d bytes, address is not a token contract", len(returnData))
	}

	data, err := parsed.Unpack("balanceOf", returnData)
	if err != nil {
		return nil, fmt.Errorf("unpack balance of error: %w", err)
	}

	// the data variable is an []interface{} value, it has the same effect as any on javascript
	// it means that in compilation time, the program doesnt know what is this return
	// because the balanceOf function on a ERC20 Token returns a uint256 and we know
	// the value type, we make a cast of the value to *big.Int, if something
	// goes wrong on the casting, it will return an error
	balance, ok := data[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unpack balance of error: return data invalid")
	}

	return balance, nil
}

func PackName(parsed abi.ABI) ([]byte, error) {
	return []byte{}, nil
}

func UnpackName(parsed abi.ABI, returnData []byte) (string, error) {
	return "", nil
}

func PackSymbol(parsed abi.ABI) ([]byte, error) {
	return []byte{}, nil
}

func UnpackSymbol(parsed abi.ABI, returnData []byte) (string, error) {
	return "", nil
}

func PackDecimals(parsed abi.ABI) ([]byte, error) {
	return []byte{}, nil
}

func UnpackDecimals(parsed abi.ABI, returnData []byte) (uint8, error) {
	return 0, nil
}

func PackTotalSupply(parsed abi.ABI) ([]byte, error) {
	return []byte{}, nil
}

func UnpackTotalSupply(parsed abi.ABI, returnData []byte) (*big.Int, error) {
	return &big.Int{}, nil
}
