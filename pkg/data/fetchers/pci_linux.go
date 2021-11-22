package fetchers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/jaypipes/ghw"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getPciData)
}

func getPciData(ctx context.Context, log logr.Logger, t *Trie) {
	prefix := []string{"Hardware", "Bus", "PCI"}

	pci, err := ghw.PCI()
	if err != nil {
		log.Error(err, "Can't read PCI data")
		t.Insert(Error(err), prefix...)
		return
	}

	for _, d := range pci.Devices {
		ss := strings.Split(d.Address, ".")
		addr := ss[0]
		function := ss[1]
		addrs := strings.Split(addr, ":")
		domain := addrs[0]
		bus := addrs[1]
		device := addrs[2]

		t.Insert(Some(domain), "Hardware", "Bus", "PCI", addr, "Domain")
		t.Insert(Some(bus), "Hardware", "Bus", "PCI", addr, "Bus")
		t.Insert(Some(device), "Hardware", "Bus", "PCI", addr, "Device")
		t.Insert(Some(d.Driver), "Hardware", "Bus", "PCI", addr, "Functions", function, "Driver")
		t.Insert(Some(d.Vendor.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Vendor")
		t.Insert(Some(d.Product.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Product")
		t.Insert(Some(d.Revision), "Hardware", "Bus", "PCI", addr, "Functions", function, "Revision")
		t.Insert(Some(d.Class.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Class")
		t.Insert(Some(d.Subclass.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Subclass")
		t.Insert(Some(d.Driver), "Hardware", "Bus", "PCI", addr, "Functions", function, "Driver")
		for _, iface := range d.Subclass.ProgrammingInterfaces {
			t.Insert(Some(iface.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Interfaces", iface.ID)
		}
	}
}
