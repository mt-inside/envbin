package renderers

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
	"gopkg.in/yaml.v2"
)

func RenderYAML(log logr.Logger, w http.ResponseWriter, r *http.Request, data *data.Trie) []byte {
	w.Header().Set("Content-Type", "text/yaml")

	bs, err := yaml.Marshal(data)
	if err != nil {
		log.Error(err, "Can't encode data to YAML")
	}

	return bs
}
