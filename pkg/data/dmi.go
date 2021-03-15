package data

import (
	"context"
	"strconv"

	"github.com/dselans/dmidecode"
	"github.com/go-logr/logr"
)

func init() {
	plugins = append(plugins, getDmiData)
}

func getDmiData(ctx context.Context, log logr.Logger, t *Trie) {
	dmi := dmidecode.New()
	// TODO: try to detect when it's a permissions error (even if we just check our UID or some /sys file access), and set Forbidden
	if err := dmi.Run(); err != nil {
		log.Error(err, "Can't read DMI")
		return
	}

	syss, _ := dmi.SearchByType(1) // system info
	if len(syss) != 1 {
		log.Info("Unexpectedly many DMI readings; skipping")
		return
	}
	sys := syss[0]
	t.Insert(sys["Manufacturer"], "Hardware", "System", "Manufacturer")
	t.Insert(sys["Product Name"], "Hardware", "System", "Product")

	mbs, _ := dmi.SearchByType(2) // mobo
	if len(mbs) != 1 {
		log.Info("Unexpectedly many DMI readings; skipping")
		return
	}
	mb := mbs[0]
	t.Insert(mb["Manufacturer"], "Hardware", "Motherboard", "Manufacturer")
	t.Insert(mb["Product Name"], "Hardware", "Motherboard", "Product")

	fws, _ := dmi.SearchByType(0) // firmware
	if len(fws) != 1 {
		log.Info("Unexpectedly many DMI readings; skipping")
		return
	}
	fw := fws[0]
	t.Insert(fw["Manufacturer"], "Hardware", "Firmware", "Manufacturer")
	t.Insert(fw["Product Name"], "Hardware", "Firmware", "Version")
	t.Insert(fw["BIOS Revision"], "Hardware", "Firmware", "Revision")
	t.Insert(fw["Release Date"], "Hardware", "Firmware", "Date")
	t.Insert(fw["ROM Size"], "Hardware", "Firmware", "ROM Size")

	cpus, _ := dmi.SearchByType(4) // CPUs
	for i, cpu := range cpus {
		t.Insert(cpu["Socket Designation"], "Hardware", "CPU", strconv.Itoa(i), "Socket")
		t.Insert(cpu["Max Speed"], "Hardware", "CPU", strconv.Itoa(i), "Max Speed")
	}

	dimms, _ := dmi.SearchByType(17) // DIMMs
	for i, dimm := range dimms {
		t.Insert(dimm["Bank Locator"], "Hardware", "Memory", strconv.Itoa(i), "Channel")
		t.Insert(dimm["Locator"], "Hardware", "Memory", strconv.Itoa(i), "Slot")
		t.Insert(dimm["Rank"], "Hardware", "Memory", strconv.Itoa(i), "Ranks")
		t.Insert(dimm["Size"], "Hardware", "Memory", strconv.Itoa(i), "Size")
		t.Insert(dimm["Speed"], "Hardware", "Memory", strconv.Itoa(i), "Speed")
		t.Insert(dimm["Type"], "Hardware", "Memory", strconv.Itoa(i), "Type")
		t.Insert(dimm["Type Detail"], "Hardware", "Memory", strconv.Itoa(i), "Sub Type")
	}

	// for _, record := range dmi.Data {
	// 	for _, v := range record {
	// 		spew.Dump(v)
	// 	}
	// }
}
