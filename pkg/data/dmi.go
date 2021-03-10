package data

import (
	"github.com/dselans/dmidecode"
	"github.com/mt-inside/go-usvc"
)

func init() {
	plugins = append(plugins, getDmiData)
}

func getDmiData() map[string]string {
	log := usvc.Global

	data := map[string]string{}

	dmi := dmidecode.New()
	if err := dmi.Run(); err != nil {
		usvc.Global.Error(err, "Can't read DMI")
		return map[string]string{}
	}

	syss, _ := dmi.SearchByType(1) // system info
	for _, sys := range syss {
		log.Info(
			"System",
			"Manufacturer", sys["Manufacturer"],
			"Product", sys["Product Name"],
		)
	}

	mbs, _ := dmi.SearchByType(2) // mobo
	for _, mb := range mbs {
		log.Info(
			"Baseboard",
			"Manufacturer", mb["Manufacturer"],
			"Product", mb["Product Name"],
		)
	}

	fws, _ := dmi.SearchByType(0) // firmware
	for _, fw := range fws {
		log.Info(
			"Firmware",
			"Manufacturer", fw["Vendor"],
			"Version", fw["Version"],
			"Revision", fw["BIOS Revision"],
			"Date", fw["Release Date"],
			"ROM Size", fw["ROM Size"],
		)
	}

	cpus, _ := dmi.SearchByType(4) // CPUs
	for _, cpu := range cpus {
		log.Info(
			"CPU",
			"Socket", cpu["Socket Designation"],
			"Max Speed", cpu["Max Speed"],
			"Cores", cpu["Core Count"],
			"Threads", cpu["Thread Count"],
		)
	}

	dimms, _ := dmi.SearchByType(17) // DIMMs
	for _, dimm := range dimms {
		log.Info(
			"RAM",
			"Channel", dimm["Bank Locator"],
			"Slot", dimm["Locator"],
			"Ranks", dimm["Rank"],
			"Size", dimm["Size"],
			"Speed", dimm["Speed"],
			"Type", dimm["Type"],
			"Sub Type", dimm["Type Detail"],
		)
	}

	// for _, record := range dmi.Data {
	// 	for _, v := range record {
	// 		spew.Dump(v)
	// 	}
	// }

	return data
}
