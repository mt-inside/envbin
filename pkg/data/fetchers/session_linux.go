package fetchers

import (
	"context"
	"os"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getSessionData)
}

func getSessionData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	vals <- Insert(Some(os.Getenv("XDG_SESSION_ID")), "Processes", "0", "Session", "ID")
	vals <- Insert(Some(os.Getenv("XDG_SESSION_CLASS")), "Processes", "0", "Session", "Class")
	vals <- Insert(Some(os.Getenv("XDG_SESSION_TYPE")), "Processes", "0", "Session", "Type")
	vals <- Insert(Some(os.Getenv("XDG_SEAT")), "Processes", "0", "Session", "Seat")
	vals <- Insert(Some(os.Getenv("XDG_VTNR")), "Processes", "0", "Session", "VT")
}