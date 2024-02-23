//go:build native && ((linux && amd64) || darwin)

package fetchers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("usb", getUsbData)
}

// For unified topology:
// - walk /sys/bus/usb/devices/usb*
// - serial is the PCI address
// - [this]/[this.devnum]-0:* are the configs & interfaces of the hub - useful; do this
// - then [this]/[this.devnum]-[^0] - the direct ports
// - then recurse
// - this lib we're using does useful stuff looking up device IDs, decoding configs etc - see if it can target/find a specific device, if not, have a func that just loops them and matches physical path
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

	log.Info("Enumerated USB devices", "count", len(devs))
	for _, dev := range devs {
		defer dev.Close()
		d := dev.Desc

		phyAddr, err := findUSBPhysicalAddr(d.Bus, d.Address)
		if err != nil {
			continue // The virtual Root Hub devices, which are annoying
		}

		addr := []string{"Hardware", "Bus", "USB", "Busses", strconv.Itoa(d.Bus)}
		for i, a := range strings.Split(phyAddr, ".") {
			if i == 0 {
				addr = append(addr, "Address", a)
				continue
			}
			addr = append(addr, "Port", a)
		}
		node := trie.PrefixChan(vals, addr...)

		node <- trie.Insert(trie.Some(strconv.Itoa(d.Bus)), "Bus")
		node <- trie.Insert(trie.Some(phyAddr), "Physical Address")
		node <- trie.Insert(trie.Some(strconv.Itoa(d.Address)), "Logical Device")

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

		node <- trie.Insert(trie.Some(d.Vendor.String()), "VendorID")
		node <- trie.Insert(trie.Some(d.Product.String()), "ProductID")
		node <- trie.Insert(trie.Some(m), "Manufacturer")
		node <- trie.Insert(trie.Some(p), "Product")
		node <- trie.Insert(trie.Some(orElse(dev.SerialNumber())), "Serial")
		node <- trie.Insert(trie.Some(d.Spec.String()), "Spec")
		node <- trie.Insert(trie.Some(d.Speed.String()), "Speed")

		for _, c := range d.Configs {
			pow := "0mA"
			if c.SelfPowered {
				pow = fmt.Sprintf("%dmA", c.MaxPower)
			}

			node <- trie.Insert(trie.Some(pow), "Configs", strconv.Itoa(c.Number), "Power")
			node <- trie.Insert(trie.Some(strconv.FormatBool(c.RemoteWakeup)), "Configs", strconv.Itoa(c.Number), "Wakeup")
			for _, i := range c.Interfaces {
				// I've never seen a device with differing properties across alts of an interface, so we just read item 0. If you do need to iterate Alts, nb that you want a.Alternate, not a.Number
				node <- trie.Insert(trie.Some(usbid.Classify(i.AltSettings[0])), "Configs", strconv.Itoa(c.Number), "Interfaces", strconv.Itoa(i.Number), "Description")
				driver, err := findUSBDriver(d.Bus, phyAddr, c.Number, i.Number)
				if err == nil {
					node <- trie.Insert(trie.Some(driver), "Configs", strconv.Itoa(c.Number), "Interfaces", strconv.Itoa(i.Number), "Driver")
				} else {
					node <- trie.Insert(trie.Error(err), "Configs", strconv.Itoa(c.Number), "Interfaces", strconv.Itoa(i.Number), "Driver")
				}
			}
		}
	}
}
