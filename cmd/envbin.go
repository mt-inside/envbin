package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	"github.com/mattn/go-isatty"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/middleware"
	"github.com/mt-inside/envbin/pkg/renderers"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	log.Println(data.RenderBuildData())

	app := cli.App("envbin", "Print environment information, sometimes, badly")
	app.Spec = "[--dev]"

	var (
		devMode = app.BoolOpt("dev", false, "Developer mode")

		logr logr.Logger
	)

	app.Before = func() {
		logr = getLogger(*devMode)
	}

	app.Command("serve", "Serve environment information as various mimetypes", func(cmd *cli.Cmd) {
		var (
			addr = cmd.StringArg("ADDR", ":8080", "Listen address")
		)
		cmd.Spec = "[ADDR]"
		cmd.Action = func() { serve(addr) }
	})

	app.Command("oneshot", "Print information to stdout and exit", func(cmd *cli.Cmd) {
		cmd.Action = func() { oneshot(logr) }
	})

	app.Run(os.Args)
}

func oneshot(log logr.Logger) {
	renderers.RenderTTY()
}

func getLogger(devMode bool) logr.Logger {
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

	return zr
}

type ZaprWriter struct{ log logr.Logger }

func (w ZaprWriter) Write(data []byte) (n int, err error) {
	w.log.Info(string(data))
	return len(data), nil
}

func serve(addr *string) {
	rootMux := mux.NewRouter()
	rootMux.Use(middleware.LoggingMiddleware)

	rootMux.Path("/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return strings.Contains(r.Header.Get("Accept"), "text/html")
	}).Handler(middleware.MiddlewareStack(renderers.RenderHTML, "text/html"))
	rootMux.Path("/").Headers("Accept", "application/json").Handler(middleware.MiddlewareStack(renderers.RenderJSON, "application/json"))
	rootMux.Path("/").Headers("Accept", "text/yaml", "Accept", "text/x-yaml", "Accept", "application/x-yaml").Handler(middleware.MiddlewareStack(renderers.RenderYAML, "text/yaml"))
	rootMux.Path("/").Handler(middleware.MiddlewareStack(renderers.RenderText, "text/plain")) // fall through

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, rootMux))

	// TODO: graceful shutdown (lower readiness)
}
