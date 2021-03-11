package renderers

import (
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
)

func RenderJSON(log logr.Logger, w http.ResponseWriter, r *http.Request, data *data.Trie) []byte {
	w.Header().Set("Content-Type", "application/json")

	bs, err := json.Marshal(data)
	if err != nil {
		log.Error(err, "Can't encode data to JSON")
		return nil
	}

	return bs
}
