package main

import (
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
