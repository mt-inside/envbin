package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/jaypipes/ghw"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getBlockData)
}

func getBlockData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	prefix := []string{"Hardware", "Block"}

	blk, err := ghw.Block()
	if err != nil {
		vals <- Insert(Error(err), prefix...)
		log.Error(err, "Can't get block device info")
		return
	}

	for _, d := range blk.Disks {
		vals <- Insert(Some(strconv.FormatUint(d.SizeBytes, 10)), "Hardware", "Block", d.Name, "SizeBytes")
		vals <- Insert(Some(d.Vendor), "Hardware", "Block", d.Name, "Vendor")
		vals <- Insert(Some(d.Model), "Hardware", "Block", d.Name, "Model")
		vals <- Insert(Some(strconv.FormatUint(d.PhysicalBlockSizeBytes, 10)), "Hardware", "Block", d.Name, "BlockSizeBytes")
		vals <- Insert(Some(strconv.FormatBool(d.IsRemovable)), "Hardware", "Block", d.Name, "Removable")
		vals <- Insert(Some(d.StorageController.String()), "Hardware", "Block", d.Name, "ControllerType")
		vals <- Insert(Some(d.SerialNumber), "Hardware", "Block", d.Name, "Serial")

		for _, p := range d.Partitions {
			vals <- Insert(Some(strconv.FormatUint(p.SizeBytes, 10)), "Hardware", "Block", d.Name, "Partitions", p.Name, "SizeBytes")
			vals <- Insert(Some(p.Type), "Hardware", "Block", d.Name, "Partitions", p.Name, "Filesystem")
			vals <- Insert(Some(p.MountPoint), "Hardware", "Block", d.Name, "Partitions", p.Name, "MountPoint")
			vals <- Insert(Some(strconv.FormatBool(p.IsReadOnly)), "Hardware", "Block", d.Name, "Partitions", p.Name, "ReadOnly")
			vals <- Insert(Some(p.UUID), "Hardware", "Block", d.Name, "Partitions", p.Name, "UUID")
		}
	}
}
