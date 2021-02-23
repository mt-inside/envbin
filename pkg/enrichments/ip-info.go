package enrichments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mt-inside/envbin/pkg/util"
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

	info, err := ipApiFetch("")
	if err != nil {
		util.GlobalLog.Error(err, "Can't get info for external IP")
		return "<unknown>"
	}

	return info.Ip
}

func EnrichIpRendered(ip string) string {
	if ip == "" {
		log.Printf("IP '%s' is a special parameter to ipapi.co (gets apparent external IP) and shouldn't be provided through this path", ip)
		return ""
	}

	info, err := ipApiFetch(ip)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s, %s, %s, %s (AS: %s, %s)", info.City, info.Region, info.Postal, info.Country, info.Asn, info.As)
}

func ipApiFetch(ip string) (IpInfo, error) {
	log := util.GlobalLog // TODO hack

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/json", baseUrl, ip), nil)
	if err != nil {
		log.Error(err, "Can't make http request?")
		return IpInfo{}, err
	}

	req.Header.Set("user-agent", "envbin")

	res, err := client.Do(req)
	if err != nil {
		log.Error(err, "Can't get IP info", "ip", ip)
		return IpInfo{}, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err, "Can't get IP info", "ip", ip)
		return IpInfo{}, err
	}

	ipInfo := IpInfo{}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		log.Error(err, "Can't get IP info", "ip", ip)
		return IpInfo{}, err
	}

	// TODO ipapi.co returns valid JSON for a lot of error cases, just with "error" and "reason" set
	if ipInfo.Error {
		log.Error(err, "Can't get IP info", "ip", ip, "message", ipInfo.Reason)
		return IpInfo{}, err
	}

	return ipInfo, nil
}
