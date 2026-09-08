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
	packed, err := parsed.Pack("name")
	if err != nil {
		return nil, fmt.Errorf("pack token name error: %w", err)
	}

	return packed, nil
}

func UnpackName(parsed abi.ABI, returnData []byte) (string, error) {
	if len(returnData) == 0 {
		return "", fmt.Errorf("unpack token name error: expected 96 bytes, got %d bytes, address is not a token contract", len(returnData))
	}

	data, err := parsed.Unpack("name", returnData)
	if err != nil {
		return "", fmt.Errorf("unpack token name error: %w", err)
	}

	name, ok := data[0].(string)
	if !ok {
		return "", fmt.Errorf("unpack token name error: return data invalid")
	}

	return name, nil
}

func PackSymbol(parsed abi.ABI) ([]byte, error) {
	packed, err := parsed.Pack("symbol")
	if err != nil {
		return nil, fmt.Errorf("pack token symbol error: %w", err)
	}

	return packed, nil
}

func UnpackSymbol(parsed abi.ABI, returnData []byte) (string, error) {
	if len(returnData) == 0 {
		return "", fmt.Errorf("unpack token symbol error: expected 96 bytes, got %d bytes, address is not a token contract", len(returnData))
	}

	data, err := parsed.Unpack("symbol", returnData)
	if err != nil {
		return "", fmt.Errorf("unpack token symbol error: %w", err)
	}

	symbol, ok := data[0].(string)
	if !ok {
		return "", fmt.Errorf("unpack token symbol error: return data invalid")
	}

	return symbol, nil
}

func PackDecimals(parsed abi.ABI) ([]byte, error) {
	packed, err := parsed.Pack("decimals")
	if err != nil {
		return nil, fmt.Errorf("pack token decimals error: %w", err)
	}

	return packed, nil
}

func UnpackDecimals(parsed abi.ABI, returnData []byte) (uint8, error) {
	if len(returnData) == 0 {
		return 0, fmt.Errorf("unpack token decimals error: expected 32 bytes, got %d bytes, address is not a token contract", len(returnData))
	}

	data, err := parsed.Unpack("decimals", returnData)
	if err != nil {
		return 0, fmt.Errorf("unpack token decimals error: %w", err)
	}

	decimals, ok := data[0].(uint8)
	if !ok {
		return 0, fmt.Errorf("unpack token decimals error: return data invalid")
	}

	return decimals, nil
}

func PackTotalSupply(parsed abi.ABI) ([]byte, error) {
	packed, err := parsed.Pack("totalSupply")
	if err != nil {
		return nil, fmt.Errorf("pack token total supply error: %w", err)
	}

	return packed, nil
}

func UnpackTotalSupply(parsed abi.ABI, returnData []byte) (*big.Int, error) {
	if len(returnData) == 0 {
		return nil, fmt.Errorf("unpack token total supply error: expected 32 bytes, got %d bytes, address is not a token contract", len(returnData))
	}

	data, err := parsed.Unpack("totalSupply", returnData)
	if err != nil {
		return nil, fmt.Errorf("unpack token total supply error: %w", err)
	}

	totalSupply, ok := data[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unpack token total supply error: return data invalid")
	}

	return totalSupply, nil
}
