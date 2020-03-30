package main

// mime type switching, if that's a thing?
// What does curl, browser, etc send?

import (
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	"github.com/mt-inside/envbin/pkg/handlers"
	"github.com/mt-inside/envbin/pkg/middleware"
	"github.com/mt-inside/envbin/pkg/renderers"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	app := cli.App("envbin", "Print environment information, sometimes, badly")
	app.Spec = "[ADDR]"
	addr := app.StringArg("ADDR", ":8080", "Listen address")

	app.Action = func() { envbinMain(addr) }

	app.Run(os.Args)
}

func envbinMain(addr *string) {
	rootMux := mux.NewRouter()
	rootMux.Use(middleware.LoggingMiddleware)

	handlers.HandleApi(rootMux.PathPrefix("/handlers").Subrouter()) //TODO rename our package away from handlers
	handlers.HandleMisc(rootMux)
	handlers.HandleProbes(rootMux)

	rootMux.Path("/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool { return strings.Contains(r.Header.Get("Accept"), "text/html") }).Handler(middleware.MiddlewareStack(renderers.RenderHTML, "text/html"))
	rootMux.Path("/").Headers("Accept", "application/json").Handler(middleware.MiddlewareStack(renderers.RenderJSON, "application/json"))
	rootMux.Path("/").Headers("Accept", "text/yaml").Handler(middleware.MiddlewareStack(renderers.RenderYAML, "text/yaml"))
	rootMux.Path("/").Headers("Accept", "text/x-yaml").Handler(middleware.MiddlewareStack(renderers.RenderYAML, "text/x-yaml"))
	rootMux.Path("/").Headers("Accept", "application/x-yaml").Handler(middleware.MiddlewareStack(renderers.RenderYAML, "application/x-yaml"))
	rootMux.Path("/").Handler(middleware.MiddlewareStack(renderers.RenderText, "text/plain")) // fall through

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, rootMux))

	// TODO: graceful shutdown (lower readiness)
}
