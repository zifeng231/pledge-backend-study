package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io"
	"pledge-backend-study/api/common/statecode"
	"pledge-backend-study/api/models/request"
)

// 定义一个空的结构体
// 在 Go 中，结构体可以没有字段，这在某些场景下很有用（比如方法接收者、占位符等）。
type PoolBaseInfo struct {
}

// &表示取地址
// *表示指向的指针，*PoolBaseInfo表示指向的PoolBaseInfo结构体的指针
// *有两个用途 用途一：声明一个指针类型  用途二：解引用 —— 通过指针访问值
func NewPoolBaseInfo() *PoolBaseInfo {
	return &PoolBaseInfo{}
}

func (v *PoolBaseInfo) PoolBaseInfo(c *gin.Context, req *request.PoolBaseInfo) int {

	err := c.ShouldBind(req)
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "ChainId" {
				return statecode.ChainIdEmpty
			}
		}
		return statecode.CommonErrServerErr
	}
	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr
	}
	return statecode.CommonSuccess

}
