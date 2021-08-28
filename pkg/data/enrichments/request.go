package enrichments

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func GetDataWithRequest(ctx context.Context, log logr.Logger, r *http.Request) *Trie {
	t := data.GetData(ctx, log)
	getRequestData(ctx, log, t, r)

	return t
}

func getRequestData(ctx context.Context, log logr.Logger, t *Trie, r *http.Request) {
	t.Insert(Some(r.RemoteAddr), "Request", "RemoteAddr") // TODO This will be the last proxy; look at x-forwarded-for if you want to be better

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		reqIpInfo, err := EnrichIpRendered(ctx, log, host)
		if err != nil {
			log.Error(err, "Can't get IP info", "ip", host)
			if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
				t.Insert(Timeout(time.Second), "Request", "RemoteAddr", "Info") // FIXME: duration
			} else {
				t.Insert(Error(err), "Request", "RemoteAddr", "Info")
			}
		} else {
			t.Insert(Some(reqIpInfo), "Request", "RemoteAddr", "Info")
		}
	}

	t.Insert(Some(r.UserAgent()), "Request", "UserAgent")

	t.Insert(Some(r.Header.Get("X-Envbin-Proxy-Chain")), "Request", "Proxies")
}
