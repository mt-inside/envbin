package middleware

import (
	"net/http"
)

func MiddlewareStack(next func(r *http.Request) []byte, mime string) http.Handler {
	return recoveryMiddleware( // let it crash
		proxyHeaders( // Sets header x-envbin-proxy-chain any and all forwarded addresses
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", mime)

				bs := next(r)

				// could just execute the template straight into the writer, but we're gonna merge back with badpod soon
				w.Write(bs)
			}),
		),
	)
}
