package data

import (
	"context"
	"net"
	"net/http"
	"net/url"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/enrichments"
)

func getRequestData(ctx context.Context, log logr.Logger, t *Trie, r *http.Request) {
	t.Insert(Some{r.RemoteAddr}, "Request", "RemoteAddr") // TODO This will be the last proxy; look at x-forwarded-for if you want to be better

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		reqIpInfo, err := enrichments.EnrichIpRendered(ctx, log, host)
		if err != nil {
			log.Error(err, "Can't get IP info", "ip", host)
			if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
				t.Insert(Some{"Timeout"}, "Request", "RemoteAddr", "Info")
			} else {
				t.Insert(Some{"Error"}, "Request", "RemoteAddr", "Info")
			}
		} else {
			t.Insert(Some{reqIpInfo}, "Request", "RemoteAddr", "Info")
		}
	}

	t.Insert(Some{r.UserAgent()}, "Request", "UserAgent")

	t.Insert(Some{r.Header.Get("X-Envbin-Proxy-Chain")}, "Request", "Proxies")
}
