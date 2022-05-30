package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/reiver/go-v4l2"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getV4l2Data)
}

func getV4l2Data(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	for i := 0; i < 64; i++ {
		nodeName := "/dev/video" + strconv.Itoa(i)
		device, err := v4l2.Open(nodeName)
		if nil != err {
			continue
		}
		defer device.Close()

		vals <- trie.Insert(trie.Some(device.MustCard()), "Hardware", "V4l2", nodeName, "Name")
		vals <- trie.Insert(trie.Some(device.MustDriver()), "Hardware", "V4l2", nodeName, "Driver")
		vals <- trie.Insert(trie.Some(device.MustBusInfo()), "Hardware", "V4l2", nodeName, "Location")
		vals <- trie.Insert(trie.Some(device.MustVersion()), "Hardware", "V4l2", nodeName, "Version")
		vals <- trie.Insert(trie.Some(strconv.FormatBool(device.MustHasCapability(v4l2.CapabilityVideoOutput))), "Hardware", "V4l2", nodeName, "Capabilities", "VideoOutput")
		vals <- trie.Insert(trie.Some(strconv.FormatBool(device.MustHasCapability(v4l2.CapabilityVideoCapture))), "Hardware", "V4l2", nodeName, "Capabilities", "VideoCapture")
		vals <- trie.Insert(trie.Some(strconv.FormatBool(device.MustHasCapability(v4l2.CapabilityStreaming))), "Hardware", "V4l2", nodeName, "Capabilities", "StreamingIO")

		ffs, err := device.FormatFamilies()
		if err != nil {
			vals <- trie.Insert(trie.Error(err), "Hardware", "V4l2", nodeName, "Formats")
			continue
		}
		defer ffs.Close()
		f := 0
		var ff v4l2.FormatFamily
		for ffs.Next() {
			err := ffs.Decode(&ff)
			if err != nil {
				vals <- trie.Insert(trie.Error(err), "Hardware", "V4l2", nodeName, "Formats", strconv.Itoa(f))
				continue
			}
			vals <- trie.Insert(trie.Some(ff.PixelFormat().String()), "Hardware", "V4l2", nodeName, "Formats", strconv.Itoa(f), "Name")
			vals <- trie.Insert(trie.Some(ff.Description()), "Hardware", "V4l2", nodeName, "Formats", strconv.Itoa(f), "Description")
			vals <- trie.Insert(trie.Some(strconv.FormatBool(ff.HasFlags(v4l2.FormatFamilyFlagCompressed))), "Hardware", "V4l2", nodeName, "Formats", strconv.Itoa(f), "Compressed")
			vals <- trie.Insert(trie.Some(strconv.FormatBool(ff.HasFlags(v4l2.FormatFamilyFlagEmulated))), "Hardware", "V4l2", nodeName, "Formats", strconv.Itoa(f), "Emulated")
			f = f + 1
		}
	}
}
