package data

import (
	"net/http"

	"github.com/mt-inside/envbin/pkg/util"
)

var (
	plugins []func() map[string]string
)

func GetData(r *http.Request) map[string]string {
	d := make(map[string]string) //TODO: strongly type me with a struct. Esp for (optional) sections

	for _, p := range plugins {
		d = util.AppendMap(d, p())
	}
	d = util.AppendMap(d, getRequestData(r))

	return d
}
