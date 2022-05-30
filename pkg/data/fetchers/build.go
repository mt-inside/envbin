package fetchers

import (
	"context"
	"runtime"
	"strconv"
	"time"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
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

func getBuildData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	vals <- trie.Insert(trie.Some(Version), "Build", "Version")
	vals <- trie.Insert(trie.Some(BuildTime().Format(time.Stamp)), "Build", "Time")
	vals <- trie.Insert(trie.Some(runtime.Version()), "Build", "Runtime")
	vals <- trie.Insert(trie.Some(runtime.GOARCH), "Hardware", "CPU", "Arch")
	vals <- trie.Insert(trie.Some(runtime.GOOS), "OS", "Kernel", "Type")
}
