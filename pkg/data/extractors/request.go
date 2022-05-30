package extractors

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data/enrichments"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func RequestData(ctx context.Context, log logr.Logger, r *http.Request, vals chan<- trie.InsertMsg) {

	// TODO This will be the last proxy; look at x-forwarded-for if you want to be better
	if host, port, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		enrichments.EnrichIp(ctx, log, host, trie.PrefixChan(vals, "RemoteAddr"))
		vals <- trie.Insert(trie.Some(host), "RemoteAddr", "Address")
		vals <- trie.Insert(trie.Some(port), "RemoteAddr", "Port")
	}

	vals <- trie.Insert(trie.Some(r.Host), "Headers", "Host") // go's http client promotes this header to here then deletes from Header map
	for k, v := range r.Header {
		vals <- trie.Insert(trie.Some(strings.Join(v, ",")), "Headers", k)
	}
}
