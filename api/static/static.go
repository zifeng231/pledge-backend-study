package static

import (
	"path"
	"runtime"
)

// 获取当前 Go 源文件所在的目录路径
func GetCurrentAbPathByCaller() string {
	var absPath string
	_, fileName, _, ok := runtime.Caller(0)
	if ok {
		absPath = path.Dir(fileName)
	}
	return absPath
}
