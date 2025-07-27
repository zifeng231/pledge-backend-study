package services

import (
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models/request"
	"pledge-backend-study/api/models/response"
	"pledge-backend-study/config"
	"pledge-backend-study/db"
	"pledge-backend-study/log"
	"pledge-backend-study/utils"
)

type UserService struct{}

func NewUser() *UserService {
	return &UserService{}
}

func (s *UserService) Login(req *request.Login, result *response.Login) int {
	log.Logger.Sugar().Info("contractService", req)
	if req.Name == "admin" && req.Password == "password" {
		token, err := utils.CreateToken(req.Name)
		if err != nil {
			log.Logger.Error("CreateToken" + err.Error())
			return statecode.CommonErrServerErr
		}
		result.TokenId = token
		//save to redis
		_ = db.RedisSet(req.Name, "login_ok", config.Config.Jwt.ExpireTime)
		return statecode.CommonSuccess
	} else {
		return statecode.NameOrPasswordErr
	}
}
