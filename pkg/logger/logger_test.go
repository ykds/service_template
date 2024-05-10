package logger

import (
	"io"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	lj := NewLumberjack(LumberjackOption{})
	opt := Option{
		Lumberjack: LumberjackOption{},
		Output:     nil,
		ErrOutput:  nil,
	}
	opt.Output = []io.Writer{lj}
	opt.ErrOutput = []io.Writer{lj}
	opt.Output = append(opt.Output, os.Stdout)
	opt.ErrOutput = append(opt.Output, os.Stderr)
	log := InitLogger(opt)

	log.Info("this is info")
	log.Warn("this is warn")
	log.Error("this is error")
}
