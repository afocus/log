package console

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/afocus/log"
)

type Console struct {
	io.Writer
}

func New() *Console {
	return &Console{
		os.Stdout,
	}
}

func (c Console) Format(ev *log.Event) []byte {
	// windows 暂时先不支持彩色输出
	if runtime.GOOS == "windows" && len(os.Getenv("MSYSTEM")) == 0 && len(os.Getenv("cygwin")) == 0 {
		return log.FormatPattern(ev)
	}
	var fcolor, bcolor = 36, 0
	switch ev.Level {
	case log.ERROR, log.FATAL:
		bcolor = 41
	case log.WARN:
		bcolor = 43
		fcolor = 30
	case log.INFO:
		bcolor = 44
	case log.DEBUG:
		bcolor = 35
	}
	data := fmt.Sprintf(
		"%s \x1b[%d;%dm %s \x1b[0m %s %s-%s \x1b[0;32m→\x1b[0m %s\n",
		ev.Timestamp, fcolor, bcolor, ev.Level, ev.File, ev.ID, ev.Action, ev.Message,
	)
	return []byte(data)
}
