package logger

import (
	"backend-rms/core/config"
	"backend-rms/util/email"
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ErrorKey = "error"
)

var (
	logger       Logger
	accessLogger Logger
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	ErrorIf(err error)
	WithContext(ctx context.Context) Logger
	Close() error
}

type sugaredLogger struct {
	*zap.SugaredLogger
}

func (s *sugaredLogger) ErrorIf(err error) {
	if err != nil {
		s.WithOptions(zap.AddCallerSkip(1)).Errorw("Error Occurred", ErrorKey, err)
	}
}

func (s *sugaredLogger) WithContext(ctx context.Context) Logger {
	log := s.WithOptions(zap.AddCallerSkip(-1))
	if loginManagerID, ok := ctx.Value("login_manager_id").(int64); ok {
		log = log.With("login_manager_id", loginManagerID)
	}
	if requestID, ok := ctx.Value("request_id").(string); ok {
		log = log.With("request_id", requestID)
	}
	return &sugaredLogger{log}
}

func (s *sugaredLogger) Close() error {
	return s.Sync()
}

func Init(logConf config.LogConf) {
	writeSyncer := getLogWriter(logConf)

	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(logConf.Level)); err != nil {
		panic("init logger failed, please check the level config item for log")
	}

	core := zapcore.NewCore(getEncoder(logConf), writeSyncer, logLevel)
	log := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	logger = &sugaredLogger{log.Sugar()}

	accessWriteSyncer := getAccessLogWriter(logConf)
	accessCore := zapcore.NewCore(getEncoder(logConf), accessWriteSyncer, logLevel)
	accessLog := zap.New(accessCore,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	accessLogger = &sugaredLogger{accessLog.Sugar()}
}

func getEncoder(logConf config.LogConf) zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,                       // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,                     // 日志级别名称大写，如 ERROR、INFO
		EncodeTime:     zapcore.TimeEncoderOfLayout(logConf.TimeFormat), // 时间格式，例如 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder,                  // 执行时间，以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,                      // Caller 短格式，如：types/converter.go:17，长格式为绝对路径
	}

	if logConf.Stdout {
		// 终端输出的关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(logConf config.LogConf) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logConf.FileName,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackup,
		MaxAge:     logConf.MaxAge,
		Compress:   logConf.Compress,
	}

	if logConf.Stdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}

	return zapcore.AddSync(lumberJackLogger)
}

func getAccessLogWriter(logConf config.LogConf) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logConf.AccessFolder,
		MaxSize:    logConf.MaxSize,
		MaxBackups: logConf.MaxBackup,
		MaxAge:     logConf.MaxAge,
		Compress:   logConf.Compress,
	}

	if logConf.Stdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}

	return zapcore.AddSync(lumberJackLogger)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
	email.SendException(fmt.Sprint(args...))
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
	errString := fmt.Sprintf(template, args...)
	email.SendException(errString)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
	errString := fmt.Sprintf(msg, keysAndValues...)
	email.SendException(errString)
}

func ErrorIf(err error) {
	logger.ErrorIf(err)
	email.SendException(fmt.Sprintf("%v", err))
}

func WithContext(ctx context.Context) Logger {
	return logger.WithContext(ctx)
}

func Close() error {
	return logger.Close()
}

func AccessLogger() Logger {
	return accessLogger
}
