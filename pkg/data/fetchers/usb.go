package fetchers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/google/gousb"
	"github.com/google/gousb/usbid"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
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

func getUsbData(ctx context.Context, log logr.Logger, t *Trie) {
	prefix := []string{"Hardware", "Bus", "USB"}

	err := usbid.LoadFromURL(usbid.LinuxUsbDotOrg)
	if err != nil {
		t.Insert(Error(err), prefix...)
	}

	usb := gousb.NewContext()
	defer usb.Close()

	devs, err := usb.OpenDevices(func(d *gousb.DeviceDesc) bool {
		return true
	})
	if err != nil {
		t.Insert(Error(err), prefix...)
	}

	for _, dev := range devs {
		defer dev.Close()
		d := dev.Desc

		addr := fmt.Sprintf("%03d:%03d:%d", d.Bus, d.Address, d.Port)

		var m, p string
		mx, ok := usbid.Vendors[d.Vendor]
		if ok {
			m = mx.Name
			px, ok := mx.Product[d.Product]
			if ok {
				p = px.Name
			} else {
				p = unwrap(dev.Product())
			}
		} else {
			m = unwrap(dev.Manufacturer())
			p = unwrap(dev.Product())
		}

		t.Insert(Some(d.Vendor.String()), "Hardware", "Bus", "USB", addr, "VendorID")
		t.Insert(Some(d.Product.String()), "Hardware", "Bus", "USB", addr, "ProductID")
		t.Insert(Some(m), "Hardware", "Bus", "USB", addr, "Manufacturer")
		t.Insert(Some(p), "Hardware", "Bus", "USB", addr, "Product")
		t.Insert(Some(unwrap(dev.SerialNumber())), "Hardware", "Bus", "USB", addr, "Serial")
		t.Insert(Some(d.Spec.String()), "Hardware", "Bus", "USB", addr, "Spec")
		t.Insert(Some(d.Speed.String()), "Hardware", "Bus", "USB", addr, "Speed")

		for _, c := range d.Configs {
			pow := "0mA"
			if c.SelfPowered {
				pow = fmt.Sprintf("%dmA", c.MaxPower)
			}

			t.Insert(Some(pow), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Power")
			t.Insert(Some(strconv.FormatBool(c.RemoteWakeup)), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Wakeup")
			for _, i := range c.Interfaces {
				// I've never seen a device with differing properties across alts of an interface, so we just read item 0. If you do need to iterate Alts, nb that you want a.Alternate, not a.Number
				t.Insert(Some(usbid.Classify(i.AltSettings[0])), "Hardware", "Bus", "USB", addr, "Configs", strconv.Itoa(c.Number), "Interfaces", strconv.Itoa(i.Number))
			}
		}
	}
}