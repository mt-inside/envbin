package data

import (
	sigar "github.com/cloudfoundry/gosigar"
	"strconv"
)

func getProcsData() map[string]string {
	data := map[string]string{}

	procs := sigar.ProcList{}
	procs.Get()

	data["OtherProcsCount"] = strconv.Itoa(len(procs.List) - 1)

	return data
}
