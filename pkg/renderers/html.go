package renderers

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
)

func RenderHTML(log logr.Logger, w http.ResponseWriter, r *http.Request, data *data.Trie) []byte {
	w.Header().Set("Content-Type", "text/html")

	t, err := template.ParseFiles("html.tpl")
	if err != nil {
		log.Error(err, "Failed to parse template html.tpl")
		return nil
	}

	var b bytes.Buffer
	t.Execute(&b, data)

	return b.Bytes()
}
