package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/middleware"
	"github.com/mt-inside/envbin/pkg/renderers"
	"github.com/mt-inside/envbin/pkg/util"
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
		logr = util.GetLogger(*devMode)
	}

	app.Command("serve", "Serve environment information as various mimetypes", func(cmd *cli.Cmd) {
		var (
			addr = cmd.StringArg("ADDR", ":8080", "Listen address")
		)
		cmd.Spec = "[ADDR]"
		cmd.Action = func() { serve(logr, addr) }
	})

	app.Command("oneshot", "Print information to stdout and exit", func(cmd *cli.Cmd) {
		cmd.Action = func() { oneshot(logr) }
	})

	app.Run(os.Args)
}

func serve(log logr.Logger, addr *string) {
	rootMux := mux.NewRouter()
	rootMux.Use(middleware.LoggingMiddleware)

	rootMux.Path("/health").HandlerFunc(healthHandler)
	rootMux.Path("/ready").HandlerFunc(healthHandler)

	rootMux.Path("/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return strings.Contains(r.Header.Get("Accept"), "text/html")
	}).Handler(middleware.MiddlewareStack(renderers.RenderHTML, "text/html"))
	rootMux.Path("/").Headers("Accept", "application/json").Handler(middleware.MiddlewareStack(renderers.RenderJSON, "application/json"))
	rootMux.Path("/").Headers("Accept", "text/yaml", "Accept", "text/x-yaml", "Accept", "application/x-yaml").Handler(middleware.MiddlewareStack(renderers.RenderYAML, "text/yaml"))
	rootMux.Path("/").Handler(middleware.MiddlewareStack(renderers.RenderText, "text/plain")) // fall through

	log.Info("Listening", "addr", *addr)
	http.ListenAndServe(*addr, rootMux)

	// TODO: graceful shutdown (lower readiness - combine with badpod first)
}

func oneshot(log logr.Logger) {
	renderers.RenderTTY()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

	log.Printf("Served health ok")
}
