package data

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/enrichments"
)

func init() {
	plugins = append(plugins, getNetworkData)
}

func getNetworkData(ctx context.Context, log logr.Logger, t *Trie) {
	hostname, _ := os.Hostname()

	t.Insert(Some{hostname}, "Network", "Hostname")

	getIfaces(t)

	t.Insert(Some{getDefaultIP()}, "Network", "DefaultIP")

	extIP, err := enrichments.ExternalIp(ctx, log)
	if err != nil {
		log.Error(err, "Can't get external IP address")
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			t.Insert(Some{"Timeout"}, "Network", "ExternalIP")
		} else {
			t.Insert(Some{"Error"}, "Network", "ExternalIP")
		}
	} else {
		t.Insert(Some{extIP}, "Network", "ExternalIP", "Address")

		extIpInfo, err := enrichments.EnrichIpRendered(ctx, log, extIP)
		if err != nil {
			log.Error(err, "Can't get IP info", "ip", extIP)
			if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
				t.Insert(Some{"Timeout"}, "Network", "ExternalIP", "Info")
			} else {
				t.Insert(Some{"Error"}, "Network", "ExternalIP", "Info")
			}
		} else {
			t.Insert(Some{extIpInfo}, "Network", "ExternalIP", "Info")
		}
	}
}

func getIfaces(t *Trie) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal("Can't show system network interfaces")
		return
	}

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if addr.(*net.IPNet).IP.To4() == nil {
				continue
			}
			k := fmt.Sprint(iface.Index)
			v := fmt.Sprintf("%s, %s, %s", iface.Name, addr.String(), iface.Flags)
			t.Insert(Some{v}, "Network", "Interfaces", string(k))
		}
	}
}

func getDefaultIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Println(err)
		return "<unknown>"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
