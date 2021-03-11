package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/middleware"
	"github.com/mt-inside/envbin/pkg/renderers"
	"github.com/mt-inside/go-usvc"
)

type serveCmd struct {
	Addr string `short:"a" long:"addr" description:"listen address" default:":8080"`
}

var serveOpts serveCmd

func init() {
	if _, err := flagParser.AddCommand(
		"serve",
		"Serves on HTTP",
		"Serves envbin info over HTTP. Various mimetypes can be requested",
		&serveOpts,
	); err != nil {
		panic(err)
	}
}

func (cmd *serveCmd) Execute(args []string) error {
	log := usvc.GetLogger(mainOpts.DevMode)
	log.Info(data.RenderBuildData())

	rootMux := mux.NewRouter()
	rootMux.Use(middleware.LoggingMiddleware)

	rootMux.Path("/health").HandlerFunc(healthHandler)
	rootMux.Path("/ready").HandlerFunc(healthHandler)

	rootMux.Path("/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		return strings.Contains(r.Header.Get("Accept"), "text/html")
	}).Handler(middleware.MiddlewareStack(log, renderers.RenderHTML))
	rootMux.Path("/").Headers("Accept", "application/json").Handler(middleware.MiddlewareStack(log, renderers.RenderJSON))
	rootMux.Path("/").Headers("Accept", "text/yaml", "Accept", "text/x-yaml", "Accept", "application/x-yaml").Handler(middleware.MiddlewareStack(log, renderers.RenderYAML))
	rootMux.Path("/").Handler(middleware.MiddlewareStack(log, renderers.RenderText)) // fall through

	log.Info("Listening", "addr", cmd.Addr)
	http.ListenAndServe(cmd.Addr, rootMux)

	// TODO: graceful shutdown (lower readiness - combine with badpod first)

	return nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))

	usvc.Global.Info("Served health ok")
}
