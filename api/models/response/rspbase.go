package response

import "github.com/gin-gonic/gin"

type Gin struct {
	Res *gin.Context
}

type Page struct {
	Code  int         `json:"code"`
	Msg   string      `json:"message"`
	Total int         `json:"total"`
	Data  interface{} `json:"data"`
}

// 统一分页格式
func (g *Gin) ResponsePages(c *gin.Context, code int, total int, data interface{}) {
	lang := statecode.LangZh

}
