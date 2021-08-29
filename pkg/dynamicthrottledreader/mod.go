package dynamicthrottledreader

import (
	"fmt"
	"io"

	"github.com/mxk/go-flowrate/flowrate"
	"github.com/spf13/viper"
)

type DynamicThrottledReader struct {
	tr       *flowrate.Reader
	oldLimit int64
}

func NewDynamicThrottledReader(r io.Reader) *DynamicThrottledReader {
	limit := viper.GetInt64("Rate")
	return &DynamicThrottledReader{
		// TODO take init arg
		tr:       flowrate.NewReader(r, limit),
		oldLimit: limit,
	}
}

func (d *DynamicThrottledReader) Read(p []byte) (int, error) {
	acc := 0
	for i := 0; i < len(p); i++ {
		// TODO move to channel, non-block read every time
		if d.oldLimit != viper.GetInt64("Rate") {
			d.tr.SetLimit(viper.GetInt64("Rate"))
			d.oldLimit = viper.GetInt64("Rate")
		}

		n, err := d.tr.Read(p[i : i+1])
		if n > 1 {
			panic(fmt.Errorf("Read too much: %d", n))
		}
		acc += n

		if err != nil || n == 0 {
			return acc, err
		}
	}
	return acc, nil
}

func (d *DynamicThrottledReader) Close() error {
	return d.tr.Close()
}

func (d *DynamicThrottledReader) SetBlocking(new bool) (old bool) {
	return d.tr.SetBlocking(new)
}
