package fetchers

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strconv"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getOsData)
}

func getOsData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(Some(strconv.Itoa(os.Getpid())), "Process", "ID")

	pid := os.Getppid()
	i := 0
	for {
		var ps sigar.ProcState
		err := ps.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process", "PID", pid)
			t.Insert(Error(err), "Process", "Parents", strconv.Itoa(i))
			break
		}

		t.Insert(Some(strconv.Itoa(pid)), "Process", "Parents", strconv.Itoa(i), "PID")
		t.Insert(Some(ps.Name), "Process", "Parents", strconv.Itoa(i), "Name")
		t.Insert(Some(strconv.Itoa(ps.Priority)), "Process", "Parents", strconv.Itoa(i), "Priority")
		t.Insert(Some(strconv.Itoa(ps.Nice)), "Process", "Parents", strconv.Itoa(i), "Nice")

		if pid == 1 {
			break
		}

		pid = ps.Ppid
		i++
	}

	// Note: these are all properties of the process, but remember procs run _as_ users (unless they have the sticky bit set), so there's not really a concept of a proc identity
	t.Insert(Some(strconv.Itoa(os.Getuid())), "Process", "User", "UID")
	u, err := user.Current()
	if err != nil {
		t.Insert(Error(err), "Process", "User", "Details")
	} else {
		t.Insert(Some(u.Username), "Process", "User", "Username")
		t.Insert(Some(u.Name), "Process", "User", "Name")
		t.Insert(Some(u.HomeDir), "Process", "User", "Home")
	}

	t.Insert(Some(strconv.Itoa(os.Getgid())), "Process", "Group", "GID")
	g, err := user.LookupGroupId(strconv.Itoa(os.Getgid()))
	if err != nil {
		t.Insert(Error(err), "Process", "Group", "Details")
	} else {
		t.Insert(Some(g.Name), "Process", "Group", "Name")
	}
	if groups, err := os.Getgroups(); err == nil {
		t.Insert(Some(fmt.Sprint(groups)), "Process", "Groups")
	}

	if exe, err := os.Executable(); err == nil {
		t.Insert(Some(exe), "Process", "Path")
	}
	if cwd, err := os.Getwd(); err == nil {
		t.Insert(Some(cwd), "Process", "CWD")
	}

}
