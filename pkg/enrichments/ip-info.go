package enrichments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const baseUrl = "https://ipapi.co"

type IpInfo struct {
	Ip      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country_name"`
	Postal  string `json:"postal"`
	Asn     string `json:"asn"`
	As      string `json:"org"`
	Error   bool   `json:"error"`
	Reason  string `json:"reason"`
}

func ExternalIp() string {
	// TODO Can force to v4?
	return ipApiFetch("").Ip
}

func EnrichIp(ip string) IpInfo {
	if ip == "" {
		log.Fatalf("IP '%s' is a special parameter to ipapi.co and shouldn't be provided through this path", ip)
	}
	return ipApiFetch(ip)
}

func EnrichIpRendered(ip string) string {
	e := EnrichIp(ip)
	return fmt.Sprintf("%s, %s, %s, %s (AS: %s, %s)", e.City, e.Region, e.Postal, e.Country, e.Asn, e.As)
}

func ipApiFetch(ip string) IpInfo {
	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/json", baseUrl, ip), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("user-agent", "envbin")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Can't get IP info for %s: %v", ip, err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Can't get IP info for %s: %v", ip, err)
	}

	ipInfo := IpInfo{}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		log.Printf("Can't get IP info for %s: %v", ip, err)
	}

	// TODO ipapi.co returns valid JSON for a lot of error cases, just with "error" and "reason" set
	if ipInfo.Error {
		log.Printf("Can't get IP info for %s: %v", ip, ipInfo.Reason)
	}

	return ipInfo
}
