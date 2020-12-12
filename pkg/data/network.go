package data

import (
	"log"
	"net"
	"os"

	"github.com/mt-inside/envbin/pkg/enrichments"
)

func getNetworkData() map[string]string {
	data := map[string]string{}

	hostname, _ := os.Hostname()

	data["Hostname"] = hostname
	data["HostIp"] = getDefaultIP()
	// basically pointless enriching the Host IP; either it's a container or VM or private network, in which case enrichment is pointless, or the host's interface has the external IP, in which case it'll equal ExternalIp, which is enriched anyway
	data["ExternalIp"] = getExternalIp()
	data["ExternalIpEnrich"] = enrichments.EnrichIpRendered(data["ExternalIp"])

	return data
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
