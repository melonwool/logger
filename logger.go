package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Logger struct {
	ZapLogger *zap.Logger
	ZapSugar  *zap.SugaredLogger
	FileName  string
	File      io.WriteCloser
	Option
}

type (
	Option struct {
		DateFormat   string
		SignalListen bool
		TickerListen *time.Ticker
		ZapOptions   []zap.Option
	}
	OptFunc func(option *Logger)
)

// NewLogger 初始化一个Logger 包含了zapLogger 和 zapSugar
func NewLogger(fileName string, optFunc ...OptFunc) (*Logger, error) {
	logger := &Logger{
		FileName: fileName,
	}
	for _, fn := range optFunc {
		fn(logger)
	}
	core, err := logger.core()
	logger.newZapLogger(core)
	logger.ZapSugar = logger.ZapLogger.Sugar()
	if logger.SignalListen {
		// 启动一个协程监听信号来重新打开文件写入日志
		go logger.listenSignal()
	}
	if logger.TickerListen != nil {
		go logger.listenTicker()
	}
	return logger, err
}

func (logger *Logger) newZapLogger(core zapcore.Core) {
	logger.ZapLogger = zap.New(core, logger.ZapOptions...)
}

func DateFormat(format string) OptFunc {
	return func(logger *Logger) {
		logger.DateFormat = format
	}
}

func SignalListen() OptFunc {
	return func(logger *Logger) {
		logger.SignalListen = true
	}
}

func TickerListen(ticker *time.Ticker) OptFunc {
	return func(logger *Logger) {
		logger.TickerListen = ticker
	}
}

func ZapOptions(option ...zap.Option) OptFunc {
	return func(logger *Logger) {
		logger.ZapOptions = option
	}
}

// listenSignal 监听SIGUSR1 信号,重新打开文件写入,用来日志切割使用
func (logger *Logger) listenSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	for range c {
		logger.reload()
	}
}

// reload 重新载入Logger 和 Sugar
func (logger *Logger) reload() {
	core, err := logger.core()
	if err != nil {
		log.Println(err)
	}
	logger.newZapLogger(core)
	logger.ZapSugar = logger.ZapLogger.Sugar()
}

// core 获取 zapcore.Core
func (logger *Logger) core() (zapcore.Core, error) {
	var err error
	var core zapcore.Core
	encoderConfig := zap.NewProductionEncoderConfig()
	if logger.DateFormat != "" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(logger.DateFormat)
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	logger.File, err = os.OpenFile(logger.FileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return core, err
	}
	core = zapcore.NewCore(encoder, zapcore.AddSync(logger.File), zapcore.InfoLevel)
	return core, err
}

func (logger *Logger) listenTicker() {
	for range logger.TickerListen.C {
		_, err := os.Stat(logger.FileName)
		if err != nil && os.IsNotExist(err) {
			logger.reload()
		}
	}
}
