package models

import "math/big"

type AddressStatus struct {
	Chain string
	// enum type of postgres ('a41', 'a41ventures', 'grant', 'b2b', 'b2c')
	Label   string
	Address string
	// enum type of postgres ('new', 'existing', 'leave', 'return')
	Type string
}

type Delegation struct {
	Address   string
	Validator string
	Chain     string
	Amount    *big.Int
}

type ValidatorIncome struct {
	Chain     string
	Validator string
	// reward from self delegated token
	Reward *big.Int
	// commission from non-self delegated token
	Commission *big.Int
}

type Reward struct {
	Chain     string
	Validator string
	Value     *big.Int
}
