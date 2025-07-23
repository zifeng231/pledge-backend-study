package log

//日志切割库
import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
)

var Logger *zap.Logger

func init() {
	//zap 不支持文件归档，如果要支持文件按大小或者时间归档，需要使用lumberjack，lumberjack也是zap官方推荐的。
	// https://github.com/uber-go/zap/blob/master/FAQ.md
	hook := lumberjack.Logger{
		Filename:   getCurrentAbPathByCaller() + "/logs/log.log", // 日志文件路径
		MaxSize:    50,                                           // 每个日志文件保存的最大尺寸 单位：M
		MaxAge:     7,                                            // 文件最多保存多少天
		MaxBackups: 20,                                           // 日志文件最多保存多少个备份
		LocalTime:  false,
		Compress:   true, // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",                         // 日志时间字段名
		LevelKey:       "level",                        // 日志级别字段名
		NameKey:        "logger",                       // 日志记录器名称字段名
		CallerKey:      "line",                         // 调用者（代码行号）字段名
		MessageKey:     "msg",                          // 日志消息字段名
		StacktraceKey:  "stacktrace",                   // 堆栈跟踪字段名
		LineEnding:     zapcore.DefaultLineEnding,      // 换行符
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 日志级别小写（如 info、error）
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // 时间格式为 ISO8601（如 2023-05-01T15:04:05.000Z）
		EncodeDuration: zapcore.SecondsDurationEncoder, // 时长以秒为单位
		EncodeCaller:   zapcore.FullCallerEncoder,      // 显示完整文件路径和行号
		EncodeName:     zapcore.FullNameEncoder,        // 显示完整的 logger 名称
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "pledge"))
	// 构造日志
	Logger = zap.New(core, caller, development, filed)
}

// // 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {

	var absPath string
	_, fileName, _, ok := runtime.Caller(0)
	if ok {
		absPath = path.Dir(fileName)
	}
	return absPath
}
