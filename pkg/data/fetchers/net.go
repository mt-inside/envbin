package fetchers

import (
	"context"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

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

	extIP, err := enrichments.ExternalIp(ctx, log)
	if err != nil {
		log.Error(err, "Can't get external IP address")
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			vals <- Insert(Timeout(time.Second), "Network", "ExternalIP")
		} else {
			vals <- Insert(Error(err), "Network", "ExternalIP")
		}
		return
	}

	enrichments.EnrichIp(ctx, log, extIP, PrefixChan(vals, "Network", "ExternalIP"))
	vals <- Insert(Some(extIP), "Network", "ExternalIP", "Address")
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
