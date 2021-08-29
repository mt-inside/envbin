package servers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
)

func GetEnv(log logr.Logger, r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		data := data.GetData(ctx, log) // TODO refresh on GET. TODO push updates to web UI (gin seems to support push)
		c.JSON(200, data)
	})
}
