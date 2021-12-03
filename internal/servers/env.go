package servers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/extractors"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func GetEnv(log logr.Logger, r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Have to register this as a plugin so that it's run in the pool with the rest, cause it needs to be in parallel with them cause it does i/o.
		data.RegisterPlugin(
			func(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
				extractors.RequestData(ctx, log, c.Request, trie.PrefixChan(vals, "Request"))
			}, // partial application
		)
		data := data.GetData(ctx, log) // TODO refresh on GET. TODO push updates to web UI (gin seems to support push)

		c.JSON(200, data)
	})
}
