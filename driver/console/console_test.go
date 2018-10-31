package console

import (
	"testing"

	"github.com/afocus/log"
)

func TestConsoleColor(t *testing.T) {
	o := log.New(log.DEBUG, New())
	o.Debug("Debug")
	o.Info("Info")
	o.Warn("Warn")
	o.Error("Error")
	o.Fatal("fatal")
}

func TestConsoleJSON(t *testing.T) {
	c := New()
	c.UseJSON(false)
	o := log.New(log.DEBUG, c)
	o.Debug("Debug")
	o.Info("Info")
	o.Warn("Warn")
	o.Error("Error")
	o.Fatal("fatal")
}
