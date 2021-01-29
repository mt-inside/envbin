package data

import (
	"net"
	"net/http"

	"github.com/mt-inside/envbin/pkg/enrichments"
)

func getRequestData(r *http.Request) map[string]string {
	if r == nil {
		return nil
	}

	data := map[string]string{}

	data["RequestIp"] = r.RemoteAddr // This will be the last proxy; look at x-forwarded-for if you want to be better
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		data["RequestIpEnrich"] = enrichments.EnrichIpRendered(host)
	}
	data["UserAgent"] = r.UserAgent()
	data["ProxyChain"] = r.Header.Get(http.CanonicalHeaderKey("X-Envbin-Proxy-Chain"))

	return data
}
