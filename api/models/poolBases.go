package models

type PoolBases struct {
	Id                     int    `json:"-" gorm:"column:id;primaryKey"`
	PoolID                 int    `json:"pool_id" gorm:"column:pool_id;"`
	AutoLiquidateThreshold string `json:"autoLiquidateThreshold" gorm:"column:auto_liquidata_threshold;"`
	BorrowSupply           string `json:"borrowSupply" gorm:"column:borrow_supply;"`
	BorrowToken            string `json:"borrowToken" gorm:"column:pool_id;"`
	BorrowTokenInfo        string `json:"borrowTokenInfo" gorm:"column:borrow_token_info;"`
	EndTime                string `json:"endTime" gorm:"end_time;"`
	InterestRate           string `json:"interestRate" gorm:"column:interest_rate;"`
	JpCoin                 string `json:"jpCoin" gorm:"column:jp_coin;"`
	LendSupply             string `json:"lendSupply" gorm:"column:lend_supply;"`
	LendToken              string `json:"lendToken" gorm:"column:lend_token;"`
	LendTokenInfo          string `json:"lendTokenInfo" gorm:"column:lend_token_info;"`
	MartgageRate           string `json:"martgageRate" gorm:"column:martgage_rate;"`
	MaxSupply              string `json:"maxSupply" gorm:"column:max_supply;"`
	SettleTime             string `json:"settleTime" gorm:"column:settle_time;"`
	SpCoin                 string `json:"spCoin" gorm:"column:sp_coin;"`
	State                  string `json:"state" gorm:"column:state;"`
}
