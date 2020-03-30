package data

import (
	"github.com/klauspost/cpuid"
	"github.com/shirou/gopsutil/host"
	"runtime"
	"strconv"
)

func getHardwareData() map[string]string {
	data := map[string]string{}

	is, _ := host.Info()

	data["Arch"] = runtime.GOARCH
	data["CpuName"] = cpuid.CPU.BrandName
	data["PhysCores"] = strconv.Itoa(cpuid.CPU.PhysicalCores)
	data["VirtCores"] = strconv.Itoa(cpuid.CPU.LogicalCores)
	data["Virt"] = is.VirtualizationSystem

	return data
}
