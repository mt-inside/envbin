package fetchers

import (
	"context"
	"runtime"
	"strconv"
	"time"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

const Binary = "envbin"

var (
	Version  string
	TimeUnix string
)

func init() {
	data.RegisterPlugin(getBuildData)
}

func BuildTime() time.Time {
	unix, err := strconv.ParseInt(TimeUnix, 10, 64)
	if err != nil {
		unix = 0
	}

	return time.Unix(unix, 0)
}

func getBuildData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(Some(Version), "Build", "Version")
	t.Insert(Some(BuildTime().Format(time.Stamp)), "Build", "Time")
	t.Insert(Some(runtime.Version()), "Build", "Runtime")
	t.Insert(Some(runtime.GOARCH), "Hardware", "CPU", "Arch")
	t.Insert(Some(runtime.GOOS), "OS", "Kernel", "Type")
}
