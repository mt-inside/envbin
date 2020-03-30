package middleware

import (
	"github.com/gorilla/handlers"
	"net/http"
)

// Ugly we have to do this and it's not in the library
// TODO: doesn't work?
func recoveryMiddleware(next http.Handler) http.Handler {
	return handlers.RecoveryHandler()(next)
}
