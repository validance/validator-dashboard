package client

import "math/big"

func BigIntToFloat(i *big.Int) *big.Float {
	return new(big.Float).SetInt(i)
}
