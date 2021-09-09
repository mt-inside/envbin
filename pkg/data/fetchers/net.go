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

func getNetworkData(ctx context.Context, log logr.Logger, t *Trie) {
	hostname, _ := os.Hostname()

	t.Insert(Some(hostname), "Network", "Hostname")

	getIfaces(log, t)

	getDefaultIP(log, t)

	extIP, err := enrichments.ExternalIp(ctx, log)
	if err != nil {
		log.Error(err, "Can't get external IP address")
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			t.Insert(Timeout(time.Second), "Network", "ExternalIP")
		} else {
			t.Insert(Error(err), "Network", "ExternalIP")
		}
		return
	}

	t.Insert(Some(extIP), "Network", "ExternalIP", "Address")

	extIpInfo, err := enrichments.EnrichIpRendered(ctx, log, extIP)
	if err != nil {
		log.Error(err, "Can't get IP info", "ip", extIP)
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			t.Insert(Timeout(time.Second), "Network", "ExternalIP", "Info")
		} else {
			t.Insert(Error(err), "Network", "ExternalIP", "Info")
		}
		return
	}

	t.Insert(Some(extIpInfo), "Network", "ExternalIP", "Info")
}

func getIfaces(log logr.Logger, t *Trie) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Error(err, "Can't show system network interfaces")
		t.Insert(Error(err), "Network", "Interfaces")
		return
	}

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.To4() == nil {
				continue
			}
			k := strconv.Itoa(iface.Index)
			t.Insert(Some(iface.Name), "Network", "Interfaces", k, "Name")
			t.Insert(Some(addr.String()), "Network", "Interfaces", k, "Address")
			t.Insert(Some(iface.Flags.String()), "Network", "Interfaces", k, "Flags")
		}
	}
}

func getDefaultIP(log logr.Logger, t *Trie) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Error(err, "Can't get default IP")
		t.Insert(Error(err), "Network", "DefaultIP")
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	t.Insert(Some(localAddr.IP.String()), "Network", "DefaultIP")
}
