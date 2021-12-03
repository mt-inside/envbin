package extractors

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data/enrichments"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func RequestData(ctx context.Context, log logr.Logger, r *http.Request, vals chan<- InsertMsg) {

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		enrichments.EnrichIp(ctx, log, host, PrefixChan(vals, "RemoteAddr"))
	}
	vals <- Insert(Some(r.RemoteAddr), "RemoteAddr", "Address") // TODO This will be the last proxy; look at x-forwarded-for if you want to be better

	vals <- Insert(Some(r.Host), "Headers", "Host") // go's http client promotes this header to here then deletes from Header map
	for k, v := range r.Header {
		vals <- Insert(Some(strings.Join(v, ",")), "Headers", k)
	}
}
