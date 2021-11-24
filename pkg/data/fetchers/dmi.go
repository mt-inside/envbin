package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/yumaojun03/dmidecode"
	"github.com/yumaojun03/dmidecode/parser/processor"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
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

func getDmiFirmware(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		log.Error(err, "Can't read DMI")
		vals <- Insert(Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.BIOS()
	if err != nil {
		log.Error(err, "Can't get DMI Firmware info")
		vals <- Insert(Error(err), "Hardware", "Firmware")
		return
	}

	vals <- Insert(Some(data[0].Vendor), "Hardware", "Firmware", "Vendor")
	vals <- Insert(Some(data[0].BIOSVersion), "Hardware", "Firmware", "Version")
	vals <- Insert(Some(data[0].ReleaseDate), "Hardware", "Firmware", "Date")
	vals <- Insert(Some(u8(data[0].SystemBIOSMajorRelease)), "Hardware", "Firmware", "Major")
	vals <- Insert(Some(u8(data[0].SystemBIOSMinorRelease)), "Hardware", "Firmware", "Minor")
}

func getDmiMotherboard(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		log.Error(err, "Can't read DMI")
		vals <- Insert(Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.BaseBoard()
	if err != nil {
		log.Error(err, "Can't get DMI Motherboard info")
		vals <- Insert(Error(err), "Hardware", "Motherboard")
		return
	}

	if len(data) == 0 {
		return
	}

	vals <- Insert(Some(data[0].Manufacturer), "Hardware", "Motherboard", "Vendor")
	vals <- Insert(Some(data[0].ProductName), "Hardware", "Motherboard", "Product")
	vals <- Insert(Some(data[0].Version), "Hardware", "Motherboard", "Version")
	vals <- Insert(Some(data[0].SerialNumber), "Hardware", "Motherboard", "Serial")
}

func getDmiSystem(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		log.Error(err, "Can't read DMI")
		vals <- Insert(Error(err), "Hardware", "DMI")
		return
	}

	chassis, err := dmi.Chassis()
	if err != nil {
		log.Error(err, "Can't get DMI Chassis info")
		vals <- Insert(Error(err), "Hardware", "Chassis")
		return
	}

	vals <- Insert(Some(chassis[0].Manufacturer), "Hardware", "Chassis", "Vendor")
	vals <- Insert(Some(chassis[0].SKUNumber), "Hardware", "Chassis", "Product")
	vals <- Insert(Some(chassis[0].Version), "Hardware", "Chassis", "Version")
	vals <- Insert(Some(chassis[0].SerialNumber), "Hardware", "Chassis", "Serial")

	system, err := dmi.System()
	if err != nil {
		log.Error(err, "Can't get DMI Chassis info")
		vals <- Insert(Error(err), "Hardware", "System")
		return
	}

	vals <- Insert(Some(system[0].Manufacturer), "Hardware", "System", "Vendor")
	vals <- Insert(Some(system[0].ProductName), "Hardware", "System", "Product")
	vals <- Insert(Some(system[0].Family), "Hardware", "System", "Family")
	vals <- Insert(Some(system[0].SKUNumber), "Hardware", "System", "SKU")
	vals <- Insert(Some(system[0].Version), "Hardware", "System", "Version")
	vals <- Insert(Some(system[0].SerialNumber), "Hardware", "System", "Serial")
}

func getDmiRAM(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		log.Error(err, "Can't read DMI")
		vals <- Insert(Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.MemoryDevice()
	if err != nil {
		log.Error(err, "Can't get DMI RAM info")
		vals <- Insert(Error(err), "Hardware", "RAM")
		return
	}

	for i, m := range data {
		vals <- Insert(Some(m.BankLocator), "Hardware", "RAM", strconv.Itoa(i), "Channel")
		vals <- Insert(Some(m.DeviceLocator), "Hardware", "RAM", strconv.Itoa(i), "Slot")
		vals <- Insert(Some(m.Manufacturer), "Hardware", "RAM", strconv.Itoa(i), "Vendor")
		vals <- Insert(Some(m.PartNumber), "Hardware", "RAM", strconv.Itoa(i), "Product")
		vals <- Insert(Some(m.SerialNumber), "Hardware", "RAM", strconv.Itoa(i), "Serial")
		vals <- Insert(Some(m.Type.String()), "Hardware", "RAM", strconv.Itoa(i), "Type")
		vals <- Insert(Some(m.TypeDetail.String()), "Hardware", "RAM", strconv.Itoa(i), "Subtype")
		vals <- Insert(Some(m.FormFactor.String()), "Hardware", "RAM", strconv.Itoa(i), "Form Factor")
		vals <- Insert(Some(u16(m.Size)), "Hardware", "RAM", strconv.Itoa(i), "Size MB")
		vals <- Insert(Some(u16(m.Speed)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Speed", "Max MT/s")
		vals <- Insert(Some(u16(m.ConfiguredMemoryClockSpeed)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Speed", "Current MT/s")
		vals <- Insert(Some(u16(m.MinimumVoltage)), "Hardware", "RAM", strconv.Itoa(i), "Voltage", "Minimum mV")
		vals <- Insert(Some(u16(m.ConfiguredVoltage)), "Hardware", "RAM", strconv.Itoa(i), "Voltage", "Current mV")
		vals <- Insert(Some(u16(m.MaximumVoltage)), "Hardware", "RAM", strconv.Itoa(i), "Voltage", "Maximum mV")
		vals <- Insert(Some(u16(m.TotalWidth)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Width", "Total")
		vals <- Insert(Some(u16(m.DataWidth)), "Hardware", "RAM", strconv.Itoa(i), "Bus", "Width", "Data") // Will be less if ECC
	}
}

func getDmiCPU(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		log.Error(err, "Can't read DMI")
		vals <- Insert(Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.Processor()
	if err != nil {
		log.Error(err, "Can't get DMI CPU info")
		vals <- Insert(Error(err), "Hardware", "CPU")
		return
	}

	vals <- Insert(Some(data[0].Manufacturer), "Hardware", "CPU", "Vendor")
	vals <- Insert(Some(data[0].Version), "Hardware", "CPU", "Product")
	vals <- Insert(Some(data[0].Family.String()), "Hardware", "CPU", "Family")
	vals <- Insert(Some(data[0].SerialNumber), "Hardware", "CPU", "Serial")
	vals <- Insert(Some(u8(data[0].CoreCount)), "Hardware", "CPU", "Cores")
	vals <- Insert(Some(u8(data[0].ThreadCount)), "Hardware", "CPU", "Threads")
	vals <- Insert(Some(data[0].Voltage.String()), "Hardware", "CPU", "Voltage")
	vals <- Insert(Some(u16(data[0].ExternalClock)), "Hardware", "CPU", "Clock", "Bus")
	vals <- Insert(Some(u16(data[0].CurrentSpeed)), "Hardware", "CPU", "Clock", "Current")
	vals <- Insert(Some(u16(data[0].MaxSpeed)), "Hardware", "CPU", "Clock", "Max")
}

func cacheSizeK(c processor.CacheSize) string {
	grans := [...]int64{
		1024,
		65536,
	}
	return strconv.FormatInt(grans[c.Granularity]*int64(c.Size), 10)
}
func getDmiCPUCache(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	dmi, err := dmidecode.New()
	if err != nil {
		log.Error(err, "Can't read DMI")
		vals <- Insert(Error(err), "Hardware", "DMI")
		return
	}

	data, err := dmi.ProcessorCache()
	if err != nil {
		log.Error(err, "Can't get DMI CPU Cache info")
		vals <- Insert(Error(err), "Hardware", "CPU", "Cache")
		return
	}

	for _, c := range data {
		vals <- Insert(Some(cacheSizeK(c.InstalledSize)), "Hardware", "CPU", "Cache", "Totals", c.SocketDesignation)
	}
}
