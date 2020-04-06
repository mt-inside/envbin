package actions

import (
	"github.com/mt-inside/envbin/pkg/data"
	"math"
	"runtime"
	"time"
)

func init() {
	useCPU()
}

// Not the best "algorithm" in the world.
// * Seems to over/undershoot by about 10%. This could be the sampling rate of top though
// * You really don't need to do anything "cpu-intensive" here; this happily loads 16 virtual cores.
func useCPU() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			period := time.Tick(1 * time.Second)
			for {
				cpus := float64(runtime.NumCPU())
				// Try to cap the high-time to rought 1s. If it's more than that, then the high duty part lasts longer than the 1s period timer, and that channel starts to fill up with ticks. As soon as there's some breathing room again that'll quickly get drained, but the channel could get full and idk what the Tick producer does then.
				// Also, if the user requests a crazy duty cycle of say 1M, then it won't respond to requests to lower that rate until after a period of 1M / #cores.
				dutyCycle := math.Min(data.GetCPUUse(), cpus) / cpus
				highTimer := time.After(
					time.Duration(
						dutyCycle*1000,
					) * time.Millisecond,
				)
			high:
				for {
					select {
					case <-highTimer:
						break high
					default:
					}
				}
				<-period // low
			}
		}()
	}
}
