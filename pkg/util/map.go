package util

func AppendMap(base map[string]string, extra map[string]string) map[string]string {
	for k, v := range extra {
		base[k] = v
	}

	return base
}
