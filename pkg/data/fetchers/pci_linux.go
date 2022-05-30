package fetchers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	"github.com/jaypipes/ghw"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getPciData)
}

func getPciData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	prefix := []string{"Hardware", "Bus", "PCI"}

	pci, err := ghw.PCI()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), prefix...)
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

		vals <- trie.Insert(trie.Some(domain), "Hardware", "Bus", "PCI", addr, "Domain")
		vals <- trie.Insert(trie.Some(bus), "Hardware", "Bus", "PCI", addr, "Bus")
		vals <- trie.Insert(trie.Some(device), "Hardware", "Bus", "PCI", addr, "Device")
		vals <- trie.Insert(trie.Some(d.Driver), "Hardware", "Bus", "PCI", addr, "Functions", function, "Driver")
		vals <- trie.Insert(trie.Some(d.Vendor.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Vendor")
		vals <- trie.Insert(trie.Some(d.Product.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Product")
		vals <- trie.Insert(trie.Some(d.Revision), "Hardware", "Bus", "PCI", addr, "Functions", function, "Revision")
		vals <- trie.Insert(trie.Some(d.Class.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Class")
		vals <- trie.Insert(trie.Some(d.Subclass.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Subclass")
		vals <- trie.Insert(trie.Some(d.Driver), "Hardware", "Bus", "PCI", addr, "Functions", function, "Driver")
		for _, iface := range d.Subclass.ProgrammingInterfaces {
			vals <- trie.Insert(trie.Some(iface.Name), "Hardware", "Bus", "PCI", addr, "Functions", function, "Interfaces", iface.ID)
		}
	}
}
