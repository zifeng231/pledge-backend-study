package validate

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"pledge-backend-study/db"
	"regexp"
	"strings"
	"sync"
)

// 定义一个读写锁
//var lock sync.RWMutex

var lock = sync.RWMutex{}

// 定义一个全局map
//var mapString = map[string][]string{}

// 推荐使用这种方式来定义，可以制定容量
var mapString = make(map[string][]string)

var splitParamsRegexString = `'[^']*'|\S+`
var splitParamsRegex = regexp.MustCompile(splitParamsRegexString)

//var m map[string][]string 一个没有初始化的map 只能读 不能写

// 这段代码的作用是注册自定义的字段校验规则，用于表单/请求参数的校验。
// 常见于基于 Gin 框架的 Go 项目，配合 go-playground/validator 使用。
func BindingValidator() {

	//binding.Validator.Engine()：获取 Gin 框架当前使用的校验引擎。
	//(*validator.Validate)：类型断言，判断校验引擎是不是 go-playground/validator 的实例。
	//RegisterValidation("IsPassword", IsPassword)：注册一个名为 "IsPassword" 的自定义校验方法，方法名为 IsPassword（函数）。
	//其他类似，都是注册自定义的校验规则。
	//  type User struct {
	//      Password string `validate:"IsPassword"`
	//      Email    string `validate:"IsEmail"`
	//  }
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("IsPassword", IsPassword)       //判断是否为合法密码
		_ = v.RegisterValidation("IsPhoneNumber", IsPhoneNumber) //检查手机号码字段是否合法
		_ = v.RegisterValidation("IsEmail", IsEmail)             //检查邮箱字段是否合法
		//_ = v.RegisterValidation("CheckUserNicknameLength", CheckUserNicknameLength) //检查用户昵称长度是否合法
		//_ = v.RegisterValidation("CheckUserAccount", CheckUserAccount)               //检查用户账号是否合法
		_ = v.RegisterValidation("OnlyOne", OnlyOne) //字段唯一性约束
	}
}

// 字段唯一性的约束 类似于下面这
//
//	type UserCreateRequest struct {
//		Username string `json:"username" validate:"required,onlyone=users,username"`
//		Email    string `json:"email" validate:"required,email,onlyone=users,email"`
//		Password string `json:"password" validate:"required,min=6"`
//	}
func OnlyOne(fl validator.FieldLevel) bool {
	//先去解析参数
	params := fl.Param()
	vals := parse(params)
	if len(vals) <= 0 {
		panic("参数错误")
	}
	tableName := vals[0]
	fieldName := vals[1]

	var data dataStruct
	sqlStr := fmt.Sprintf("`%s`=?", fieldName)
	db.Mysql.Table(tableName).Select("count(*)").Where(sqlStr, fl.Field().Interface().(string)).Scan(&data.DataCount)
	if data.DataCount > 0 {
		return false
	}
	return true
}

type dataStruct struct {
	DataCount int //这个结构体用来保存查询到的记录条数
}

// 字符串转数组
func parse(params string) []string {
	//先使用读锁
	lock.RLock()
	//在 Go 中，从 map 中取值时可以返回两个值：
	strings1, ok := mapString[params]
	lock.RUnlock()
	if !ok {
		lock.Lock()
		allString := splitParamsRegex.FindAllString(params, -1)
		for i := 0; i < len(allString); i++ {
			allString[i] = strings.Replace(allString[i], "'", "", -1)
		}
		mapString[params] = allString
		lock.Unlock()
	}
	return strings1

}

func IsPassword(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9!@#￥%^&*]{6,20}$`, fl.Field().Interface().(string)); matched {
			return matched
		}
		return false
	}
	return false
}

// IsPhoneNumber 检查手机号码字段是否合法
func IsPhoneNumber(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		if isOk, _ := regexp.MatchString(`^1[0-9]{10}$`, fl.Field().Interface().(string)); isOk {
			return isOk
		}
	}
	return false
}

// IsEmail 检查手机号码字段是否合法
func IsEmail(fl validator.FieldLevel) bool {
	if fl.Field().Interface().(string) != "" {
		if isOk, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, fl.Field().Interface().(string)); isOk {
			return isOk
		}
	}
	return false
}
