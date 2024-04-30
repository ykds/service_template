package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var defaultLogger *Logger

func Info(msg string) {
	defaultLogger.Info(msg)
}

func Infof(msg string) {
	defaultLogger.Infof(msg)
}

func Warn(msg string) {
	defaultLogger.Warnf(msg)
}

func Warnf(msg string) {
	defaultLogger.Warnf(msg)
}

func Error(msg string) {
	defaultLogger.Error(msg)
}

func Errorf(msg string) {
	defaultLogger.Errorf(msg)
}

func Panic(msg string) {
	defaultLogger.Panic(msg)
}

func Panicf(msg string) {
	defaultLogger.Panicf(msg)
}

func Fatal(msg string) {
	defaultLogger.Fatal(msg)
}

func Fatalf(msg string) {
	defaultLogger.Fatalf(msg)
}

type Logger struct {
	*zap.SugaredLogger
	out io.Writer
}

func (l *Logger) Info(msg string) {
	l.SugaredLogger.Info(msg)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.SugaredLogger.Infof(msg, args...)
}

func (l *Logger) Warn(msg string) {
	l.SugaredLogger.Warn(msg)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.SugaredLogger.Warnf(msg, args...)
}

func (l *Logger) Error(msg string) {
	l.SugaredLogger.Error(msg)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.SugaredLogger.Errorf(msg, args...)
}

func (l *Logger) Panic(msg string) {
	l.SugaredLogger.Panic(msg)
}

func (l *Logger) Panicf(msg string, args ...interface{}) {
	l.SugaredLogger.Panicf(msg, args...)
}

func (l *Logger) Fatal(msg string) {
	l.SugaredLogger.Fatal(msg)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.SugaredLogger.Fatalf(msg, args...)
}

func InitLogger() {
	defaultLogger = NewLogger(DefaultOption())
}

func NewLogger(opt Option) *Logger {
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	out := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stderr), zapcore.AddSync(newLumberjack(opt.Lumberjack)))
	core := zapcore.NewCore(enc, out, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
	return &Logger{
		out:           out,
		SugaredLogger: logger.Sugar(),
	}
}

func newLumberjack(opt LumberjackOption) io.Writer {
	return &lumberjack.Logger{
		Filename:   opt.Filename,
		MaxSize:    opt.MaxSize,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
		MaxBackups: opt.MaxBackups,
		LocalTime:  true,
	}
}
