package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models/request"
	"pledge-backend-study/api/models/response"
	"pledge-backend-study/api/services"
	"pledge-backend-study/api/validate"
	"pledge-backend-study/log"
)

type MultiSignPoolController struct {
}

func (c *MultiSignPoolController) SetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.SetMultiSign{}
	log.Logger.Sugar().Info("SetMultiSign req ", req)

	errCode := validate.NewMutiSign().SetMultiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, err := services.NewMutiSign().SetMultiSign(&req)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, nil)
	return
}

func (c *MultiSignPoolController) GetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.GetMultiSign{}
	result := response.MultiSign{}
	log.Logger.Sugar().Info("GetMultiSign req ", nil)

	errCode := validate.NewMutiSign().GetMultiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, err := services.NewMutiSign().GetMultiSign(&result, req.ChainId)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}
