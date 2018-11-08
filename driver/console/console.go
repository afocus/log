package console

import (
	"encoding/json"
	"io"
	"os"

	"github.com/afocus/log"
	"github.com/gookit/color"
)

type Console struct {
	io.Writer
	usejson       bool
	usejsonIndent bool
}

func New() *Console {
	return &Console{
		Writer: os.Stdout,
	}
}

func (c *Console) UseJSON(indent bool) {
	c.usejson = true
	c.usejsonIndent = indent
}

func (c *Console) toJSON(ev *log.Event) []byte {
	if c.usejsonIndent {
		b, _ := json.MarshalIndent(ev, "", "	")
		b = append(b, '\n')
		return b
	}
	b, _ := json.Marshal(ev)
	return b
}

func (c *Console) Format(ev *log.Event) []byte {
	if c.usejson {
		return c.toJSON(ev)
	}
	var theme *color.Theme
	data := string(log.FormatPattern(ev))
	switch ev.Level {
	case log.FATAL:
		theme = color.Error
	case log.ERROR:
		theme = color.Danger
	case log.WARN:
		theme = color.Warn
	case log.INFO:
		theme = color.Primary
	default:
		return []byte(data)
	}
	return []byte(theme.Sprint(data))
}
