package data

import (
	"context"
	"sync"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data/trie"
)

type pluginfn func(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg)

var (
	plugins []pluginfn
)

func RegisterPlugin(p pluginfn) {
	plugins = append(plugins, p)
}

func GetData(ctx context.Context, log logr.Logger) *trie.Trie {
	wg := sync.WaitGroup{}
	vals := make(chan trie.InsertMsg)

	for _, p := range plugins {
		wg.Add(1)
		f := p // Avoid capture of the loop variable
		go func() {
			f(ctx, log, vals)
			log.V(1).Info("Plugin done")
			// If one of the plugins is sticking and you wanna see which one it is
			// err := pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
			// if err != nil {
			// 	panic(err)
			// }
			wg.Done()
		}()
	}

	go func() {
		// TODO: This relies on them all returning, and doing so asap after the deadline pops. To guard against rogue plugins we should also wait for the deadline ourselves, add a grace period, then close the channel. Will need some way to stop an errant plugin (which finally tries to write to the channel) from taking the whole system down when it panics trying to write to the closed channel
		//<-ctx.Done()
		wg.Wait()
		close(vals) // Causes BuildFromInsertMsgs to finish and return
		// TODO: plugins need to deal with trying to write to a closed channel (what actually happens when you panic on a background goroutine, and can we recover them individually to catch it?)
		// TODO: Would be nice to mark plugins as timed-out. Atm they do it (if they notice), but a) they have to notice and b) the channel will be closed at that point. When plugins are registered, that should be with the prefix - we record that so we can write TimedOut entires on their behalf. Also means they don't have to use their full prefix everywhere (pass them a PrefixChan()?). Yes several write to several parts of the tree; they should register multiple times for each prefix they "own" (one compilation unit per /source/, but not per tree prefix). Enforce that they don't clash - nothing should add a value or child as a sibling of a _value_ (or a value as a sibling of a child?)
	}()

	t := trie.BuildFromInsertMsgs(log.WithName("trie"), vals)

	// TODO: could even stream results to the client (but it'd have to build the trie itself or something)
	return t
}
