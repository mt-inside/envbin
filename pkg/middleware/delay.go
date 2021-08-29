package middleware

import (
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func Delay(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Duration(viper.GetInt64("DelayMS")) * time.Millisecond)
}
