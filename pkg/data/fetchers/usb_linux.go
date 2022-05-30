//go:build native
// +build native

package fetchers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getUsbData)
}

func unwrap(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func orElse(s string, err error) string {
	if err != nil {
		return err.trie.Error()
	}
	return s
}

func getUsbData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	prefix := []string{"Hardware", "Bus", "USB"}

	err := usbid.LoadFromURL(usbid.LinuxUsbDotOrg)
	if err != nil {
		log.Error(err, "Can't load USB IDs")
		// Don't return, will just have worse device naming
	}

	usb := gousb.NewContext()
	defer usb.Close()

	devs, err := usb.OpenDevices(func(d *gousb.DeviceDesc) bool {
		return true
	})
	if err != nil {
		vals <- trie.Insert(trie.Error(err), prefix...)
		return
	}

	for _, dev := range devs {
		defer dev.Close()
		d := dev.Desc

		phyAddr, err := findPhysicalAddr(d.Bus, d.Address)
		if err != nil {
			continue // The virtual Root Hub devices, which are annoying
		}

		addr := fmt.Sprintf("%d-%s", d.Bus, phyAddr)

		vals <- trie.Insert(trie.Some(strconv.Itoa(d.Bus)), "Hardware", "Bus", "USB", addr, "Bus")
		vals <- trie.Insert(trie.Some(phyAddr), "Hardware", "Bus", "USB", addr, "Physical Address")
		vals <- trie.Insert(trie.Some(strconv.Itoa(d.Address)), "Hardware", "Bus", "USB", addr, "Logical Device")

		var m, p string
		mx, ok := usbid.Vendors[d.Vendor]
		if ok {
			m = mx.Name
			px, ok := mx.Product[d.Product]
			if ok {
				p = px.Name
			} else {
				p = orElse(dev.Product())
			}
		} else {
			m = unwrap(dev.Manufacturer())
			p = unwrap(dev.Product())
		}

		vals <- trie.Insert(trie.Some(d.Vendor.String()), "Hardware", "Bus", "USB", addr, "VendorID")
		vals <- trie.Insert(trie.Some(d.Product.String()), "Hardware", "Bus", "USB", addr, "ProductID")
		vals <- trie.Insert(trie.Some(m), "Hardware", "Bus", "USB", addr, "Manufacturer")
		vals <- trie.Insert(trie.Some(p), "Hardware", "Bus", "USB", addr, "Product")
		vals <- trie.Insert(trie.Some(orElse(dev.SerialNumber())), "Hardware", "Bus", "USB", addr, "Serial")
		vals <- trie.Insert(trie.Some(d.Spec.String()), "Hardware", "Bus", "USB", addr, "Spec")
		vals <- trie.Insert(trie.Some(d.Speed.String()), "Hardware", "Bus", "USB", addr, "Speed")

		for _, c := range d.Configs {
			pow := "0mA"
			if c.SelfPowered {
				pow = fmt.Sprintf("%dmA", c.MaxPower)
			}

			vals <- trie.Insert(trie.Some(pow), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Power")
			vals <- trie.Insert(trie.Some(strconv.FormatBool(c.RemoteWakeup)), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Wakeup")
			for _, i := range c.Interfaces {
				// I've never seen a device with differing properties across alts of an interface, so we just read item 0. If you do need to iterate Alts, nb that you want a.Alternate, not a.Number
				vals <- trie.Insert(trie.Some(usbid.Classify(i.AltSettings[0])), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Interfaces", strconv.Itoa(i.Number), "Description")
				driver, err := findDriver(d.Bus, phyAddr, c.Number, i.Number)
				if err == nil {
					vals <- trie.Insert(trie.Some(driver), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Interfaces", strconv.Itoa(i.Number), "Driver")
				} else {
					vals <- trie.Insert(trie.Error(err), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Interfaces", strconv.Itoa(i.Number), "Driver")
				}
			}
		}
	}
}

func findPhysicalAddr(bus int, dev int) (string, error) {
	// filename format: .../devices/<bus>-<address>[.<port>[.<port>...]]:<configuration>.<interface>
	glob := fmt.Sprintf("/sys/bus/usb/devices/%d-*", bus)
	devPaths, err := filepath.Glob(glob)
	if err != nil {
		return "", err
	}

	for _, devPath := range devPaths {
		if strings.Contains(devPath, ":") {
			continue
		}

		devnumRaw, err := os.ReadFile(filepath.Join(devPath, "devnum"))
		if err != nil {
			return "", err
		}
		devnum, err := strconv.Atoi(strings.TrimSpace(string(devnumRaw)))
		if err != nil {
			return "", err
		}

		if devnum != dev {
			continue
		}

		devpathRaw, err := os.ReadFile(filepath.Join(devPath, "devpath"))
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(devpathRaw)), nil
	}

	return "", fmt.Errorf("Can't find USB device at bus %d addr %d", bus, dev)
}

func findDriver(bus int, phyAddr string, config int, iface int) (string, error) {
	ifacePath := fmt.Sprintf("/sys/bus/usb/devices/%d-%s:%d.%d", bus, phyAddr, config, iface)
	driverPath := filepath.Join(ifacePath, "driver")
	targetPath, err := os.Readlink(driverPath)
	if err != nil {
		return "", err
	}
	driverName := filepath.Base(targetPath)

	return driverName, nil
}
