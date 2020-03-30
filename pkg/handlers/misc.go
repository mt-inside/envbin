package handlers

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleMisc(mux *mux.Router) {
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Can't read body: %v", err)
		}
		w.Write(bs)
	})
}
