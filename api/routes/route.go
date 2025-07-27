package routes

import (
	"github.com/gin-gonic/gin"
	"pledge-backend-study/api/controllers"
	"pledge-backend-study/api/models/middlewares"
	"pledge-backend-study/config"
)

func InitRoute(e *gin.Engine) *gin.Engine {
	group := e.Group("/api/v" + config.Config.Env.Version)
	poolController := controllers.PoolController{}
	//获取所有的池子的基本信息
	group.GET("/poolBaseInfo", poolController.PoolBaseInfo) //pool base information

	//获取所有池子的数据信息
	group.GET("/poolDataInfo", poolController.PoolDataInfo) //pool data information

	group.GET("/token", poolController.TokenList) //pool token information

	group.POST("/pool/debtTokenList", middlewares.CheckToken(), poolController.DebtTokenList) //pool debtTokenList
	group.POST("/pool/search", middlewares.CheckToken(), poolController.Search)

	// plgr-usdt price
	priceController := controllers.PriceController{}
	group.GET("/price", priceController.NewPrice) //new price on ku-coin-exchange

	// pledge-defi admin backend
	multiSignPoolController := controllers.MultiSignPoolController{}
	group.POST("/pool/setMultiSign", middlewares.CheckToken(), multiSignPoolController.SetMultiSign) //multi-sign set
	group.POST("/pool/getMultiSign", middlewares.CheckToken(), multiSignPoolController.GetMultiSign) //multi-sign get

	userController := controllers.UserController{}
	group.POST("/user/login", userController.Login)                             // login
	group.POST("/user/logout", middlewares.CheckToken(), userController.Logout) // logout
	return e

}
