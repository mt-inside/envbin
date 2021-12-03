package fetchers

import (
	"context"
	"net"
	"os"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/enrichments"

	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getNetworkData)
}

func getNetworkData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	hostname, _ := os.Hostname()

	vals <- Insert(Some(hostname), "Network", "Hostname")

	getIfaces(log, vals)

	getDefaultIP(log, vals)

	enrichments.EnrichedExternalIp(ctx, log, PrefixChan(vals, "Network", "ExternalIP"))
}

func getIfaces(log logr.Logger, vals chan<- InsertMsg) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Error(err, "Can't show system network interfaces")
		vals <- Insert(Error(err), "Network", "Interfaces")
		return
	}

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.To4() == nil {
				continue
			}
			k := strconv.Itoa(iface.Index)
			vals <- Insert(Some(iface.Name), "Network", "Interfaces", k, "Name")
			vals <- Insert(Some(addr.String()), "Network", "Interfaces", k, "Address")
			vals <- Insert(Some(iface.Flags.String()), "Network", "Interfaces", k, "Flags")
		}
	}
}

func getDefaultIP(log logr.Logger, vals chan<- InsertMsg) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Error(err, "Can't get default IP")
		vals <- Insert(Error(err), "Network", "DefaultIP")
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	vals <- Insert(Some(localAddr.IP.String()), "Network", "DefaultIP")
}
