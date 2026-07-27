package main

import (
	"fmt"
	"math/big"
)

var divisor = big.NewInt(1e18)

func FormatWeiToEther(balance *big.Int) string {
	remain := new(big.Int)
	quo := new(big.Int)

	quo.QuoRem(balance, divisor, remain)

	return fmt.Sprintf("%s.%018s", quo, remain)
}
