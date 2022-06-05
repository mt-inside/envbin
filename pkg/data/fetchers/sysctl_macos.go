package fetchers

import (
	"context"
	"encoding/json"
	"os/exec"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/enrichments"
	"github.com/mt-inside/envbin/pkg/data/trie"
	"howett.net/plist"
)

func init() {
	data.RegisterPlugin(getIoreg)
	data.RegisterPlugin(getSPHw)
	data.RegisterPlugin(getSPSw)
}

/* Ways to get more info:
* - system_profiler will give you all you need and loads more. Dunno what API(s) it uses under the hood - could strace it)
*   - The internet seems to thing that running it will lead to it writing its info in ~/Library/Prefs/com.apple.SystemProfiler.plist but doesn't seem to happen
* - IOKitLib, callable from C (eg golang FFI: https://gist.github.com/csexton/56121dbb613df68f143162b60a2c694a)
*   - `ioreg -l` dump some stuff to the terminal - all of ^^ ? (serial number is in there, along with model like from sysctl, nothing else)
 */

// Data also available in sysctl, but everything in there is duplicated by one of the others
// func getSysctl(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
// 	// https://reincubate.com/support/deviceidentifier/apple-identifiers/#understanding-codes
// 	// https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man3/sysctl.3.html

// 	vals <- trie.Insert(trie.Some("Apple"), "Hardware", "System", "Vendor")

// 	s, err := syscall.Sysctl("hw.model")
// 	if err != nil {
// 		vals <- trie.Insert(trie.Error(err), "Hardware", "System", "Product")
// 		return
// 	}

// 	vals <- trie.Insert(trie.Some(s), "Hardware", "System", "Product")
// }

func getIoreg(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	ioregOut, err := exec.Command(
		"ioreg", // Can't find a proper API for this. TODO: strace what this command calls
		"-a",    // output as XML plist "archive", rather than pretty-printing
		//"-l", // Show all properties
		"-c",                     // Output only objects with class name...
		"IOPlatformExpertDevice", // <<
	).Output()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "System")
	}

	var ioreg map[string]interface{}
	_, err = plist.Unmarshal(ioregOut, &ioreg)
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "System")
	}

	ioregPlatExpDev := ioreg["IORegistryEntryChildren"].([]interface{})[0].(map[string]interface{})

	vals <- trie.Insert(trie.Some(ioregPlatExpDev["IORegistryEntryName"].(string)), "Hardware", "System", "SKU") // (ie key of this object) Model?
}

func getSPHw(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	spHwOut, err := exec.Command(
		"system_profiler", // Can't find a proper API for this. TODO: strace what this command calls
		"SPHardwareDataType",
		"-json",
	).Output()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "System")
		return
	}

	var spHw map[string]interface{}
	err = json.Unmarshal(spHwOut, &spHw)
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "System")
		return
	}

	spHwOverview := spHw["SPHardwareDataType"].([]interface{})[0].(map[string]interface{})

	vals <- trie.Insert(trie.Some("Apple"), "Hardware", "System", "Vendor")
	vals <- trie.Insert(trie.Some(spHwOverview["machine_name"].(string)), "Hardware", "System", "Product")
	vals <- trie.Insert(trie.Some(spHwOverview["machine_model"].(string)), "Hardware", "System", "Family")
	vals <- trie.Insert(trie.Some(spHwOverview["serial_number"].(string)), "Hardware", "System", "Serial")
	vals <- trie.Insert(trie.Some(spHwOverview["platform_UUID"].(string)), "Hardware", "System", "UUID")

	vals <- trie.Insert(trie.Some(spHwOverview["chip_type"].(string)), "Hardware", "CPU", "Product")
	enrichments.EnrichMacProcs(ctx, log, spHwOverview["number_processors"].(string), trie.PrefixChan(vals, "Hardware", "CPU"))

	vals <- trie.Insert(trie.Some(spHwOverview["boot_rom_version"].(string)), "Hardware", "Firmware", "Version")
}

func getSPSw(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	// Also available at SPInstallHistoryDataType[@._name=="macOS *"]

	spSwOut, err := exec.Command(
		"system_profiler", // Can't find a proper API for this. TODO: strace what this command calls
		"SPSoftwareDataType",
		"-json",
	).Output()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "OS", "Distro")
	}

	var spSw map[string]interface{}
	err = json.Unmarshal(spSwOut, &spSw)
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "OS", "Distro")
	}

	spSwOSOverview := spSw["SPSoftwareDataType"].([]interface{})[0].(map[string]interface{})

	vals <- trie.Insert(trie.Some(spSwOSOverview["os_version"].(string)), "OS", "Distro")
}
