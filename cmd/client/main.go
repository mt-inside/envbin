package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/antchfx/jsonquery"
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

	if false {
		renderNetIfaces(root)
		norm.Println()
	}

	if false {
		renderPCI(root)
		norm.Println()
	}

	if false {
		renderUSB(root)
		norm.Println()
	}

	if true {
		renderBlock(root)
		norm.Println()
	}

	return nil
}

func renderSummary(root *jsonquery.Node) {
	whiteBold.Print(jsonquery.FindOne(root, "Network/Hostname").InnerText())
	norm.Print(" " + jsonquery.FindOne(root, "Network/DefaultIP").InnerText())
	norm.Print(" / " + jsonquery.FindOne(root, "Network/ExternalIP/Address").InnerText())
	grey.Printf(" (%s)", jsonquery.FindOne(root, "Network/ExternalIP/Info").InnerText())
	norm.Println()
	whiteBold.Print(jsonquery.FindOne(root, "OS/Distro/Release").InnerText())
	norm.Print(" " + jsonquery.FindOne(root, "OS/Distro/Version").InnerText())
	grey.Printf(" (%s %s)", jsonquery.FindOne(root, "OS/Kernel/Type").InnerText(), jsonquery.FindOne(root, "OS/Kernel/Version").InnerText())
	norm.Print(" up " + jsonquery.FindOne(root, "OS/Uptime").InnerText())
	norm.Println()
	whiteBold.Print(jsonquery.FindOne(root, "Hardware/CPU/Model").InnerText())
	white.Printf(" %s/%s", jsonquery.FindOne(root, "Hardware/CPU/Cores").InnerText(), jsonquery.FindOne(root, "Hardware/CPU/Threads").InnerText())
	white.Printf(" %s", jsonquery.FindOne(root, "Hardware/CPU/Arch").InnerText())
	norm.Println()
	whiteBold.Print(jsonquery.FindOne(root, "Hardware/Memory/Total").InnerText())
	norm.Println()
}

func renderNetIfaces(root *jsonquery.Node) {
	for _, iface := range jsonquery.Find(root, "Network/Interfaces/*") {
		whiteBold.Print(jsonquery.FindOne(iface, "Name").InnerText())
		norm.Print(" " + jsonquery.FindOne(iface, "Address").InnerText())
		grey.Print(" " + jsonquery.FindOne(iface, "Flags").InnerText())
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
			whiteBold.Print(jsonquery.FindOne(fn, "Vendor").InnerText())
			whiteBold.Print(" " + jsonquery.FindOne(fn, "Product").InnerText())
			grey.Print(" rev " + jsonquery.FindOne(fn, "Revision").InnerText())
			grey.Printf(" (%s / %s)", jsonquery.FindOne(fn, "Class").InnerText(), jsonquery.FindOne(fn, "Subclass").InnerText())
			norm.Println()
		}
	}
}

func renderUSB(root *jsonquery.Node) {
	for _, dev := range jsonquery.Find(root, "Hardware/Bus/USB/*") {
		white.Print(dev.Data)
		whiteBold.Printf(" %s %s", jsonquery.FindOne(dev, "Manufacturer").InnerText(), jsonquery.FindOne(dev, "Product").InnerText())

		serial := jsonquery.FindOne(dev, "Serial").InnerText()
		spec := jsonquery.FindOne(dev, "Spec").InnerText()
		speed := jsonquery.FindOne(dev, "Speed").InnerText()

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
			norm.Printf(" %s", jsonquery.FindOne(fn, "Power").InnerText())
			norm.Print(" [")
			if jsonquery.FindOne(fn, "Wakeup").InnerText() == "true" {
				norm.Printf("Wakeup")
			}
			norm.Print("]")

			ifaces := jsonquery.Find(fn, "Interfaces/*")
			if len(ifaces) != 1 {
				norm.Println()
			}
			for _, iface := range ifaces {
				norm.Print("    ")
				norm.Print(iface.InnerText())
				norm.Println()
			}
		}
	}
}

func renderBlock(root *jsonquery.Node) {
	for _, blk := range jsonquery.Find(root, "Hardware/Block/*") {
		white.Print(blk.Data)
		whiteBold.Printf(" %s %s", jsonquery.FindOne(blk, "Vendor").InnerText(), jsonquery.FindOne(blk, "Model").InnerText())

		serial := jsonquery.FindOne(blk, "Serial").InnerText()
		if serial != "" {
			grey.Printf(" serial %s", serial)
		}

		norm.Printf(" [%s, %s bytes", jsonquery.FindOne(blk, "ControllerType").InnerText(), jsonquery.FindOne(blk, "SizeBytes").InnerText())
		if jsonquery.FindOne(blk, "Removable").InnerText() == "true" {
			norm.Printf("Removable")
		}
		norm.Print("]")

		norm.Println()

		ps := jsonquery.Find(blk, "Partitions/*")
		for _, p := range ps {
			norm.Print("  ")
			white.Print(p.Data)
			whiteBold.Printf(" %s on %s", jsonquery.FindOne(p, "Filesystem").InnerText(), jsonquery.FindOne(p, "MountPoint").InnerText())
			grey.Printf(" uuid %s", jsonquery.FindOne(p, "UUID").InnerText())
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
