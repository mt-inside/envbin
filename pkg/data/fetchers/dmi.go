package fetchers

import (
	"context"
	"strconv"

	"github.com/dselans/dmidecode"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getDmiData)
}

func getDmiData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	dmi := dmidecode.New()
	// TODO: try to detect when it's a permissions error (even if we just check our UID or some /sys file access), and set Forbidden
	if err := dmi.Run(); err != nil {
		vals <- Insert(Error(err), "Hardware", "DMI")
		log.Error(err, "Can't read DMI")
		return
	}

	getDmiSystem(log, dmi, vals)

	getDmiMobo(log, dmi, vals)

	getDmiFw(log, dmi, vals)

	cpus, _ := dmi.SearchByType(4) // CPUs
	for i, cpu := range cpus {
		vals <- Insert(Some(cpu["Socket Designation"]), "Hardware", "CPU", strconv.Itoa(i), "Socket")
		vals <- Insert(Some(cpu["Max Speed"]), "Hardware", "CPU", strconv.Itoa(i), "Max Speed")
	}

	dimms, _ := dmi.SearchByType(17) // DIMMs
	for i, dimm := range dimms {
		vals <- Insert(Some(dimm["Bank Locator"]), "Hardware", "Memory", strconv.Itoa(i), "Channel")
		vals <- Insert(Some(dimm["Locator"]), "Hardware", "Memory", strconv.Itoa(i), "Slot")
		vals <- Insert(Some(dimm["Rank"]), "Hardware", "Memory", strconv.Itoa(i), "Ranks")
		vals <- Insert(Some(dimm["Size"]), "Hardware", "Memory", strconv.Itoa(i), "Size")
		vals <- Insert(Some(dimm["Speed"]), "Hardware", "Memory", strconv.Itoa(i), "Speed")
		vals <- Insert(Some(dimm["Type"]), "Hardware", "Memory", strconv.Itoa(i), "Type")
		vals <- Insert(Some(dimm["Type Detail"]), "Hardware", "Memory", strconv.Itoa(i), "Sub Type")
	}

	// for _, record := range dmi.Data {
	// 	for _, v := range record {
	// 		spew.Dump(v)
	// 	}
	// }
}

func getDmiSystem(log logr.Logger, dmi *dmidecode.DMI, vals chan<- InsertMsg) {
	syss, _ := dmi.SearchByType(1) // system info
	if len(syss) != 1 {
		log.Info("Unexpectedly many DMI 'system info' entries; skipping all")
		return
	}
	sys := syss[0]
	vals <- Insert(Some(sys["Manufacturer"]), "Hardware", "System", "Manufacturer")
	vals <- Insert(Some(sys["Product Name"]), "Hardware", "System", "Product")
}

func getDmiMobo(log logr.Logger, dmi *dmidecode.DMI, vals chan<- InsertMsg) {
	mbs, _ := dmi.SearchByType(2) // mobo
	if len(mbs) != 1 {
		log.Info("Unexpectedly many DMI 'motherboard' entries; skipping all")
		return
	}
	mb := mbs[0]
	vals <- Insert(Some(mb["Manufacturer"]), "Hardware", "Motherboard", "Manufacturer")
	vals <- Insert(Some(mb["Product Name"]), "Hardware", "Motherboard", "Product")
}

func getDmiFw(log logr.Logger, dmi *dmidecode.DMI, vals chan<- InsertMsg) {
	fws, _ := dmi.SearchByType(0) // firmware
	if len(fws) != 1 {
		log.Info("Unexpectedly many DMI 'firmware' entries; skipping all")
		return
	}
	fw := fws[0]
	vals <- Insert(Some(fw["Manufacturer"]), "Hardware", "Firmware", "Manufacturer")
	vals <- Insert(Some(fw["Product Name"]), "Hardware", "Firmware", "Version")
	vals <- Insert(Some(fw["BIOS Revision"]), "Hardware", "Firmware", "Revision")
	vals <- Insert(Some(fw["Release Date"]), "Hardware", "Firmware", "Date")
	vals <- Insert(Some(fw["ROM Size"]), "Hardware", "Firmware", "ROM Size")
}
