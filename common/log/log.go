package zaplog

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// const (
// 	logDir  = "./logs/"
// 	logFile = "log.txt"
// )

func LoggerInit() error {

	execPath, _ := os.Executable()
	logDir := filepath.Join(filepath.Dir(execPath), "logs")
	logFile := filepath.Join(logDir, "log.txt")
	// 检查目录是否存在
	if _, err := os.Stat(logDir); !os.IsNotExist(err) {
		// // 目录已存在
		// fmt.Println("目录已存在：")
	} else {
		// 目录不存在，创建目录
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Failed to create the log directory:%v\n", err)
			return nil
		}
		// fmt.Println("目录已创建：")
	}
	fileConfig := &lumberjack.Logger{
		Filename:   logFile, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    5,       //文件大小限制,单位MB
		MaxBackups: 5,       //最大保留日志文件数量
		MaxAge:     30,      //日志文件保留天数
		Compress:   false,   //是否压缩处理
	}
	//zapcore.AddSync 函数通常用于将一个 io.Writer 转换为 zapcore.WriteSyncer 接口的实现。
	//zapcore.WriteSyncer 是 zap 包中的一个接口，它扩展了 io.Writer 接口，增加了一个 Sync 方法，
	//该方法用于确保所有已写入的数据都被正确地刷新到它们的最终目的地
	FileWriteSyncer := zapcore.AddSync(fileConfig)
	stdioWriteSyncer := zapcore.AddSync(os.Stdout)

	//设置日志编码器
	EncoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewJSONEncoder(EncoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, FileWriteSyncer, zap.InfoLevel),
		zapcore.NewCore(encoder, stdioWriteSyncer, zap.ErrorLevel),
	)
	//初始化实例
	logger := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)
	logger.Info("log init success")
	return nil
}
