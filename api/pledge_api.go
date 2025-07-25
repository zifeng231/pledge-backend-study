package main

import (
	"github.com/gin-gonic/gin"
	"pledge-backend-study/api/models"
	"pledge-backend-study/api/models/kucoin"
	"pledge-backend-study/api/models/middlewares"
	"pledge-backend-study/api/models/ws"
	"pledge-backend-study/api/static"
	"pledge-backend-study/api/validate"
	"pledge-backend-study/config"
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
	// 为 StartServer 的函数，其作用是启动一个服务器，用于监听价格信息并通过 WebSocket 将其发送给所有连接的客户端。
	go ws.StartServer()
	//持续从 KuCoin 交易所获取 PLGR/USDT 交易对的最新行情价格，并将价格保存到全局变量和 Redis 缓存中。
	go kucoin.GetExchangePrice()

	gin.SetMode(gin.ReleaseMode)
	app := gin.Default()
	staticPath := static.GetCurrentAbPathByCaller()
	app.Static("/storage/", staticPath)
	app.Use(middlewares.Cors()) // 「 Cross domain Middleware 」
	routes.initRoutes(app)
	app.Run(":" + config.Config.Env.Port)

}
