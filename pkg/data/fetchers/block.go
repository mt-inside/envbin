package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/jaypipes/ghw"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getBlockData)
}

func getBlockData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	prefix := []string{"Hardware", "Block"}

	blk, err := ghw.Block()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), prefix...)
		return
	}

	for _, d := range blk.Disks {
		vals <- trie.Insert(trie.Some(strconv.FormatUint(d.SizeBytes, 10)), "Hardware", "Block", d.Name, "SizeBytes")
		vals <- trie.Insert(trie.Some(d.Vendor), "Hardware", "Block", d.Name, "Vendor")
		vals <- trie.Insert(trie.Some(d.Model), "Hardware", "Block", d.Name, "Model")
		vals <- trie.Insert(trie.Some(strconv.FormatUint(d.PhysicalBlockSizeBytes, 10)), "Hardware", "Block", d.Name, "BlockSizeBytes")
		vals <- trie.Insert(trie.Some(strconv.FormatBool(d.IsRemovable)), "Hardware", "Block", d.Name, "Removable")
		vals <- trie.Insert(trie.Some(d.StorageController.String()), "Hardware", "Block", d.Name, "ControllerType")
		vals <- trie.Insert(trie.Some(d.SerialNumber), "Hardware", "Block", d.Name, "Serial")

		for _, p := range d.Partitions {
			vals <- trie.Insert(trie.Some(strconv.FormatUint(p.SizeBytes, 10)), "Hardware", "Block", d.Name, "Partitions", p.Name, "SizeBytes")
			vals <- trie.Insert(trie.Some(p.Type), "Hardware", "Block", d.Name, "Partitions", p.Name, "Filesystem")
			vals <- trie.Insert(trie.Some(p.MountPoint), "Hardware", "Block", d.Name, "Partitions", p.Name, "MountPoint")
			vals <- trie.Insert(trie.Some(strconv.FormatBool(p.IsReadOnly)), "Hardware", "Block", d.Name, "Partitions", p.Name, "ReadOnly")
			vals <- trie.Insert(trie.Some(p.UUID), "Hardware", "Block", d.Name, "Partitions", p.Name, "UUID")
		}
	}
}
