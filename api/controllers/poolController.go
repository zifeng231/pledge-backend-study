package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models"
	"pledge-backend-study/api/models/request"
	"pledge-backend-study/api/models/response"
	"pledge-backend-study/api/services"
	"pledge-backend-study/api/validate"
	"pledge-backend-study/config"
	"regexp"
	"strings"
	"time"
)

type PoolController struct {
}

func (c *PoolController) PoolBaseInfo(ctx *gin.Context) {
	//相当于定义了一个结构体 封装返回参数
	res := response.Gin{Res: ctx}
	//接受请求参数
	req := request.PoolBaseInfo{}

	//定义返回参数data
	var result []models.PoolBaseInfoRes
	//验证请求参数
	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	errCode = services.NewPoolService().PoolBaseInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

func (c *PoolController) PoolDataInfo(context *gin.Context) {
	r := response.Gin{Res: context}
	req := request.PoolBaseInfo{}

	var result []models.PoolDataInfoRes

	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(context, &req)
	if errCode != statecode.CommonSuccess {
		r.Response(context, errCode, nil)
		return
	}

	errCode = services.NewPool().PoolDataInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		r.Response(context, errCode, nil)
		return
	}

	r.Response(context, statecode.CommonSuccess, result)
	return
}

func (c *PoolController) DebtTokenList(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.PoolBaseInfo{}

	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, result := services.NewTokenList().DebtTokenList(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

func (c *PoolController) TokenList(context *gin.Context) {

	req := request.PoolBaseInfo{}

	result := response.TokenList{}
	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(context, &req)
	if errCode != statecode.CommonSuccess {
		context.JSON(200, map[string]string{
			"error": "chainId error",
		})
		return
	}

	errCode, data := services.NewTokenList().GetTokenList(&req)
	if errCode != statecode.CommonSuccess {
		context.JSON(200, map[string]string{
			"error": "chainId error",
		})
		return
	}
	var BaseUrl = c.GetBaseUrl()
	result.Name = "Pledge Token List"
	result.LogoURI = BaseUrl + "storage/img/Pledge-project-logo.png"
	result.Timestamp = time.Now()
	result.Version = response.Version{
		Major: 2,
		Minor: 16,
		Patch: 12,
	}
	for _, v := range data {
		result.Tokens = append(result.Tokens, response.Token{
			Name:     v.Symbol,
			Symbol:   v.Symbol,
			Decimals: v.Decimals,
			Address:  v.Token,
			ChainID:  v.ChainId,
			LogoURI:  v.Logo,
		})
	}

	context.JSON(200, result)
	return

}

func (c *PoolController) GetBaseUrl() string {

	domainName := config.Config.Env.DomainName
	domainNameSlice := strings.Split(domainName, "")
	pattern := "\\d+"
	isNumber, _ := regexp.MatchString(pattern, domainNameSlice[0])
	if isNumber {
		return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + ":" + config.Config.Env.Port + "/"
	}
	return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + "/"
}

func (c *PoolController) Search(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.Search{}
	result := response.Search{}

	errCode := validate.NewSearch().Search(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, count, pools := services.NewSearch().Search(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	result.Rows = pools
	result.Count = count
	res.Response(ctx, statecode.CommonSuccess, result)
	return
}
