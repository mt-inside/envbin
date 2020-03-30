package data

import (
	"fmt"
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/docker/go-units"
	"runtime"
)

func getMemData() map[string]string {
	data := map[string]string{}

	mem := sigar.Mem{}
	mem.Get()

	ms := new(runtime.MemStats)
	runtime.ReadMemStats(ms)

	data["MemTotal"] = units.BytesSize(float64(mem.Total))
	data["GcRuns"] = fmt.Sprintf("%d (%d forced)", ms.NumGC, ms.NumForcedGC)

	return data
}