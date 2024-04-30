package logger

type Option struct {
	Lumberjack LumberjackOption `json:"lumberjack" yaml:"lumberjack"`
}

type LumberjackOption struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
	MaxBackups int    `json:"max_backups"`
}

func DefaultOption() Option {
	return Option{
		Lumberjack: LumberjackOption{
			Filename:   "err.log",
			MaxSize:    10,
			MaxAge:     7,
			Compress:   false,
			MaxBackups: 3,
		},
	}
}
