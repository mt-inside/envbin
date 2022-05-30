package fetchers

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/enrichments"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getNetworkData)
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

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.To4() == nil {
				continue
			}
			k := strconv.Itoa(iface.Index)
			vals <- trie.Insert(trie.Some(iface.Name), "Network", "Interfaces", k, "Name")
			vals <- trie.Insert(trie.Some(addr.String()), "Network", "Interfaces", k, "Address")
			vals <- trie.Insert(trie.Some(iface.Flags.String()), "Network", "Interfaces", k, "Flags")
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
