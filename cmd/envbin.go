package main

// mime type switching, if that's a thing?
// What does curl, browser, etc send?

import (
	"fmt"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	"github.com/mt-inside/envbin/pkg/handlers"
	"github.com/mt-inside/envbin/pkg/middleware"
	"github.com/mt-inside/envbin/pkg/renderers"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var (
	Version   string
	GitCommit string
	BuildTime string
)

func main() {
	fmt.Printf("envbin %s git %s, built at %s with %s\n", Version, GitCommit, BuildTime, runtime.Version())
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
	rootMux.Path("/").Headers("Accept", "text/yaml", "Accept", "text/x-yaml", "Accept", "application/x-yaml").Handler(middleware.MiddlewareStack(renderers.RenderYAML, "text/yaml"))
	rootMux.Path("/").Handler(middleware.MiddlewareStack(renderers.RenderText, "text/plain")) // fall through

	log.Printf("Listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, rootMux))

	// TODO: graceful shutdown (lower readiness)
}
