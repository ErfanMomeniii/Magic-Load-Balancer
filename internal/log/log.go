package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	Level  zap.AtomicLevel
)

const logLayout = "2020-01-02 15:04:05.000"

// init: Set NewProduction as default logger. Config depend on Logger instance
func init() {
	var err error
	Level = zap.NewAtomicLevel()
	Logger, err = zap.Config{
		Level:             Level,
		Development:       false,
		Encoding:          "json",
		DisableStacktrace: true,
		DisableCaller:     true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			EncodeTime:     zapcore.TimeEncoderOfLayout(logLayout),
			EncodeDuration: zapcore.StringDurationEncoder,

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			NameKey:     "key",
			FunctionKey: zapcore.OmitKey,

			MessageKey: "msg",
			LineEnding: zapcore.DefaultLineEnding,
		},
	}.Build()

	if err != nil {
		panic(err.(any))
	}
}

func CloseLogger() error {
	return Logger.Sync()
}
