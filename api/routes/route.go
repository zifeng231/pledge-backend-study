package routes

import (
	"github.com/gin-gonic/gin"
	"pledge-backend-study/api/controllers"
	"pledge-backend-study/config"
)

func InitRoute(e *gin.Engine) *gin.Engine {
	e.Group("/api/v" + config.Config.Env.Version)
	poolController := controllers.PoolController{}

}
