//go:build native
// +build native

package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/yobert/alsa"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getAlsaData)
}

func getAlsaData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	cards, err := alsa.OpenCards()
	if err != nil {
		vals <- Insert(Error(err), "Hardware", "Sound", "Alsa")
		return
	}
	defer alsa.CloseCards(cards)

	for _, card := range cards {
		vals <- Insert(Some(card.Title), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Name")
		vals <- Insert(Some(card.Path), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Path")

		devices, err := card.Devices()
		if err != nil {
			vals <- Insert(Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices")
			continue
		}
		for _, device := range devices {
			vals <- Insert(Some(device.Title), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Name")
			vals <- Insert(Some(device.Path), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Path")
			vals <- Insert(Some(device.Type.String()), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Type")
			vals <- Insert(Some(strconv.FormatBool(device.Play)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Play")
			vals <- Insert(Some(strconv.FormatBool(device.Record)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Record")

			if err := device.Open(); err != nil {
				vals <- Insert(Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample")
				continue
			}
			if err := device.Prepare(); err != nil {
				vals <- Insert(Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample")
				continue
			}
			vals <- Insert(Some(strconv.Itoa(device.BufferFormat().Channels)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Channels")
			vals <- Insert(Some(device.BufferFormat().SampleFormat.String()), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample", "Format")
			vals <- Insert(Some(strconv.Itoa(device.BufferFormat().Rate)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample", "Rate")
		}
	}
}
