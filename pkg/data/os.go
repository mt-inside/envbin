package data

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-logr/logr"
)

func init() {
	plugins = append(plugins, getOsData)
}

func getOsData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(strconv.Itoa(os.Getpid()), "Process", "ID")
	t.Insert(strconv.Itoa(os.Getppid()), "Process", "ParentID")

	t.Insert(strconv.Itoa(os.Getuid()), "Process", "UID")
	t.Insert(strconv.Itoa(os.Getgid()), "Process", "GID")
	if groups, err := os.Getgroups(); err == nil {
		t.Insert(fmt.Sprint(groups), "Process", "Groups")
	}

	if exe, err := os.Executable(); err == nil {
		t.Insert(exe, "Process", "Path")
	}
	if cwd, err := os.Getwd(); err == nil {
		t.Insert(cwd, "Process", "CWD")
	}

	/* TODO: capabilities */
}
