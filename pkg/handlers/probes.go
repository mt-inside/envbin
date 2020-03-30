package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mt-inside/envbin/pkg/data"
	"net/http"
)

func HandleProbes(mux *mux.Router) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if data.GetLiveness() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ok")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error")
		}
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if data.GetReadiness() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "ok")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error")
		}

	})
}
