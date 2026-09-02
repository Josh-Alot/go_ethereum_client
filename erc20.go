package main

import (
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

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
