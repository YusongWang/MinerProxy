package utils

import (
	"os"

	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func init() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)
	var writes = []zapcore.WriteSyncer{}

	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stderr.Fd()) {
		// REWORK 移除这个大括号。将注释打开。目前全部输出为文件日志方便调试
	}

	path := "./logs/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	hook := lumberjack.Logger{
		Filename:   path + "MinerProxy.log", // 日志文件路径
		MaxSize:    128,                     // 每个日志文件保存的大小 单位:M
		MaxAge:     7,                       // 文件最多保存多少天
		MaxBackups: 30,                      // 日志文件最多保存多少个备份
		Compress:   true,                    // 是否压缩
	}

	writes = append(writes, zapcore.AddSync(&hook))
	// } else {
	// 	writes = append(writes, zapcore.AddSync(os.Stdout))
	// }

	core := zapcore.NewCore(
		//zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	//field := zap.Fields(zap.String("app", app_name))
	// 构造日志
	Logger = zap.New(core, caller, development)
}
