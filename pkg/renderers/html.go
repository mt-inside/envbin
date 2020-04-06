package renderers

import (
	"bytes"
	"html/template"
	"log"
)

func RenderHTML(data map[string]string) (bs []byte) {
	var b bytes.Buffer
	t, err := template.ParseFiles("html.html")
	if err != nil {
		log.Fatalf("Failed to parse template html.html: %v", err)
	}
	t.Execute(&b, data)
	bs = b.Bytes()

	return
}
