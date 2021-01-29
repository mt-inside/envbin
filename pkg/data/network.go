package data

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mt-inside/envbin/pkg/enrichments"
)

func getNetworkData() map[string]string {
	data := map[string]string{}

	hostname, _ := os.Hostname()

	data["Hostname"] = hostname
	getIfaces(data)
	data["HostIp"] = getDefaultIP()
	// basically pointless enriching the Host IP; either it's a container or VM or private network, in which case enrichment is pointless, or the host's interface has the external IP, in which case it'll equal ExternalIp, which is enriched anyway
	data["ExternalIp"] = getExternalIp()
	data["ExternalIpEnrich"] = enrichments.EnrichIpRendered(data["ExternalIp"])

	return data
}

func getIfaces(data map[string]string) {
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
			k := fmt.Sprintf("Interface%d", iface.Index)
			v := fmt.Sprintf("%s, %s, %s", iface.Name, addr.String(), iface.Flags)
			data[k] = v
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

func getExternalIp() string {
	// Using the same service to get this too
	return enrichments.ExternalIp()
}
