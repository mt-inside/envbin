package renderers

import (
	"gopkg.in/yaml.v2"
	"log"
)

func RenderYAML(data map[string]string) (bs []byte) {
	var err error
	bs, err = yaml.Marshal(data)
	if err != nil {
		log.Fatalf("Can't encode data to YAML: %v", err)
	}

	return
}
