package handlers

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/gorilla/mux"
	"github.com/mt-inside/envbin/pkg/actions"
	"github.com/mt-inside/envbin/pkg/data"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func HandleApi(apiMux *mux.Router) {
	apiMux.NotFoundHandler = notFound(apiMux)

	apiMux.HandleFunc("/exit", exit).Methods("GET")
	apiMux.HandleFunc("/delay", delay).Methods("GET") /* Latency to first byte */
	apiMux.HandleFunc("/bandwidth", bandwidth).Methods("GET") /* Latency between bytes */
	apiMux.HandleFunc("/errorrate", errorrate).Methods("GET") /* Proportion of 500s */
	apiMux.HandleFunc("/allocate", allocate).Methods("GET") /* Allocate (and use) memory */
	apiMux.HandleFunc("/free", free).Methods("GET") /* Free all the extra memory */
	apiMux.HandleFunc("/cpu", cpu).Methods("GET") /* Use CPU at a given rate */
	apiMux.HandleFunc("/liveness", liveness).Methods("GET")
	apiMux.HandleFunc("/readiness", readiness).Methods("GET")
}

func notFound(apiMux *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiMux.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			methods, err := route.GetMethods()
			if err != nil {
				methods = []string{"GET"}
			}
			pathTemplate, err := route.GetPathTemplate()
			if err == nil {
				queriesTemplates, err := route.GetQueriesTemplates()
				if err == nil {
					// TODO: should return JSON? Is there a standard / convention for self-discoverable REST APIs?
					fmt.Fprintf(w, "%s %s?%s\n", methods, pathTemplate, strings.Join(queriesTemplates, ","))
				}
			}
			return nil
		})
	}
}

func exit(w http.ResponseWriter, r *http.Request) {
	rc, err := strconv.ParseInt(r.URL.Query().Get("code"), 0, 32)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		fmt.Fprintf(w, "Exiting %d\n", rc)
		w.(http.Flusher).Flush()
		os.Exit(int(rc))
	}
}

func delay (w http.ResponseWriter, r *http.Request) {
	d, err := strconv.ParseInt(r.URL.Query().Get("value"), 0, 64)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		data.SetDelay(d)
		fmt.Fprintf(w, "Delay set to %v\n", d)
	}
}

func bandwidth(w http.ResponseWriter, r *http.Request) {
	b, err := strconv.ParseInt(r.URL.Query().Get("value"), 0, 64)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		data.SetBandwidth(b)
		fmt.Fprintf(w, "Bandwidth set to %s/s\n", units.BytesSize(float64(b)))
	}
}

func errorrate(w http.ResponseWriter, r *http.Request) {
	e, err := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		data.SetErrorRate(e)
		fmt.Fprintf(w, "Error rate set to %v\n", e)
	}
}

func allocate(w http.ResponseWriter, r *http.Request) {
	a, err := strconv.ParseInt(r.URL.Query().Get("value"), 0, 64)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		fmt.Fprintf(w, "Allocating %s bytes\n", units.BytesSize(float64(a)))
		actions.AllocAndTouch(a)
	}
}

func free(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Freeing\n")
	actions.FreeAllocs()
}

func cpu(w http.ResponseWriter, r *http.Request) {
	c, err := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		data.SetCPUUse(c)
		fmt.Fprintf(w, "CPU usage set to %v\n", c)
	}
}

func liveness(w http.ResponseWriter, r *http.Request) {
	l, err := strconv.ParseBool(r.URL.Query().Get("value"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		data.SetLiveness(l)
		fmt.Fprintf(w, "Liveness check set to %v\n", l)
	}
}

func readiness(w http.ResponseWriter, r *http.Request) {
	ready, err := strconv.ParseBool(r.URL.Query().Get("value"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	} else {
		data.SetReadiness(ready)
		fmt.Fprintf(w, "Readiness check set to %v\n", ready)
	}
}