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

func getMemData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		log.Error(err, "Can't read memory information")
		vals <- Insert(Error(err), "Hardware", "Memory")
		return
	}

	vals <- Insert(Some(strconv.FormatUint(mem.Total, 10)), "Hardware", "Memory", "Total")
}

func getProcsData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	procs := sigar.ProcList{}
	err := procs.Get()
	if err != nil {
		log.Error(err, "Can't read process information")
		vals <- Insert(Error(err), "OS", "Processes")
		return
	}

	vals <- Insert(Some(strconv.Itoa(len(procs.List)-1)), "OS", "Processes", "Count")
}

func getOsData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	pid := os.Getpid()
	i := 0
	for {
		var ps sigar.ProcState
		err := ps.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process", "PID", pid)
			vals <- Insert(Error(err), "Processes", strconv.Itoa(i))
			break
		}
		vals <- Insert(Some(strconv.Itoa(pid)), "Processes", strconv.Itoa(i), "PID")
		vals <- Insert(Some(strconv.Itoa(ps.Pgid)), "Processes", strconv.Itoa(i), "PGID")
		vals <- Insert(Some(ps.Name), "Processes", strconv.Itoa(i), "Name")
		vals <- Insert(Some(strconv.Itoa(ps.Priority)), "Processes", strconv.Itoa(i), "Priority")
		vals <- Insert(Some(strconv.Itoa(ps.Nice)), "Processes", strconv.Itoa(i), "Nice")

		var args sigar.ProcArgs
		err = args.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process args", "PID", pid)
			vals <- Insert(Error(err), "Processes", strconv.Itoa(i), "Details")
			break
		}
		vals <- Insert(Some(strings.Join(args.List, " ")), "Processes", strconv.Itoa(i), "Cmdline")

		var exe sigar.ProcExe
		err = exe.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process exe", "PID", pid)
			vals <- Insert(Error(err), "Processes", strconv.Itoa(i), "Details")
			break
		}
		vals <- Insert(Some(exe.Name), "Processes", strconv.Itoa(i), "Exe")
		vals <- Insert(Some(exe.Cwd), "Processes", strconv.Itoa(i), "Cwd")
		vals <- Insert(Some(exe.Root), "Processes", strconv.Itoa(i), "Root")

		var env sigar.ProcEnv
		err = env.Get(pid)
		if err != nil {
			log.Error(err, "Can't get process env", "PID", pid)
			vals <- Insert(Error(err), "Processes", strconv.Itoa(i), "Details")
			break
		}
		vals <- Insert(Some(env.Vars["PATH"]), "Processes", strconv.Itoa(i), "Path")

		// Note: these are all properties of the process, but remember procs run _as_ users (unless they have the sticky bit set), so there's not really a concept of a proc identity
		vals <- Insert(Some(ps.Username), "Processes", strconv.Itoa(i), "User", "Name")
		u, err := user.Lookup(ps.Username)
		if err != nil {
			log.Error(err, "Can't get process user", "PID", pid)
			vals <- Insert(Error(err), "Processes", strconv.Itoa(i), "User", "Details")
			break
		}
		vals <- Insert(Some(u.Uid), "Processes", strconv.Itoa(i), "User", "UID")
		vals <- Insert(Some(u.Name), "Processes", strconv.Itoa(i), "User", "Full Name")
		vals <- Insert(Some(u.HomeDir), "Processes", strconv.Itoa(i), "User", "Home")

		vals <- Insert(Some(u.Gid), "Processes", strconv.Itoa(i), "Group", "GID")
		g, err := user.LookupGroupId(u.Gid)
		if err != nil {
			log.Error(err, "Can't get process group", "PID", pid)
			vals <- Insert(Error(err), "Processes", strconv.Itoa(i), "Group", "Details")
			break
		}
		vals <- Insert(Some(g.Name), "Processes", strconv.Itoa(i), "Group", "Name")

		if pid == 1 {
			break
		}

		pid = ps.Ppid
		i++
	}

	if groups, err := os.Getgroups(); err == nil {
		vals <- Insert(Some(fmt.Sprint(groups)), "Processes", "0", "Group", "Others")
	}
}
