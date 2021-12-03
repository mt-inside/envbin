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
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		data := data.GetData(ctx, log) // TODO refresh on GET. TODO push updates to web UI (gin seems to support push)

		reqData := trie.BuildFromSyncFn(
			log,
			func(vals chan<- trie.InsertMsg) { extractors.RequestData(ctx, log, c.Request, vals) }, // partial application
		)
		data.InsertTree(reqData, "Request")

		c.JSON(200, data)
	})
}
