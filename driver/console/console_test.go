package console

import (
	"testing"

	"github.com/afocus/log"
)

func TestConsoleColor(t *testing.T) {
	o := log.New(log.DEBUG, Console)
	o.Debug("mo")
	o.Info("mo")
	o.Warn("mo")
	o.Error("mo")
}
