package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
)

func MiddlewareStack(
	log logr.Logger,
	next func(log logr.Logger, w http.ResponseWriter, r *http.Request, d *data.Trie) []byte,
) http.Handler {
	return recoveryMiddleware( // let it crash
		proxyHeaders( // Sets header x-envbin-proxy-chain any and all forwarded addresses
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()

				d := data.GetDataWithRequest(ctx, log, r)

				bs := next(log, w, r, d)

				// could just execute the template straight into the writer, but we're gonna merge back with badpod soon
				_, err := w.Write(bs)
				if err != nil {
					panic(err)
				}
			}),
		),
	)
}
