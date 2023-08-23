package middleware

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// Ugly we have to do this and it's not in the library
func LoggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}
