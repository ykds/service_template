package logger

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *Logger

func init() {
	InitLogger(Option{})
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

func Panic(msg string) {
	logger.Panic(msg)
}

func Panicf(msg string, args ...interface{}) {
	logger.Panicf(msg, args)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args)
}

type Option struct {
	Lumberjack LumberjackOption `json:"lumberjack" yaml:"lumberjack"`
	Output     []io.Writer
	ErrOutput  []io.Writer
}

type LumberjackOption struct {
	Filename   string `json:"filename" yaml:"filename"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	Compress   bool   `json:"compress" yaml:"compress"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
}

type Logger struct {
	*zap.SugaredLogger
	output io.Writer
}

func InitLogger(opt Option) *Logger {
	if len(opt.Output) == 0 {
		opt.Output = []io.Writer{os.Stdout}
	}
	if len(opt.ErrOutput) == 0 {
		opt.ErrOutput = []io.Writer{os.Stderr}
	}
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	enc := zapcore.NewJSONEncoder(config)
	syncers := make([]zapcore.WriteSyncer, len(opt.Output))
	errSyncers := make([]zapcore.WriteSyncer, len(opt.ErrOutput))
	for i, out := range opt.Output {
		syncers[i] = zapcore.AddSync(out)
	}
	for i, out := range opt.ErrOutput {
		errSyncers[i] = zapcore.AddSync(out)
	}
	syncer := zapcore.NewMultiWriteSyncer(syncers...)
	core := zapcore.NewCore(enc, syncer, zapcore.InfoLevel)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.ErrorOutput(zapcore.NewMultiWriteSyncer(errSyncers...)))
	zap.ReplaceGlobals(zapLogger)
	logger = &Logger{
		output:        zapcore.NewMultiWriteSyncer(syncers...),
		SugaredLogger: zapLogger.Sugar(),
	}
	return logger
}

func (l *Logger) GetOutput() io.Writer {
	return l.output
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

func defaultLumberjackOption() LumberjackOption {
	return LumberjackOption{
		Filename:   "api.log",
		MaxSize:    5,
		MaxAge:     3,
		Compress:   false,
		MaxBackups: 0,
	}
}

func NewLumberjack(opt LumberjackOption) io.Writer {
	defaultOpt := defaultLumberjackOption()
	if opt.Filename != "" {
		defaultOpt.Filename = opt.Filename
	}
	if opt.MaxSize != 0 {
		defaultOpt.MaxSize = opt.MaxSize
	}
	if opt.MaxAge != 0 {
		defaultOpt.MaxAge = opt.MaxAge
	}
	defaultOpt.Compress = opt.Compress
	if opt.MaxBackups != 0 {
		defaultOpt.MaxBackups = opt.MaxBackups
	}
	return &lumberjack.Logger{
		Filename:   defaultOpt.Filename,
		MaxSize:    defaultOpt.MaxSize,
		MaxAge:     defaultOpt.MaxAge,
		Compress:   defaultOpt.Compress,
		MaxBackups: defaultOpt.MaxBackups,
		LocalTime:  true,
	}
}
