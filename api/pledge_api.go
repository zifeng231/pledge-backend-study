package main

import (
	"pledge-backend-study/api/models"
	"pledge-backend-study/api/models/ws"
	"pledge-backend-study/api/validate"
	"pledge-backend-study/db"
)

func main() {
	//先初始化mysql
	db.InitMysql()
	//初始化redis
	db.InitRedis()
	//初始化表
	models.InitTable()
	//表单验证之类的
	validate.BindingValidator()
	// websocket server
	go ws.StartServer()
}
