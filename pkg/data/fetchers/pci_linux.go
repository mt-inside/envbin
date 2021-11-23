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

func getPciData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	prefix := []string{"Hardware", "Bus", "PCI"}

	pci, err := ghw.PCI()
	if err != nil {
		log.Error(err, "Can't read PCI data")
		vals <- Insert(Error(err), prefix...)
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

		vals <- Insert(Some(domain), "Hardware", "Bus", "PCI", addr, "Domain")
		vals <- Insert(Some(bus), "Hardware", "Bus", "PCI", addr, "Bus")
		vals <- Insert(Some(device), "Hardware", "Bus", "PCI", addr, "Device")
		vals <- Insert(Some(d.Driver), "Hardware", "Bus", "PCI", addr, "Functions", function, "Driver")
		vals <- Insert(Some(d.Vendor.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Vendor")
		vals <- Insert(Some(d.Product.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Product")
		vals <- Insert(Some(d.Revision), "Hardware", "Bus", "PCI", addr, "Functions", function, "Revision")
		vals <- Insert(Some(d.Class.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Class")
		vals <- Insert(Some(d.Subclass.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Subclass")
		vals <- Insert(Some(d.Driver), "Hardware", "Bus", "PCI", addr, "Functions", function, "Driver")
		for _, iface := range d.Subclass.ProgrammingInterfaces {
			vals <- Insert(Some(iface.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Interfaces", iface.ID)
		}
	}
}
