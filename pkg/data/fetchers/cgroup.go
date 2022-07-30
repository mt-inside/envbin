package fetchers

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("cgroup", getCgroupData)
}

func getCgroupData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	fsen, err := os.Open("/proc/filesystems")
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "OS", "Isolation", "CGroups")
		return
	}
	defer fsen.Close()

	fsen_scanner := bufio.NewScanner(fsen)
	for fsen_scanner.Scan() {
		if strings.HasSuffix(fsen_scanner.Text(), "cgroup") {
			vals <- trie.Insert(trie.Some("Yes"), "OS", "Isolation", "CGroups", "v1", "Supported")
		} else if strings.HasSuffix(fsen_scanner.Text(), "cgroup2") {
			vals <- trie.Insert(trie.Some("Yes"), "OS", "Isolation", "CGroups", "v2", "Supported")
		}
	}

	mounts, err := os.Open("/proc/mounts")
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "OS", "Isolation", "CGroups")
		return
	}
	defer mounts.Close()

	mounts_scanner := bufio.NewScanner(mounts)
	for mounts_scanner.Scan() {
		if strings.HasPrefix(mounts_scanner.Text(), "cgroup2") {
			vals <- trie.Insert(trie.Some("Yes"), "OS", "Isolation", "CGroups", "v2", "Enabled")
		} else if strings.HasPrefix(mounts_scanner.Text(), "cgroup") {
			vals <- trie.Insert(trie.Some("Yes"), "OS", "Isolation", "CGroups", "v1", "Enabled")
		}
	}
}
