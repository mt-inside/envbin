package data

import "os"

func getFirmwareData() map[string]string {
	data := map[string]string{}

	_, err := os.Stat("/sys/firmware/efi")
	if os.IsNotExist(err) {
		data["FirmwareType"] = "BIOS"
	} else {
		data["FirmwareType"] = "EFI"
	}

	return data
}