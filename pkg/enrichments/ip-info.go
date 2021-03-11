package enrichments

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-logr/logr"
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

func ExternalIp(ctx context.Context, log logr.Logger) (string, error) {
	// TODO Can force to v4?

	info, err := ipApiFetch(ctx, log, "")
	if err != nil {
		return "", err
	}

	return info.Ip, nil
}

func EnrichIpRendered(ctx context.Context, log logr.Logger, ip string) (string, error) {
	if ip == "" {
		panic("Empty IP is a special parameter to ipapi.co (gets apparent external IP) and shouldn't be provided through this path")
	}

	info, err := ipApiFetch(ctx, log, ip)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s, %s, %s, %s (AS: %s, %s)", info.City, info.Region, info.Postal, info.Country, info.Asn, info.As), nil
}

func ipApiFetch(ctx context.Context, log logr.Logger, ip string) (IpInfo, error) {
	client := http.Client{
		Timeout: time.Second * 2,
	}

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
