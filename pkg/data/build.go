package data

import (
	"fmt"
	"runtime"
)

var (
	Version   string
	GitCommit string
	BuildTime string
)

func getSessionData() map[string]string {
	data := map[string]string{}

	data["Version"] = Version
	data["GitCommit"] = GitCommit
	data["BuildTime"] = BuildTime

	return data
}

func RenderSessionData() (ret []string) {
	ret = append(ret, fmt.Sprintf("envbin %s: git %s, built at %s with %s", Version, GitCommit, BuildTime, runtime.Version()))

	return
}
