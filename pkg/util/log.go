package util

import (
	"log"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/mattn/go-isatty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var GlobalLog logr.Logger

func GetLogger(devMode bool) logr.Logger {
	var zapLog *zap.Logger
	var err error

	if isatty.IsTerminal(os.Stdout.Fd()) || devMode {
		c := zap.NewDevelopmentConfig()
		c.EncoderConfig.EncodeCaller = nil
		c.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("15:04:05"))
		}
		zapLog, err = c.Build()
	} else {
		zapLog, err = zap.NewProduction()
	}
	if err != nil {
		panic(err.Error())
	}

	zr := zapr.NewLogger(zapLog)

	if devMode {
		zr.Info("Logging in dev mode; remove --dev flag for structured json output")
	}

	log.SetFlags(0) // don't add date and timestamps to the message, as the zapr writer will do that
	log.SetOutput(ZaprWriter{zr.WithValues("source", "go log")})

	GlobalLog = zr.WithName("GLOBAL!")

	return zr
}

type ZaprWriter struct{ log logr.Logger }

func (w ZaprWriter) Write(data []byte) (n int, err error) {
	w.log.Info(string(data))
	return len(data), nil
}
