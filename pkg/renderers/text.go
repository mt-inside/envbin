package renderers

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
)

func RenderText(log logr.Logger, w http.ResponseWriter, r *http.Request, data *data.Trie) (bs []byte) {
	w.Header().Set("Content-Type", "text/plain")

	return []byte(data.Render())

	/* This does the application/text output quite nicely, but for a fancy HTML page we probably want:
	* - gorilla mux SPA example
	* - SPA (react etc) which can be made elsewhere and loaded with gobindata (to avoid the complexity of hosting it behing a separate web server. Or maybe we do, in the same container / Pod?)
	* - JSON handlers for this struct (make it a struct and JSON serialse it) so it can be read by the SPA
	 */
	// TODO: actually marshall to YAML. Make the map[s]s be a real struct
	//
	// t, err := template.ParseFiles("text.tpl")
	// if err != nil {
	// 	log.Error(err, "Failed to parse template text.tpl")
	// }

	// var b bytes.Buffer
	// t.Execute(&b, data)

	// return b.Bytes()
}
