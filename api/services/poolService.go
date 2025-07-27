package services

import (
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models"
	"pledge-backend-study/log"
)

type poolService struct {
}

func NewPoolService() *poolService {
	return &poolService{}
}

func NewPool() *poolService {
	return &poolService{}
}

func (s *poolService) PoolBaseInfo(chainId int, result *[]models.PoolBaseInfoRes) int {

	err := models.NewPoolBases().PoolBaseInfo(chainId, result)
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr
	}
	return statecode.CommonSuccess
}

func (s *poolService) PoolDataInfo(chainId int, result *[]models.PoolDataInfoRes) int {

	err := models.NewPoolData().PoolDataInfo(chainId, result)
	if err != nil {
		log.Logger.Error(err.Error())
		return statecode.CommonErrServerErr
	}
	return statecode.CommonSuccess
}
