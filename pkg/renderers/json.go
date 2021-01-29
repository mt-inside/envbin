package renderers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mt-inside/envbin/pkg/data"
)

func RenderJSON(r *http.Request) (bs []byte) {
	data := data.GetData(r)
	var err error
	bs, err = json.Marshal(data)
	if err != nil {
		log.Fatalf("Can't encode data to JSON: %v", err)
	}

	return
}
