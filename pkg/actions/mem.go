package actions

import (
	"os"
	"runtime"
)

var allocs = make([][]byte, 0, 0)

func AllocAndTouch(bytes int64) {
	bs := make([]byte, bytes, bytes)
	var i int64
	for i = 0; i < bytes; i += int64(os.Getpagesize()) {
		bs[i] = 69
	}
	allocs = append(allocs, bs)
}

/* This doesn't have an effect on virtual memory usage. I /think/ its working, but I think it'll just reduce physical usage, if that? free() just sticks the virtual address space back on the free list. The OS will attempt to reclaim anon pages, but a) can they actually go on the free list not just get swapped out, and b) will it only do this under memory pressure? */
func FreeAllocs() {
	for i := range allocs {
		allocs[i] = nil
	}
	runtime.GC()
}
