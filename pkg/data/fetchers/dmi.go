package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/yumaojun03/dmidecode"
	"github.com/yumaojun03/dmidecode/parser/processor"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/enrichments"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getDmiFirmware)
	data.RegisterPlugin(getDmiMotherboard)
	data.RegisterPlugin(getDmiSystem)
	data.RegisterPlugin(getDmiRAM)
	data.RegisterPlugin(getDmiCPU)
	data.RegisterPlugin(getDmiCPUCache)
}

func u8(n byte) string {
	return strconv.FormatInt(int64(n), 10)
}
func u16(n uint16) string {
	return strconv.FormatInt(int64(n), 10)
}

func getDmiFirmware(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.BIOS()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "Firmware")
		return
	}

	vals <- trie.Insert(trie.Some(data[0].Vendor), "Hardware", "Firmware", "Vendor")
	vals <- trie.Insert(trie.Some(data[0].BIOSVersion), "Hardware", "Firmware", "Version")
	vals <- trie.Insert(trie.Some(u8(data[0].SystemBIOSMajorRelease)), "Hardware", "Firmware", "Version", "Major")
	vals <- trie.Insert(trie.Some(u8(data[0].SystemBIOSMinorRelease)), "Hardware", "Firmware", "Version", "Minor")
	vals <- trie.Insert(trie.Some(data[0].ReleaseDate), "Hardware", "Firmware", "Date")
}

func getDmiMotherboard(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.BaseBoard()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "Motherboard")
		return
	}

	if len(data) == 0 {
		return
	}

	vals <- trie.Insert(trie.Some(data[0].Manufacturer), "Hardware", "Motherboard", "Vendor")
	vals <- trie.Insert(trie.Some(data[0].ProductName), "Hardware", "Motherboard", "Model")
	vals <- trie.Insert(trie.Some(data[0].Version), "Hardware", "Motherboard", "Version")
	vals <- trie.Insert(trie.Some(data[0].SerialNumber), "Hardware", "Motherboard", "Serial")
}

func getDmiSystem(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "DMI")
		return
	}

	chassis, err := dmi.Chassis()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "Chassis")
		return
	}

	vals <- trie.Insert(trie.Some(chassis[0].Manufacturer), "Hardware", "Chassis", "Vendor")
	vals <- trie.Insert(trie.Some(chassis[0].SKUNumber), "Hardware", "Chassis", "SKU")
	vals <- trie.Insert(trie.Some(chassis[0].Version), "Hardware", "Chassis", "Version")
	vals <- trie.Insert(trie.Some(chassis[0].SerialNumber), "Hardware", "Chassis", "Serial")

	system, err := dmi.System()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "System")
		return
	}

	vals <- trie.Insert(trie.Some(system[0].Manufacturer), "Hardware", "System", "Vendor")
	vals <- trie.Insert(trie.Some(system[0].Family), "Hardware", "System", "Family")
	vals <- trie.Insert(trie.Some(system[0].ProductName), "Hardware", "System", "Model")
	vals <- trie.Insert(trie.Some(system[0].SKUNumber), "Hardware", "System", "SKU")
	vals <- trie.Insert(trie.Some(system[0].Version), "Hardware", "System", "Version")
	vals <- trie.Insert(trie.Some(system[0].SerialNumber), "Hardware", "System", "Serial")
}

