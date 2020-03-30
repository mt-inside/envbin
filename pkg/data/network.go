package data

import (
	"log"
	"net"
	"os"
)

func getNetworkData() map[string]string {
	data := map[string]string{}

	hostname, _ := os.Hostname()

	data["Hostname"] = hostname
	data["Ip"] = getDefaultIP()

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