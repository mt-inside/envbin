package renderers

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mt-inside/envbin/pkg/data"
)

func RenderTTY() {
	data := data.GetData(nil)

	title("Request")
	kv("Session", "%s", data["Session"])

	title("Hardware")
	kv("Virtualisation", "%s", data["Virt"])
	kv("Firmware", "%s", data["FirmwareType"])
	kv("Apparent hardware", "%s, %s, %s/%s cores, %s RAM", data["Arch"], data["CpuName"], data["PhysCores"], data["VirtCores"], data["MemTotal"])

	title("Operating Environment")
	kv("OS", "%s %s, up %s", data["OsType"], data["KernelVersion"], data["OsUptime"])
	kv("Distro", "%s (%s) %s", data["OsDistro"], data["OsFamily"], data["OsVersion"])
	kv("REAL", "%s", data["OsReal"])
	kv("PID", "%s, parent %s, #others %s", data["Pid"], data["Ppid"], data["OtherProcsCount"])
	kv("User", "UID %s (effective %s)", data["Uid"], data["Euid"])
	kv("Groups", "Primary %s (effective %s), others %s", data["Gid"], data["Egid"], data["Groups"])

	title("Network")
	kv("Hostname", "%s", data["Hostname"])
	kv("Primary IP", "%s", data["HostIp"])
	kv("External IP", "%s %s", data["ExternalIp"], data["ExternalIpEnrich"])
	// TODO: we control both ends of this interface and it's horrid!
	// FIXME: doesn't even work, cause interface indecies aren't necc sequential
	for i := 0; i < 128; i++ {
		v, ok := data[fmt.Sprintf("Interface%d", i)]
		if !ok {
			continue
		}
		kv(fmt.Sprintf("Iface[%d]", i), "%s", v)
	}
}

var (
	whiteBold = color.New(color.FgHiWhite).Add(color.Bold)
	white     = color.New(color.FgHiWhite)
	norm      = color.New(color.FgWhite)
	grey      = color.New(color.FgHiBlack)
)

func title(s string) {
	fmt.Println()
	whiteBold.Printf("== %s ==\n", s)
}

func kv(key string, valFmt string, vals ...interface{}) {
	white.Printf("%s: ", key)
	norm.Printf(valFmt, vals...)
	fmt.Println()
}
