package data

import (
	"github.com/docker/go-units"
	"strconv"
)

type Settings struct {
	delay     int64
	bandwidth int64
	errorRate float64
	cpuUse    float64
	liveness  bool
	readiness bool
}

func NewSettings() *Settings {
	return &Settings{
		delay:     0,
		bandwidth: 0,
		errorRate: 0.0,
		cpuUse:    0.0,
		liveness:  true,
		readiness: true,
	}
}

var (
	settings *Settings
)

func init() {
	settings = NewSettings()
}

func getSettingsData() map[string]string {
	data := map[string]string{}

	data["SettingLiveness"] = strconv.FormatBool(settings.liveness)
	data["SettingReadiness"] = strconv.FormatBool(settings.readiness)
	data["SettingLatency"] = strconv.Itoa(int(settings.delay))
	data["SettingBandwidth"] = units.BytesSize(float64(settings.bandwidth))
	data["SettingErrorRate"] = strconv.FormatFloat(settings.errorRate, 'f', 2, 64)
	data["SettingCpuUse"] = strconv.FormatFloat(settings.cpuUse, 'f', 2, 64)

	return data
}

// TODO make methods
func SetDelay(d int64) {
	settings.delay = d
}

func SetBandwidth(b int64) {
	settings.bandwidth = b
}

func SetErrorRate(e float64) {
	settings.errorRate = e
}

func SetCPUUse(c float64) {
	settings.cpuUse = c
}

func SetLiveness(l bool) {
	settings.liveness = l
}

func SetReadiness(r bool)  {
	settings.readiness = r
}


func GetDelay() int64 {
	return settings.delay
}

func GetBandwidth() int64 {
	return settings.bandwidth
}

func GetErrorRate() float64 {
	return settings.errorRate
}

func GetCPUUse() float64 {
	return settings.cpuUse
}

func GetLiveness() bool {
	return settings.liveness
}

func GetReadiness() bool {
	return settings.readiness
}