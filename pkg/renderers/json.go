package renderers

import (
	"encoding/json"
	"log"
)

func RenderJSON(data map[string]string) (bs []byte) {
	var err error
	bs, err = json.Marshal(data)
	if err != nil {
		log.Fatalf("Can't encode data to JSON: %v", err)
	}

	return
}