package data

import (
	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/go-units"
	"strconv"
	"time"
)

var (
	version string
)

var (
	name string
	startTime time.Time
	reqNo     int
)

func init() {
name =      namesgenerator.GetRandomName(0)
	startTime = time.Now()
	reqNo = 0
}

func getSessionData() map[string]string {
	data := map[string]string{}

	reqNo++ // TODO: not the right place for this. Go MVC

	data["Version"] = version
	data["StartTime"] = startTime.Format("2006-01-02 15:04:05")
	data["RunTime"] = units.HumanDuration(time.Now().Sub(startTime))
	data["SessionName"] = name
	data["RequestNumber"] = strconv.Itoa(reqNo)

	return data
}