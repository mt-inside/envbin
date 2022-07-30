package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/jsonquery"
	"github.com/docker/go-units"
	"github.com/fatih/color"
	"github.com/mt-inside/envbin/pkg/data/fetchers"
	"github.com/mt-inside/go-usvc"
	"github.com/urfave/cli/v2"
)

func main() {
	log := usvc.GetLogger(true, 0)

	app := &cli.App{
		Name:     "envbinctl",
		Usage:    "A CLI client for envbin",
		Version:  fetchers.Version,
		Compiled: fetchers.BuildTime(),

		UseShortOptionHandling: true,
		EnableBashCompletion:   true, // TODO not working

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Value: "http://localhost:8081",
				Usage: "URL of the envbin daemon",
			},
		},

		Metadata: map[string]interface{}{
			"log": log,
		},

		Action: render,
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

var (
	whiteBold = color.New(color.FgHiWhite).Add(color.Bold)
	white     = color.New(color.FgHiWhite)
	norm      = color.New(color.FgWhite)
	grey      = color.New(color.FgHiBlack)
)

// TODO: really want own printf formatter, where if any of the interpolations error, return "" for the whole thing.
// - This is doable: return custom types and impliment Format() for them: https://stackoverflow.com/questions/61539121/golang-custom-type-fmt-printing
// TODO: when we get type info with the values, just impliment Format() for jsonquery.Node (as the callsite won't need to state which type it is)
func s(root *jsonquery.Node, path string) string {
	node := jsonquery.FindOne(root, path)
	if node != nil {
		return node.Value().(string)
	}
	return "<none>"
}
func b(node *jsonquery.Node, path string) bool {
	item := jsonquery.FindOne(node, path)
	if item == nil {
		return false // TODO this is why we need a proper formatter, because need to be able to print "unknown"
	}

	t, err := strconv.ParseBool(item.Value().(string))
	if err != nil {
		return false
	}

	return t
}
func i(node *jsonquery.Node, path string) int64 { //nolint:deadcode
	item := jsonquery.FindOne(node, path)
	if item == nil {
		return 0
	}

	n, err := strconv.ParseInt(item.Value().(string), 10, 64)
	if err != nil {
		return 0
	}

	return n
}
func f(node *jsonquery.Node, path string) float64 {
	item := jsonquery.FindOne(node, path)
	if item == nil {
		return 0
	}

	n, err := strconv.ParseFloat(item.Value().(string), 64)
	if err != nil {
		return 0
	}

	return n
}

