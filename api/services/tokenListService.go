package services

import (
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models"
	"pledge-backend-study/api/models/request"
)

type TokenList struct {
}

func NewTokenList() *TokenList {
	return &TokenList{}
}

func (c *TokenList) DebtTokenList(req *request.PoolBaseInfo) (int, []models.TokenInfo) {
	err, res := models.NewTokenInfo().GetTokenInfo(req)
	if err != nil {
		return statecode.CommonErrServerErr, nil
	}
	return statecode.CommonSuccess, res
}

func (c *TokenList) GetTokenList(req *request.PoolBaseInfo) (int, []models.TokenList) {
	err, tokenList := models.NewTokenInfo().GetTokenList(req)
	if err != nil {
		return statecode.CommonErrServerErr, nil
	}
	return statecode.CommonSuccess, tokenList

}
