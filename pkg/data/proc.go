package data

import (
	"fmt"
	"os"
	"strconv"
)

func getProcData() map[string]string {
	data := map[string]string{}

	data["Pid"] = strconv.Itoa(os.Getpid())
	data["Ppid"] = strconv.Itoa(os.Getppid())

	data["Uid"] = strconv.Itoa(os.Getuid())
	data["Euid"] = strconv.Itoa(os.Geteuid())
	data["Gid"] = strconv.Itoa(os.Getgid())
	data["Egid"] = strconv.Itoa(os.Getegid())
	if groups, err := os.Getgroups(); err == nil {
		data["Groups"] = fmt.Sprint(groups)
	}

	if cwd, err := os.Getwd(); err == nil {
		data["Cwd"] = cwd
	}

	/* TODO: capabilities */

	return data
}
