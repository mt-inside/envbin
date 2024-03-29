package main

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

func main() {
	var err error

	var adapter = bluetooth.DefaultAdapter
	err = adapter.Enable()
	if err != nil {
		panic(err)
	}

	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		fmt.Println("found device:", device.Address.String(), device.RSSI, device.LocalName())
	})
	if err != nil {
		panic(err)
	}
}
