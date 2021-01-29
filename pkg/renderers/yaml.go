package renderers

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/mt-inside/envbin/pkg/data"
)

func RenderYAML(r *http.Request) (bs []byte) {
	data := data.GetData(r)
	var err error
	bs, err = yaml.Marshal(data)
	if err != nil {
		log.Fatalf("Can't encode data to YAML: %v", err)
	}

	return
}
