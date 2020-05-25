package data

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/docker/go-units"
)

var (
	Version   string
	GitCommit string
	BuildTime string
)

var (
	name      string
	startTime time.Time
	reqNo     int
)

func init() {
	name = namesgenerator.GetRandomName(0)
	startTime = time.Now()
	reqNo = 0
}

func getSessionData() map[string]string {
	data := map[string]string{}

	reqNo++ // TODO: not the right place for this. Go MVC

	data["Version"] = Version
	data["GitCommit"] = GitCommit
	data["BuildTime"] = BuildTime
	data["StartTime"] = startTime.Format("2006-01-02 15:04:05 -0700")
	data["RunTime"] = units.HumanDuration(time.Now().Sub(startTime))
	data["SessionName"] = name
	data["RequestNumber"] = strconv.Itoa(reqNo)

	return data
}

func RenderSessionData() (ret []string) {
	ret = append(ret, fmt.Sprintf("envbin %s: git %s, built at %s with %s", Version, GitCommit, BuildTime, runtime.Version()))
	ret = append(ret, fmt.Sprintf("session: %s", name))

	return
}
