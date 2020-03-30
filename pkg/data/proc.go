package data

import (
	"fmt"
	"github.com/docker/go-units"
	"github.com/prometheus/procfs"
	"log"
	"os"
	"runtime"
	"strconv"
)

func getProcData() map[string]string {
	data := map[string]string{}

	ms := new(runtime.MemStats)
	runtime.ReadMemStats(ms)

	proc, err := procfs.Self()
	if err != nil {
		log.Fatalf("Can't get proc info: %v", err)
	}
	stat, err := proc.NewStat()
	if err != nil {
		log.Fatalf("Can't get proc stat: %v", err)
	}

	data["Pid"] = strconv.Itoa(os.Getpid())
	data["Uid"] = strconv.Itoa(os.Getuid())
	data["Gid"] = strconv.Itoa(os.Getgid())
	data["MemUseVirtual"] = fmt.Sprintf("%s (of which %s golang runtime)",
		units.BytesSize(float64(stat.VirtualMemory())),
		units.BytesSize(float64(ms.Sys)),
	)
	data["MemUsePhysical"] = units.BytesSize(float64(stat.ResidentMemory()))
	data["CpuSelfTime"] = strconv.FormatFloat(stat.CPUTime(), 'f', 2, 64) + "s"

	return data
}
