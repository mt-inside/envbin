package fetchers

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	sigar "github.com/elastic/gosigar"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getMemData)
	data.RegisterPlugin(getProcsData)
	data.RegisterPlugin(getOsData)
}

func getMemData(ctx context.Context, log logr.Logger, t *Trie) {
	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		log.Error(err, "Can't read memory information")
		t.Insert(Error(err), "Hardware", "Memory")
		return
	}

	t.Insert(Some(strconv.FormatUint(mem.Total, 10)), "Hardware", "Memory", "Total")
}

func getProcsData(ctx context.Context, log logr.Logger, t *Trie) {
	procs := sigar.ProcList{}
	err := procs.Get()
	if err != nil {
		log.Error(err, "Can't read process information")
		t.Insert(Error(err), "OS", "Processes")
		return
	}

	t.Insert(Some(strconv.Itoa(len(procs.List)-1)), "OS", "Processes", "Count")
}

func getOsData(ctx context.Context, log logr.Logger, t *Trie) {
	pid := os.Getpid()
	i := 0
	for {
		var ps sigar.ProcState
		err := ps.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process", "PID", pid)
			t.Insert(Error(err), "Processes", strconv.Itoa(i))
			break
		}
		t.Insert(Some(strconv.Itoa(pid)), "Processes", strconv.Itoa(i), "PID")
		t.Insert(Some(strconv.Itoa(ps.Pgid)), "Processes", strconv.Itoa(i), "PGID")
		t.Insert(Some(ps.Name), "Processes", strconv.Itoa(i), "Name")
		t.Insert(Some(strconv.Itoa(ps.Priority)), "Processes", strconv.Itoa(i), "Priority")
		t.Insert(Some(strconv.Itoa(ps.Nice)), "Processes", strconv.Itoa(i), "Nice")

		var args sigar.ProcArgs
		err = args.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process args", "PID", pid)
			t.Insert(Error(err), "Processes", strconv.Itoa(i), "Details")
			break
		}
		t.Insert(Some(strings.Join(args.List, " ")), "Processes", strconv.Itoa(i), "Cmdline")

		var exe sigar.ProcExe
		err = exe.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process exe", "PID", pid)
			t.Insert(Error(err), "Processes", strconv.Itoa(i), "Details")
			break
		}
		t.Insert(Some(exe.Name), "Processes", strconv.Itoa(i), "Exe")
		t.Insert(Some(exe.Cwd), "Processes", strconv.Itoa(i), "Cwd")
		t.Insert(Some(exe.Root), "Processes", strconv.Itoa(i), "Root")

		var env sigar.ProcEnv
		err = env.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process env", "PID", pid)
			t.Insert(Error(err), "Processes", strconv.Itoa(i), "Details")
			break
		}
		t.Insert(Some(env.Vars["PATH"]), "Processes", strconv.Itoa(i), "Path")

		// Note: these are all properties of the process, but remember procs run _as_ users (unless they have the sticky bit set), so there's not really a concept of a proc identity
		t.Insert(Some(ps.Username), "Processes", strconv.Itoa(i), "User", "Name")
		u, err := user.Lookup(ps.Username)
		if err != nil {
			log.Error(err, "Can't get process user", "PID", pid)
			t.Insert(Error(err), "Processes", strconv.Itoa(i), "User", "Details")
			break
		}
		t.Insert(Some(u.Uid), "Processes", strconv.Itoa(i), "User", "UID")
		t.Insert(Some(u.Name), "Processes", strconv.Itoa(i), "User", "Full Name")
		t.Insert(Some(u.HomeDir), "Processes", strconv.Itoa(i), "User", "Home")

		t.Insert(Some(u.Gid), "Processes", strconv.Itoa(i), "Group", "GID")
		g, err := user.LookupGroupId(u.Gid)
		if err != nil {
			log.Error(err, "Can't get process group", "PID", pid)
			t.Insert(Error(err), "Processes", strconv.Itoa(i), "Group", "Details")
			break
		}
		t.Insert(Some(g.Name), "Processes", strconv.Itoa(i), "Group", "Name")

		if pid == 1 {
			break
		}

		pid = ps.Ppid
		i++
	}

	if groups, err := os.Getgroups(); err == nil {
		t.Insert(Some(fmt.Sprint(groups)), "Processes", "0", "Group", "Others")
	}
}
