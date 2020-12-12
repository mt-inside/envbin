package data

import (
	"os"
	"strconv"
)

func getProcData() map[string]string {
	data := map[string]string{}

	data["Pid"] = strconv.Itoa(os.Getpid())
	//ppid
	data["Uid"] = strconv.Itoa(os.Getuid())
	data["Gid"] = strconv.Itoa(os.Getgid())

	return data
}
