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

type TokenSummary struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply *big.Int
}

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

func TokenInfo(client *Client, token string) (*TokenSummary, error) {
	parsed, err := ParseERC20ABI()
	if err != nil {
		return nil, fmt.Errorf("parse erc20 abi contract error: %w", err)
	}

	namePacked, err := packName(parsed)
	if err != nil {
		return nil, fmt.Errorf("pack name error: %w", err)
	}

	symbolPacked, err := packSymbol(parsed)
	if err != nil {
		return nil, fmt.Errorf("pack symbol error: %w", err)
	}

	decimalsPacked, err := packDecimals(parsed)
	if err != nil {
		return nil, fmt.Errorf("pack decimals error: %w", err)
	}

	totalSupplyPacked, err := packTotalSupply(parsed)
	if err != nil {
		return nil, fmt.Errorf("pack total supply error: %w", err)
	}

	returnNameData, err := client.Call(token, namePacked)
	if err != nil {
		return nil, fmt.Errorf("call name error: %w", err)
	}

	name, err := unpackName(parsed, returnNameData)
	if err != nil {
		return nil, fmt.Errorf("unpack name error: %w", err)
	}

	returnSymbolData, err := client.Call(token, symbolPacked)
	if err != nil {
		return nil, fmt.Errorf("call symbol error: %w", err)
	}

	symbol, err := unpackSymbol(parsed, returnSymbolData)
	if err != nil {
		return nil, fmt.Errorf("unpack symbol error: %w", err)
	}

	returnDecimalsData, err := client.Call(token, decimalsPacked)
	if err != nil {
		return nil, fmt.Errorf("call decimals error: %w", err)
	}

	decimals, err := unpackDecimals(parsed, returnDecimalsData)
	if err != nil {
		return nil, fmt.Errorf("unpack decimals error: %w", err)
	}

	returnTotalSupplyData, err := client.Call(token, totalSupplyPacked)
	if err != nil {
		return nil, fmt.Errorf("call total supply error: %w", err)
	}

	totalSupply, err := unpackTotalSupply(parsed, returnTotalSupplyData)
	if err != nil {
		return nil, fmt.Errorf("unpack total supply error: %w", err)
	}

	summary := TokenSummary{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
	}

	return &summary, nil
}

func packBalanceOf(parsed abi.ABI, address string) ([]byte, error) {
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

func unpackBalanceOf(parsed abi.ABI, returnData []byte) (*big.Int, error) {
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

func packName(parsed abi.ABI) ([]byte, error) {
	packed, err := parsed.Pack("name")
	if err != nil {
		return nil, fmt.Errorf("pack token name error: %w", err)
	}

	return packed, nil
}

func unpackName(parsed abi.ABI, returnData []byte) (string, error) {
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

func packSymbol(parsed abi.ABI) ([]byte, error) {
	packed, err := parsed.Pack("symbol")
	if err != nil {
		return nil, fmt.Errorf("pack token symbol error: %w", err)
	}

	return packed, nil
}

func unpackSymbol(parsed abi.ABI, returnData []byte) (string, error) {
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

func packDecimals(parsed abi.ABI) ([]byte, error) {
	packed, err := parsed.Pack("decimals")
	if err != nil {
		return nil, fmt.Errorf("pack token decimals error: %w", err)
	}

	return packed, nil
}

func unpackDecimals(parsed abi.ABI, returnData []byte) (uint8, error) {
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

func packTotalSupply(parsed abi.ABI) ([]byte, error) {
	packed, err := parsed.Pack("totalSupply")
	if err != nil {
		return nil, fmt.Errorf("pack token total supply error: %w", err)
	}

	return packed, nil
}

func unpackTotalSupply(parsed abi.ABI, returnData []byte) (*big.Int, error) {
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
