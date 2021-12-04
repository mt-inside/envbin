package enrichments

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/go-logr/logr"

	. "github.com/mt-inside/envbin/pkg/data/trie"
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

func EnrichedExternalIp(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	// TODO Can force to v4?

	info, err := ipApiFetch(ctx, log, "")
	if err != nil {
		log.Error(err, "Can't get external IP and its info")
		vals <- Insert(Error(fmt.Errorf("Can't get external IP info from ipapi.co: %w", err)), "Details")
		return
	}

	vals <- Insert(Some(info.Ip), "Address")
	enrichFromInfo(info, vals)
	reverseDNS(info.Ip, vals)
}

func EnrichIp(ctx context.Context, log logr.Logger, ip string, vals chan<- InsertMsg) {
	if ip == "" {
		panic("Empty IP is a special parameter to ipapi.co (gets apparent external IP) and shouldn't be provided through this path")
	}

	info, err := ipApiFetch(ctx, log, ip)
	if err != nil {
		log.Error(err, "Can't get IP info", "ip", ip)
		vals <- Insert(Error(fmt.Errorf("Can't get IP info from ipapi.co: %w", err)), "Details")
		return
	}

	enrichFromInfo(info, vals)
	reverseDNS(info.Ip, vals)
}

func enrichFromInfo(info IpInfo, vals chan<- InsertMsg) {
	vals <- Insert(Some(info.City), "City")
	vals <- Insert(Some(info.Region), "Region")
	vals <- Insert(Some(info.Postal), "Postal") // postal code
	vals <- Insert(Some(info.Country), "Country")
	vals <- Insert(Some(info.As), "AS")
	vals <- Insert(Some(info.Asn), "ASN")
}

func reverseDNS(ip string, vals chan<- InsertMsg) {
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		vals <- Insert(Error(err), "ReverseDNS")
	}

	vals <- Insert(Some(strings.Join(hosts, ",")), "ReverseDNS")
}

func ipApiFetch(ctx context.Context, log logr.Logger, ip string) (IpInfo, error) {
	client := http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s/json", baseUrl, ip), nil)
	if err != nil {
		return IpInfo{}, err
	}

	req.Header.Set("user-agent", "envbin")

	res, err := client.Do(req)
	if err != nil {
		return IpInfo{}, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return IpInfo{}, err
	}

	ipInfo := IpInfo{}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return IpInfo{}, err
	}

	// TODO ipapi.co returns valid JSON for a lot of error cases, just with "error" and "reason" set
	if ipInfo.Error {
		return IpInfo{}, fmt.Errorf("IpInfo error: %s", ipInfo.Reason)
	}

	return ipInfo, nil
}
