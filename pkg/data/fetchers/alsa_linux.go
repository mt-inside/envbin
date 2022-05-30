//go:build native
// +build native

package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/yobert/alsa"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getAlsaData)
}

func getAlsaData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	cards, err := alsa.OpenCards()
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "Sound", "Alsa")
		return
	}
	defer alsa.CloseCards(cards)

	for _, card := range cards {
		vals <- trie.Insert(trie.Some(card.Title), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Name")
		vals <- trie.Insert(trie.Some(card.Path), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Path")

		devices, err := card.Devices()
		if err != nil {
			vals <- trie.Insert(trie.Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices")
			continue
		}
		for _, device := range devices {
			vals <- trie.Insert(trie.Some(device.Title), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Name")
			vals <- trie.Insert(trie.Some(device.Path), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Path")
			vals <- trie.Insert(trie.Some(device.Type.String()), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Type")
			vals <- trie.Insert(trie.Some(strconv.FormatBool(device.Play)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Play")
			vals <- trie.Insert(trie.Some(strconv.FormatBool(device.Record)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Record")

			if err := device.Open(); err != nil {
				vals <- trie.Insert(trie.Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample")
				continue
			}
			if err := device.Prepare(); err != nil {
				vals <- trie.Insert(trie.Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample")
				continue
			}
			vals <- trie.Insert(trie.Some(strconv.Itoa(device.BufferFormat().Channels)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Channels")
			vals <- trie.Insert(trie.Some(device.BufferFormat().SampleFormat.String()), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample", "Format")
			vals <- trie.Insert(trie.Some(strconv.Itoa(device.BufferFormat().Rate)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample", "Rate")
		}
	}
}
