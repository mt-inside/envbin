package middleware

import (
	"github.com/mt-inside/envbin/pkg/data"
	"net/http"
	"time"
)

func latencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(data.GetDelay()) * time.Second)
		next.ServeHTTP(w, r)
	})
}