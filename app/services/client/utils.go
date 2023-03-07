package client

import "math/big"

func BigIntToFloat(i *big.Int) *big.Float {
	return new(big.Float).SetInt(i)
}

func StringToFloat(s string) *big.Float {
	f, _ := new(big.Float).SetString(s)
	return f
}

// FilterLowAmount is return value is true, skip delegated value
func FilterLowAmount(val *big.Float) bool {
	if val.Cmp(big.NewFloat(1)) >= 0 {
		return false
	}
	return true
}
