package fetchers

import "fmt"

// returns address[.port[.port[...]]]
func findUSBPhysicalAddr(bus int, dev int) (string, error) {
	return fmt.Sprintf("%d:%d", bus, dev), nil
}

func findUSBDriver(bus int, phyAddr string, config int, iface int) (string, error) {
	return "unknown on Darwin", nil
}
