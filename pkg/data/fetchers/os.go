package fetchers

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getOsData)
}

func getOsData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(Some(strconv.Itoa(os.Getpid())), "Process", "ID")
	t.Insert(Some(strconv.Itoa(os.Getppid())), "Process", "ParentID")

	t.Insert(Some(strconv.Itoa(os.Getuid())), "Process", "UID")
	t.Insert(Some(strconv.Itoa(os.Getgid())), "Process", "GID")
	if groups, err := os.Getgroups(); err == nil {
		t.Insert(Some(fmt.Sprint(groups)), "Process", "Groups")
	}

	if exe, err := os.Executable(); err == nil {
		t.Insert(Some(exe), "Process", "Path")
	}
	if cwd, err := os.Getwd(); err == nil {
		t.Insert(Some(cwd), "Process", "CWD")
	}

	/* TODO: capabilities */
}
