package fetchers

import (
	"context"
	"fmt"
	"math"
	"net"
	"os"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/safchain/ethtool"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/enrichments"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("network", getNetworkData)
}

func getNetworkData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	hostname, _ := os.Hostname()

	vals <- trie.Insert(trie.Some(hostname), "Network", "Hostname")

	getIfaces(log, vals)

	getDefaultIP(log, vals)

	enrichments.EnrichedExternalIp(ctx, log, trie.PrefixChan(vals, "Network", "ExternalIP"))
}

func getIfaces(log logr.Logger, vals chan<- trie.InsertMsg) {
	ifaces, err := net.Interfaces()
	if err != nil {
		vals <- trie.Insert(trie.Error(fmt.Errorf("can't get network interfaces: %w", err)), "Network", "Interfaces")
		return
	}

	e, err := ethtool.NewEthtool()
	if err != nil {
		vals <- trie.Insert(trie.Error(fmt.Errorf("can't construct ethtool link: %w", err)), "Network", "Interfaces")
	}
	defer e.Close()

	for _, iface := range ifaces {
		k := strconv.Itoa(iface.Index)
		vals <- trie.Insert(trie.Some(iface.Name), "Network", "Interfaces", k, "Name")
		vals <- trie.Insert(trie.Some(iface.Flags.String()), "Network", "Interfaces", k, "IfaceFlags")

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.To4() != nil {
				vals <- trie.Insert(trie.Some(addr.String()), "Network", "Interfaces", k, "IPv4Address")
				break
			}
		}

		// features, err := e.Features(iface.Name)
		// if err == nil {
		// 	fmt.Printf("features: %+v\n", features)
		// }

		// stats, err := e.Stats(iface.Name)
		// if err == nil {
		// 	fmt.Printf("stats: %+v\n", stats)
		// }

		drvInfo, err := e.DriverInfo(iface.Name)
		if err == nil {
			//fmt.Printf("drvrinfo: %+v\n", drvInfo)
			// TODO: use this to insert into the unified topology tree at the right place
			vals <- trie.Insert(trie.Some(drvInfo.BusInfo), "Network", "Interfaces", k, "BusAddr")
			vals <- trie.Insert(trie.Some(drvInfo.Driver), "Network", "Interfaces", k, "Driver")
			vals <- trie.Insert(trie.Some(drvInfo.FwVersion), "Network", "Interfaces", k, "Firmware")
		}

		cmdGet := ethtool.EthtoolCmd{}
		speed, err := e.CmdGet(&cmdGet, iface.Name)
		if err == nil {
			//fmt.Printf("cmd get: %+v\n", cmdGet)

			if speed != math.MaxUint32 {
				vals <- trie.Insert(trie.Some(strconv.FormatUint(uint64(speed), 10)), "Network", "Interfaces", k, "SpeedMbits")
			}

			if drvInfo.BusInfo != "N/A" {
				vals <- trie.Insert(trie.Some(strconv.FormatUint(uint64(cmdGet.Phy_address), 10)), "Network", "Interfaces", k, "PHY")
				vals <- trie.Insert(trie.Some(strconv.FormatUint(uint64(cmdGet.Transceiver), 10)), "Network", "Interfaces", k, "Transceiver")
				vals <- trie.Insert(trie.Some(strconv.FormatUint(uint64(cmdGet.Port), 10)), "Network", "Interfaces", k, "Port")
			}
		}

		linkState, err := e.LinkState(iface.Name)
		if err == nil {
			vals <- trie.Insert(trie.Some(strconv.FormatBool(linkState == 1)), "Network", "Interfaces", k, "PhysicalLink")
		}

		permAddr, err := e.PermAddr(iface.Name)
		if err == nil {
			vals <- trie.Insert(trie.Some(permAddr), "Network", "Interfaces", k, "MACAddr")
		}
	}
}

func getDefaultIP(log logr.Logger, vals chan<- trie.InsertMsg) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		vals <- trie.Insert(trie.Error(fmt.Errorf("can't get default IP: %w", err)), "Network", "DefaultIP")
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	vals <- trie.Insert(trie.Some(localAddr.IP.String()), "Network", "DefaultIP")
}