func getDmiRAM(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.MemoryDevice()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "RAM")
		return
	}

	for i, m := range data {
		// Product
		vals <- trie.Insert(trie.Some(m.Manufacturer), "Hardware", "RAM", strconv.Itoa(i), "Vendor")
		vals <- trie.Insert(trie.Some(m.PartNumber), "Hardware", "RAM", strconv.Itoa(i), "SKU")
		vals <- trie.Insert(trie.Some(m.SerialNumber), "Hardware", "RAM", strconv.Itoa(i), "Serial")
		vals <- trie.Insert(trie.Some(m.FormFactor.String()), "Hardware", "RAM", strconv.Itoa(i), "Package")

		// Location
		vals <- trie.Insert(trie.Some(m.BankLocator), "Hardware", "RAM", strconv.Itoa(i), "Channel")
		vals <- trie.Insert(trie.Some(m.DeviceLocator), "Hardware", "RAM", strconv.Itoa(i), "Slot")

		// Bus
		vals <- trie.Insert(trie.Some(u16(m.Speed)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Speed", "Max MT/s")
		vals <- trie.Insert(trie.Some(u16(m.ConfiguredMemoryClockSpeed)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Speed", "Current MT/s")
		vals <- trie.Insert(trie.Some(u16(m.TotalWidth)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Width", "Total")
		vals <- trie.Insert(trie.Some(u16(m.DataWidth)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Width", "Data") // Will be less if ECC

		// Module
		vals <- trie.Insert(trie.Some(m.Type.String()), "Hardware", "RAM", strconv.Itoa(i), "Type")
		vals <- trie.Insert(trie.Some(m.TypeDetail.String()), "Hardware", "RAM", strconv.Itoa(i), "Subtype")
		vals <- trie.Insert(trie.Some(u16(m.Size)), "Hardware", "RAM", strconv.Itoa(i), "SizeMB")
		vals <- trie.Insert(trie.Some(u16(m.MinimumVoltage)), "Hardware", "RAM", strconv.Itoa(i), "Voltage", "MinimummV")
		vals <- trie.Insert(trie.Some(u16(m.ConfiguredVoltage)), "Hardware", "RAM", strconv.Itoa(i), "Voltage", "CurrentmV")
		vals <- trie.Insert(trie.Some(u16(m.MaximumVoltage)), "Hardware", "RAM", strconv.Itoa(i), "Voltage", "MaximummV")

		// DMI-reported "speed" is bus's MT/s (double bus frequency for DDR)
		enrichments.EnrichRamSpecs(ctx, log, m.Type, uint(m.ConfiguredMemoryClockSpeed), uint(m.DataWidth), trie.PrefixChan(vals, "Hardware", "RAM", strconv.Itoa(i)))
	}
}

func getDmiCPU(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.Processor()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "CPU")
		return
	}

	vals <- trie.Insert(trie.Some(data[0].Manufacturer), "Hardware", "CPU", "Vendor")
	vals <- trie.Insert(trie.Some(data[0].Family.String()), "Hardware", "CPU", "Family")
	vals <- trie.Insert(trie.Some(data[0].Version), "Hardware", "CPU", "Model")
	vals <- trie.Insert(trie.Some(data[0].SerialNumber), "Hardware", "CPU", "Serial")
	vals <- trie.Insert(trie.Some(data[0].SocketDesignation), "Hardware", "CPU", "Package")

	vals <- trie.Insert(trie.Some(u8(data[0].CoreCount)), "Hardware", "CPU", "Cores")
	vals <- trie.Insert(trie.Some(u8(data[0].ThreadCount)), "Hardware", "CPU", "Threads")
	vals <- trie.Insert(trie.Some(data[0].Voltage.String()), "Hardware", "CPU", "Voltage")
	vals <- trie.Insert(trie.Some(u16(data[0].ExternalClock)), "Hardware", "CPU", "Clock", "Bus")
	vals <- trie.Insert(trie.Some(u16(data[0].CurrentSpeed)), "Hardware", "CPU", "Clock", "Current")
	vals <- trie.Insert(trie.Some(u16(data[0].MaxSpeed)), "Hardware", "CPU", "Clock", "Max")
}

func cacheSizeK(c processor.CacheSize) string {
	grans := [...]int64{
		1024,
		65536,
	}
	return strconv.FormatInt(grans[c.Granularity]*int64(c.Size), 10)
}
func getDmiCPUCache(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.ProcessorCache()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "CPU", "Cache", "Totals")
		return
	}

	for _, c := range data {
		level := c.Configuration.Level.String()
		vals <- trie.Insert(trie.Some(cacheSizeK(c.InstalledSize)), "Hardware", "CPU", "Cache", "Totals", level)
		//vals <- trie.Insert(trie.Some(c.CacheSpeed)), "Hardware", "CPU", "Cache", "Totals", level)
	}
}
