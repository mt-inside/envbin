package servers

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/spf13/viper"

	"github.com/mt-inside/envbin/pkg/dynamicthrottledreader"
)

const loremIpsum = `Sed ut perspiciatis, unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam eaque ipsa, quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt, explicabo.
Nemo enim ipsam voluptatem, quia voluptas sit, aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos, qui ratione voluptatem sequi nesciunt, neque porro quisquam est, qui do lorem ipsum, quia dolor sit amet consectetur adipisci velit, sed quia non numquam eius modi tempora incidunt, ut labore et dolore magnam aliquam quaerat voluptatem.
Ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur? Quis autem vel eum iure reprehenderit, qui in ea voluptate velit esse, quam nihil molestiae consequatur, vel illum, qui dolorem eum fugiat, quo voluptas nulla pariatur? At vero eos et accusamus et iusto odio dignissimos ducimus, qui blanditiis praesentium voluptatum deleniti atque corrupti, quos dolores et quas molestias excepturi sint, obcaecati cupiditate non provident, similique sunt in culpa, qui officia deserunt mollitia animi, id est laborum et dolorum fuga.
Et harum quidem rerum facilis est et expedita distinctio.
Nam libero tempore, cum soluta nobis est eligendi optio, cumque nihil impedit, quo minus id, quod maxime placeat, facere possimus, omnis voluptas assumenda est, omnis dolor repellendus.
Temporibus autem quibusdam et aut officiis debitis aut rerum necessitatibus saepe eveniet, ut et voluptates repudiandae sint et molestiae non recusandae.
Itaque earum rerum hic tenetur a sapiente delectus, ut aut reiciendis voluptatibus maiores alias consequatur aut perferendis doloribus asperiores repellat.
`

func init() {
	viper.SetDefault("BodyLen", len(loremIpsum)) // euw
}

func GetOutput(log logr.Logger, r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		tr := dynamicthrottledreader.NewDynamicThrottledReader(strings.NewReader(loremIpsum))
		c.DataFromReader(http.StatusOK, int64(len(loremIpsum)), mime.TypeByExtension(".txt"), tr, nil) // len(str) IS byte-len
	})
}
