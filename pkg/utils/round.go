package utils

import (
	"fmt"
	"math"
)

func Round(x int64, base int64, sigfig int64) int64 {
	size := int64(math.Floor(math.Log(float64(x)) / math.Log(float64(base))))
	lump_size := size - sigfig + 1
	lump := math.Pow(float64(base), float64(lump_size))
	frac := float64(x) / lump
	sigfigs := math.Round(frac)
	apprx := int64(sigfigs * lump)

	return int64(apprx)
}

func FormatIEC(x int64, decimals int64) string {
	if x == 0 {
		return "0"
	}

	base := math.Log(float64(x)) / math.Log(1024)
	suffixes := []string{"", "ki", "Mi", "Gi", "Ti"}

	return fmt.Sprintf("%.[1]*f%s", decimals, float64(x)/math.Pow(1024, math.Floor(base)), suffixes[int(math.Floor(base))])
}

func FormatSI(x int64, decimals int64) string {
	if x == 0 {
		return "0"
	}

	base := math.Log(float64(x)) / math.Log(1000)
	suffixes := []string{"", "k", "M", "G", "T"}

	return fmt.Sprintf("%.[1]*f%s", decimals, float64(x)/math.Pow(1000, math.Floor(base)), suffixes[int(math.Floor(base))])
}
