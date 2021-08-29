package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

func GetProbes(log logr.Logger, r *gin.RouterGroup) {
	r.GET("/live", func(c *gin.Context) {
		if viper.GetBool("Healthy") {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"ok": false,
			})
		}
	})

	r.GET("/ready", func(c *gin.Context) {
		if viper.GetBool("Ready") {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"ok": false,
			})
		}
	})
}
