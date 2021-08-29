package servers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type configPath struct {
	Name string `uri:"name" binding:"required"`
}
type configValue struct {
	Value interface{} `form:"value" json:"value" binding:"required"`
}

func GetConfig(log logr.Logger, r *gin.RouterGroup) {
	r.POST("/:name", func(c *gin.Context) {
		var path configPath
		if err := c.ShouldBindUri(&path); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		old := viper.Get(path.Name)
		if old == nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "No such config item", "item": path.Name})
			return
		}

		var value configValue
		if err := c.ShouldBind(&value); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
			return
		}

		// TODO: fix me? {"Value": 1000} is being unmarshaled as float64. Think the way to do it is:
		// - get the item name
		// - get the old item
		// - switch old.(type)
		// - bind to a different model for string, bool, int, etc
		// - set accordingly
		//
		// if reflect.TypeOf(old) != reflect.TypeOf(value.Value) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"msg": "Attempt to change type", "old": reflect.TypeOf(old).String(), "new": reflect.TypeOf(value.Value).String()})
		// 	return
		// }

		viper.Set(path.Name, value.Value)
		log.Info("Config changed", "item", path.Name, "old", old, "new", value.Value)
		c.JSON(http.StatusOK, gin.H{"item": path.Name, "old": old, "new": value.Value})
	})
}
