package renderers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/mt-inside/envbin/pkg/data"
)

func RenderHTML(r *http.Request) (bs []byte) {
	data := data.GetData(r)
	var b bytes.Buffer
	t, err := template.ParseFiles("html.tpl")
	if err != nil {
		log.Fatalf("Failed to parse template html.tpl: %v", err)
	}
	t.Execute(&b, data)
	bs = b.Bytes()

	return
}
