package middleware

import (
	"net/http"

	"github.com/mt-inside/envbin/pkg/data"
)

func MiddlewareStack(next func(map[string]string) []byte, mime string) http.Handler {
	return recoveryMiddleware(
		proxyHeaders( // Sets header x-envbin-proxy-chain any and all forwarded addresses
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", mime)

				bs := next(data.GetData(r)) // TODO: refactor. This really shouldn't be here

				// Templates can be executed straight into writers, so we could pump the template into the httpResponseWriter. Problem is, it only flushes on the boundaries into and out of {{}} template substitutions, which makes the output sporadic. So we dump into a string and write that one byte at a time.
				for i := 0; i < len(bs); i++ {
					w.Write(bs[i : i+1])
				}
			}),
		),
	)
}
