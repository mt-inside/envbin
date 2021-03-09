package data

import (
	"net/http"

	"github.com/mt-inside/envbin/pkg/util"
)

func GetData(r *http.Request) map[string]string {
	d := make(map[string]string) //TODO: strongly type me with a struct. Esp for (optional) sections

	d = util.AppendMap(d, getBuildData())
	d = util.AppendMap(d, getNetworkData())
	d = util.AppendMap(d, getHardwareData())
	d = util.AppendMap(d, getDmiData())
	d = util.AppendMap(d, getProcData()) //TODO: darwin (or optional)
	d = util.AppendMap(d, getProcsData())
	d = util.AppendMap(d, getMemData())
	d = util.AppendMap(d, getOsData())
	//d = util.AppendMap(d, getK8sData()) //TODO: handle no permissions and other errors
	d = util.AppendMap(d, getRequestData(r))
	d = util.AppendMap(d, getTerminalData())
	d = util.AppendMap(d, getFirmwareData())
	d = util.AppendMap(d, getOsDistributionData())

	return d
}
