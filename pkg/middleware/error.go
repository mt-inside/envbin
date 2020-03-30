package middleware

import (
	"fmt"
	"github.com/mt-inside/envbin/pkg/data"
	"math/rand"
	"net/http"
)

func errorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if data.GetErrorRate() < rand.Float64() {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, "error")
		}
	})
}