func render(c *cli.Context) error {
	//log := c.App.Metadata["log"].(logr.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := &http.Client{}

	base, err := url.Parse(c.String("url"))
	if err != nil {
		return err
	}
	path, _ := url.Parse("/api/v1/env")

	req, err := http.NewRequestWithContext(ctx, "GET", base.ResolveReference(path).String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	// Also unmarshals
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// data := data.NewTrie(log)
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// 	return err
	// }

	root, err := jsonquery.Parse(resp.Body)
	if err != nil {
		return err
	}

	renderSummary(root)
	norm.Println()

	if true {
		renderCache(root)
		norm.Println()
	}

	if true {
		renderRAM(root)
		norm.Println()
	}

	if true {
		renderNetIfaces(root)
		norm.Println()
	}

	if true {
		renderPCI(root)
		norm.Println()
	}

	if true {
		renderUSB(root)
		norm.Println()
	}

	if true {
		renderAlsa(root)
		norm.Println()
	}

	if true {
		renderV4l2(root)
		norm.Println()
	}

	if true {
		renderBlock(root)
		norm.Println()
	}

	return nil
}

func renderProduct(node *jsonquery.Node) {
	white.Print(s(node, "Vendor"))
	whiteBold.Print(" " + s(node, "Model"))
	grey.Printf(" (%s, %s)", s(node, "SKU"), s(node, "Version"))
}

func renderSummary(root *jsonquery.Node) {
	// FIXME: jsonquery doesn't seem to like a number as path element (map key), so we manually navigate. Except I don't think it's gaurenteed that '0' will be the first child
	thisProc := jsonquery.FindOne(root, "Processes").FirstChild
	whiteBold.Print(s(thisProc, "Exe"))
	norm.Print(" user ")
	whiteBold.Print(s(thisProc, "User/Name"))
	norm.Print(" group ")
	whiteBold.Print(s(thisProc, "Group/Name"))
	norm.Println()

	norm.Println()

	whiteBold.Print(s(root, "Network/Hostname"))
	norm.Print(" " + s(root, "Network/DefaultIP"))
	norm.Print(" / " + s(root, "Network/ExternalIP/Address"))
	grey.Printf(" (%s: %s, %s, %s, %s, %s)", s(root, "Network/ExternalIP/ReverseDNS"), s(root, "Network/ExternalIP/AS"), s(root, "Network/ExternalIP/City"), s(root, "Network/ExternalIP/Region"), s(root, "Network/ExternalIP/Postal"), s(root, "Network/ExternalIP/Country"))
	norm.Println()

	whiteBold.Print(s(root, "OS/Distro/Release"))
	norm.Print(" " + s(root, "OS/Distro/Version"))
	grey.Printf(" (%s %s)", s(root, "OS/Kernel/Type"), s(root, "OS/Kernel/Version"))
	norm.Print(" up " + s(root, "OS/Uptime"))
	norm.Println()

	whiteBold.Print(s(root, "Hardware/Firmware/BootType"))
	norm.Print(" boot ")
	whiteBold.Print(s(root, "Hardware/Virtualisation"))
	norm.Println()

	norm.Println()

	sys := jsonquery.FindOne(root, "Hardware/System")
	renderProduct(sys)
	norm.Println()

	cpu := jsonquery.FindOne(root, "Hardware/CPU")
	renderProduct(cpu)
	white.Printf(" %s/%s", s(root, "Hardware/CPU/Cores"), s(root, "Hardware/CPU/Threads"))
	white.Printf(" %s", s(root, "Hardware/CPU/Arch"))
	grey.Printf(" (%s, %s)", s(root, "Hardware/CPU/Model/Microarchitecture"), s(root, "Hardware/CPU/Package"))
	white.Printf(" %d/%dMHz", i(root, "Hardware/CPU/Clock/Current"), i(root, "Hardware/CPU/Clock/Max"))
	norm.Println()

	whiteBold.Print(units.BytesSize(f(root, "Hardware/Memory/Total")))
	norm.Print(" RAM")
	norm.Println()
}

func renderCache(root *jsonquery.Node) {
	l1each := i(root, "Hardware/CPU/Cache/Individual/Level1")
	l2each := i(root, "Hardware/CPU/Cache/Individual/Level2")
	l3each := i(root, "Hardware/CPU/Cache/Individual/Level3")
	l1total := i(root, "Hardware/CPU/Cache/Totals/Level1")
	l2total := i(root, "Hardware/CPU/Cache/Totals/Level2")
	l3total := i(root, "Hardware/CPU/Cache/Totals/Level3")

	white.Print("Cache")
	whiteBold.Printf(" L1 %dx%dkB", l1total/l1each, l1each>>10) // en attendant the unit system
	whiteBold.Printf(" L2 %dx%dkB", l2total/l2each, l2each>>10)
	whiteBold.Printf(" L3 %dx%dMB", l3total/l3each, l3each>>20)
	norm.Println()
}

func renderRAM(root *jsonquery.Node) {
	for _, dimm := range jsonquery.Find(root, "Hardware/RAM/*") {
		if i(dimm, "SizeMB") == 0 {
			continue
		}

		white.Printf("%s %s", s(dimm, "Channel"), s(dimm, "Slot"))
		whiteBold.Printf(" %sMB %s", s(dimm, "SizeMB"), s(dimm, "Standard"))
		norm.Printf(" %dmV", i(dimm, "Voltage/CurrentmV"))
		grey.Printf(" (%d/%dbits)", i(dimm, "Bus/Width/Data"), i(dimm, "Bus/Width/Total"))
		norm.Println()
	}
}

func renderNetIfaces(root *jsonquery.Node) {
	for _, iface := range jsonquery.Find(root, "Network/Interfaces/*") {
		whiteBold.Print(s(iface, "Name"))
		grey.Printf(" %s %s", s(iface, "Driver"), s(iface, "IfaceFlags"))
		norm.Print(" | ")
		norm.Print(s(iface, "IPv4Address"))
		norm.Print(" | ")
		norm.Print(s(iface, "MACAddr"))
		norm.Print(" | ")
		norm.Printf("link %t, speed %d", b(iface, "PhysicalLink"), i(iface, "SpeedMbits"))
		grey.Print("Mb/s")
		norm.Println()
	}
}

func renderPCI(root *jsonquery.Node) {
	for _, dev := range jsonquery.Find(root, "Hardware/Bus/PCI/*") {
		white.Print(dev.Data)
		fns := jsonquery.Find(dev, "Functions/*")
		if len(fns) != 1 {
			norm.Println()
		}

		for _, fn := range fns {
			norm.Print("  ")
			whiteBold.Print(s(fn, "Vendor"))
			whiteBold.Print(" " + s(fn, "Product"))
			grey.Print(" rev " + s(fn, "Revision"))
			grey.Printf(" (%s / %s, driver %s)", s(fn, "Class"), s(fn, "Subclass"), s(fn, "Driver"))
			norm.Println()
		}
	}
}

func renderUSB(root *jsonquery.Node) {
	for _, dev := range jsonquery.Find(root, "Hardware/Bus/USB/*") {
		white.Print(dev.Data)
		whiteBold.Printf(" %s %s", s(dev, "Manufacturer"), s(dev, "Product"))

		serial := s(dev, "Serial")
		spec := s(dev, "Spec")
		speed := s(dev, "Speed")

		if serial != "" {
			grey.Printf(" serial %s", serial)
		}
		grey.Printf(" [usb %s", spec)
		if speed != "" {
			grey.Printf(" %s speed", speed)
		}
		grey.Printf("]")

		fns := jsonquery.Find(dev, "Configs/*")
		if len(fns) != 1 {
			norm.Println()
		}

		for _, fn := range fns {
			norm.Print("  ")
			norm.Printf(" %s", s(fn, "Power"))
			norm.Print(" [")
			if s(fn, "Wakeup") == "true" {
				norm.Printf("Wakeup")
			}
			norm.Print("]")

			ifaces := jsonquery.Find(fn, "Interfaces/*")
			if len(ifaces) != 1 {
				norm.Println()
			}
			for _, iface := range ifaces {
				norm.Print("    ")
				norm.Print(s(iface, "Description"))
				grey.Printf(" driver %s", s(iface, "Driver"))
				norm.Println()
			}
		}
	}
}

func renderV4l2(root *jsonquery.Node) {
	for _, dev := range jsonquery.Find(root, "Hardware/V4l2/*") {
		white.Print(dev.Data)
		whiteBold.Printf(" %s", s(dev, "Name"))
		white.Printf(" driver %s", s(dev, "Driver"))
		grey.Printf(" video out %s, capture %s, streaming %s", s(dev, "Capabilities/VideoOutput"), s(dev, "Capabilities/VideoCapture"), s(dev, "Capabilities/StreamingIO"))

		fmts := jsonquery.Find(dev, "Formats/*")
		if len(fmts) != 1 {
			norm.Println()
		}

		for _, fmt := range fmts {
			norm.Print("  ")
			norm.Print(s(fmt, "Name"))
			grey.Printf(" compressed %s, emulated %s", s(fmt, "Compressed"), s(fmt, "Emulated"))
			norm.Println()
		}
	}
}

func renderAlsa(root *jsonquery.Node) {
	for _, dev := range jsonquery.Find(root, "Hardware/Sound/Alsa/Cards/*") {
		white.Print(s(dev, "Path"))
		whiteBold.Printf(" %s", s(dev, "Name"))

		devs := jsonquery.Find(dev, "Devices/*")
		if len(devs) != 1 {
			norm.Println()
		} else {
			norm.Print(" | ")
		}

		for _, dev := range devs {
			norm.Print("  ")

			white.Print(s(dev, "Name"))
			norm.Printf(" %s", s(dev, "Type"))

			norm.Printf(" %sch %s/s x %s", s(dev, "Channels"), s(dev, "Sample/Rate"), s(dev, "Sample/Format"))

			flags := []string{}
			if s(dev, "Play") == "true" {
				flags = append(flags, "play")
			}
			if s(dev, "Record") == "true" {
				flags = append(flags, "record")
			}
			grey.Printf(" [%s]", strings.Join(flags, ", "))

			norm.Println()
		}
	}
}

func renderBlock(root *jsonquery.Node) {
	for _, blk := range jsonquery.Find(root, "Hardware/Block/*") {
		white.Print(blk.Data)
		whiteBold.Printf(" %s %s", s(blk, "Vendor"), s(blk, "Model"))

		serial := s(blk, "Serial")
		if serial != "" {
			grey.Printf(" serial %s", serial)
		}

		norm.Printf(" [%s, %s bytes", s(blk, "ControllerType"), units.HumanSize(f(blk, "SizeBytes")))
		if s(blk, "Removable") == "true" {
			norm.Printf("Removable")
		}
		norm.Print("]")

		norm.Println()

		ps := jsonquery.Find(blk, "Partitions/*")
		for _, p := range ps {
			norm.Print("  ")
			white.Print(p.Data)
			if s(p, "Filesystem") != "NotPresent" {
				whiteBold.Printf(" %s", s(p, "Filesystem"))
			}
			if s(p, "MountPoint") != "NotPresent" {
				norm.Print(" on")
				whiteBold.Printf(" %s", s(p, "MountPoint"))
			}
			grey.Printf(" uuid %s", s(p, "UUID"))
			norm.Println()
		}
	}
}

// func title(s string) {
// 	fmt.Println()
// 	whiteBold.Printf("== %s ==\n", s)
// }

// func kv(key string, valFmt string, vals ...interface{}) {
// 	white.Printf("%s: ", key)
// 	norm.Printf(valFmt, vals...)
// 	fmt.Println()
// }

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
