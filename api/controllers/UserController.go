package controllers

import (
	"github.com/gin-gonic/gin"
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models/request"
	"pledge-backend-study/api/models/response"
	"pledge-backend-study/api/services"
	"pledge-backend-study/api/validate"
	"pledge-backend-study/db"
)

type UserController struct {
}

func (c *UserController) Login(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.Login{}
	result := response.Login{}

	errCode := validate.NewUser().Login(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode = services.NewUser().Login(&req, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

func (c *UserController) Logout(ctx *gin.Context) {
	res := response.Gin{Res: ctx}

	usernameIntf, _ := ctx.Get("username")

	//delete username in redis
	_, _ = db.RedisDelete(usernameIntf.(string))

	res.Response(ctx, statecode.CommonSuccess, nil)
	return
}
