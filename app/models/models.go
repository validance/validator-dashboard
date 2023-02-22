package models

type AddressStatus struct {
	Chain string
	// enum type of postgres ('a41', 'a41ventures', 'grant', 'b2b', 'b2c', 'unknown')
	Label   string
	Address string
	// enum type of postgres ('new', 'existing', 'leave', 'return')
	Type string
}

type Delegation struct {
	Address   string
	Validator string
	Chain     string
	Amount    float64
}

type ValidatorIncome struct {
	Chain     string
	Validator string
	// reward from self delegated token
	Reward float64
	// commission from non-self delegated token
	Commission float64
}

type GrantReward struct {
	Chain     string
	Validator string
	Reward    float64
}

type DelegationSummary struct {
	// 기존 전일 위임량
	YesterdayDelegationAmount *DelegationSummaryLabel
	// 기존 전일 위임액
	YesterdayDelegationValue *DelegationSummaryLabel
	// 기존 당일 추가 위임량
	TodayExistingIncreasedDelegationAmount *DelegationSummaryLabel
	// 기존 당일 추가 위임액
	TodayExistingIncreasedDelegationValue *DelegationSummaryLabel
	// 신규 당일 추가 위임량
	TodayNewIncreasedDelegationAmount *DelegationSummaryLabel
	// 신규 당일 추가 위임액
	TodayNewIncreasedDelegationValue *DelegationSummaryLabel
	// 재방문 당일 위임량
	TodayReturnIncreasedDelegationAmount *DelegationSummaryLabel
	// 재방문 당일 위임액
	TodayReturnIncreasedDelegationValue *DelegationSummaryLabel
	// 기존 당일 줄어든 위임량
	TodayExistingDecreasedDelegationAmount *DelegationSummaryLabel
	// 기존 당일 줄어든 위임액
	TodayExistingDecreasedDelegationValue *DelegationSummaryLabel
	// 이탈 당일 이탈한 위임량
	TodayLeftDecreasedDelegationAmount *DelegationSummaryLabel
	// 이탈 당일 이탈한 위임액
	TodayLeftDecreasedDelegationValue *DelegationSummaryLabel
}

type DelegationSummaryLabel struct {
	B2B     float64
	B2C     float64
	Unknown float64
}
