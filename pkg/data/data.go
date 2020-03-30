package data

import "github.com/mt-inside/envbin/pkg/util"

func GetData() map[string]string {
	d := make(map[string]string) //TODO: strongly type me with a struct. Esp for (optional) sections

	d = util.AppendMap(d, getSessionData())
	d = util.AppendMap(d, getNetworkData())
	d = util.AppendMap(d, getHardwareData())
	//d = util.AppendMap(d, getProcData()) TODO: darwin (or optional)
	d = util.AppendMap(d, getProcsData())
	d = util.AppendMap(d, getMemData())
	d = util.AppendMap(d, getOsData())
	d = util.AppendMap(d, getSettingsData())
	d = util.AppendMap(d, getK8sData())

	return d
}
