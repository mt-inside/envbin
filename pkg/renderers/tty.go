package renderers

import (
	"strings"

	"github.com/fatih/color"
	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
)

var (
	//whiteBold = color.New(color.FgHiWhite).Add(color.Bold)
	white = color.New(color.FgHiWhite)
	norm  = color.New(color.FgWhite)
	//grey      = color.New(color.FgHiBlack)
)

// func title(s string) {
// 	fmt.Println()
// 	whiteBold.Printf("== %s ==\n", s)
// }

// TODO: cuirated list, and --all
func RenderTTY(log logr.Logger, data *data.Trie) {
	data.Walk(renderCb)
}

func renderCb(entry data.Entry) {
	depth := len(entry.Path)
	if depth == 0 {
		return
	}
	norm.Print(strings.Repeat("  ", depth-1))
	if _, ok := entry.Value.(data.None); ok {
		white.Printf("%s\n", entry.Path[depth-1])
	} else {
		white.Printf("%s: ", entry.Path[depth-1])
		norm.Printf("%s\n", entry.Value.Render())
	}
}

// title("Request")
// kv("Session", "%s", data["Session"])

// title("Hardware")
// kv("Virtualisation", "%s", data["Virt"])
// kv("Firmware", "%s", data["FirmwareType"])
// kv("Apparent hardware", "%s, %s, %s/%s cores, %s RAM", data["Arch"], data["CpuName"], data["PhysCores"], data["VirtCores"], data["MemTotal"])

// title("Operating Environment")
// kv("OS", "%s %s, up %s", data["OsType"], data["KernelVersion"], data["OsUptime"])
// kv("Distro", "%s (%s) %s (%s)", data["OsDistro"], data["OsFamily"], data["OsVersion"], data["OsRelease"])
// kv("PID", "%s, parent %s, #others %s", data["Pid"], data["Ppid"], data["OtherProcsCount"])
// kv("User", "UID %s (effective %s)", data["Uid"], data["Euid"])
// kv("Groups", "Primary %s (effective %s), others %s", data["Gid"], data["Egid"], data["Groups"])

// title("Network")
// kv("Hostname", "%s", data["Hostname"])
// kv("Primary IP", "%s", data["HostIp"])
// kv("External IP", "%s %s", data["ExternalIp"], data["ExternalIpEnrich"])
// // TODO: we control both ends of this interface and it's horrid!
// // FIXME: doesn't even work, cause interface indecies aren't necc sequential
// for i := 0; i < 128; i++ {
// 	v, ok := data[fmt.Sprintf("Interface%d", i)]
// 	if !ok {
// 		continue
// 	}
// 	kv(fmt.Sprintf("Iface[%d]", i), "%s", v)
// }

// func kv(key string, valFmt string, vals ...interface{}) {
// 	white.Printf("%s: ", key)
// 	norm.Printf(valFmt, vals...)
// 	fmt.Println()
// }
