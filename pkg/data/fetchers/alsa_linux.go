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

func getAlsaData(ctx context.Context, log logr.Logger, t *Trie) {
	cards, err := alsa.OpenCards()
	if err != nil {
		t.Insert(Error(err), "Hardware", "Sound", "Alsa")
		return
	}
	defer alsa.CloseCards(cards)

	for _, card := range cards {
		t.Insert(Some(card.Title), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Name")
		t.Insert(Some(card.Path), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Path")

		devices, err := card.Devices()
		if err != nil {
			t.Insert(Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices")
			continue
		}
		for _, device := range devices {
			t.Insert(Some(device.Title), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Name")
			t.Insert(Some(device.Path), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Path")
			t.Insert(Some(device.Type.String()), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Type")
			t.Insert(Some(strconv.FormatBool(device.Play)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Play")
			t.Insert(Some(strconv.FormatBool(device.Record)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Record")

			if err := device.Open(); err != nil {
				t.Insert(Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample")
				continue
			}
			if err := device.Prepare(); err != nil {
				t.Insert(Error(err), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample")
				continue
			}
			t.Insert(Some(strconv.Itoa(device.BufferFormat().Channels)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Channels")
			t.Insert(Some(device.BufferFormat().SampleFormat.String()), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample", "Format")
			t.Insert(Some(strconv.Itoa(device.BufferFormat().Rate)), "Hardware", "Sound", "Alsa", "Cards", strconv.Itoa(card.Number), "Devices", strconv.Itoa(device.Number), "Sample", "Rate")
		}
	}
}
