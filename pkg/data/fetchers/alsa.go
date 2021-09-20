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
		}
	}
}
