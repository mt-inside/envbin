package renderers

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/mt-inside/envbin/pkg/data"
)

func RenderText(r *http.Request) (bs []byte) {
	data := data.GetData(r)
	/* This does the application/text output quite nicely, but for a fancy HTML page we probably want:
	* - gorilla mux SPA example
	* - SPA (react etc) which can be made elsewhere and loaded with gobindata (to avoid the complexity of hosting it behing a separate web server. Or maybe we do, in the same container / Pod?)
	* - JSON handlers for this struct (make it a struct and JSON serialse it) so it can be read by the SPA
	 */
	// TODO: actually marshall to YAML. Make the map[s]s be a real struct
	var b bytes.Buffer
	t, err := template.ParseFiles("text.tpl")
	if err != nil {
		log.Fatalf("Failed to parse template text.tpl: %v", err)
	}
	t.Execute(&b, data)
	bs = b.Bytes()

	return
}
