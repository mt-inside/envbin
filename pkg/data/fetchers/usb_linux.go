package fetchers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// returns address[.port[.port[...]]]
func findUSBPhysicalAddr(bus int, dev int) (string, error) {
	// filename format: .../devices/<bus>-<address>[.<port>[.<port>...]]:<configuration>.<interface>
	glob := fmt.Sprintf("/sys/bus/usb/devices/%d-*", bus)
	devPaths, err := filepath.Glob(glob)
	if err != nil {
		return "", err
	}

	for _, devPath := range devPaths {
		if strings.Contains(devPath, ":") {
			continue
		}

		devnumRaw, err := os.ReadFile(filepath.Join(devPath, "devnum"))
		if err != nil {
			return "", err
		}
		devnum, err := strconv.Atoi(strings.TrimSpace(string(devnumRaw)))
		if err != nil {
			return "", err
		}

		if devnum != dev {
			continue
		}

		devpathRaw, err := os.ReadFile(filepath.Join(devPath, "devpath"))
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(devpathRaw)), nil
	}

	return "", fmt.Errorf("can't find USB device at bus %d addr %d", bus, dev)
}

func findUSBDriver(bus int, phyAddr string, config int, iface int) (string, error) {
	ifacePath := fmt.Sprintf("/sys/bus/usb/devices/%d-%s:%d.%d", bus, phyAddr, config, iface)
	driverPath := filepath.Join(ifacePath, "driver")
	targetPath, err := os.Readlink(driverPath)
	if err != nil {
		return "", err
	}
	driverName := filepath.Base(targetPath)

	return driverName, nil
}
