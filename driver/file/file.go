package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/afocus/log"
)

// LogFile 实现log.FormatWriter
// 主要做日志分割
type LogFile struct {
	fd     *os.File
	option *Option
	rsize  uint32
	name   string
}

// Option 日志文件配置信息
type Option struct {
	// 日志文件的路径 如server.log,/data/logs/xxx.log
	Path string
	// 最多保存的日志文件数量
	// 默认不限制
	MaxFileCount int
	// 单个日志文件的最大大小 单位MB
	// 默认不限制 当单个文件不限制的情况下 MaxFileCount无效 因为永远不会分割
	MaxFileSize uint64
	// 是否使用json格式输出
	UseJSON bool
}

func New(opt *Option) (*LogFile, error) {
	if opt.MaxFileCount < 0 {
		opt.MaxFileCount = 0
	}
	if opt.MaxFileSize != 0 {
		opt.MaxFileSize = opt.MaxFileSize * (1 << 20)
	}
	opt.Path, _ = filepath.Abs(strings.TrimSuffix(opt.Path, "/"))
	f := &LogFile{
		option: opt,
		name:   opt.Path,
	}
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(opt.Path), 0777); err != nil {
		return nil, err
	}
	return f, f.rotate()
}

func (f *LogFile) Write(data []byte) (int, error) {
	n, err := f.fd.Write(data)
	if err != nil {
		return n, err
	}
	f.rsize += uint32(n)
	if f.option.MaxFileSize > 0 && uint64(f.rsize) > f.option.MaxFileSize {
		err = f.rotate()
	}
	return n, err
}

func (f *LogFile) Format(ev *log.Event) []byte {
	if f.option.UseJSON {
		b, _ := json.MarshalIndent(ev, "", "	")
		return b
	}
	return log.FormatPattern(ev)
}

// 获取目录下指定前缀的所有日志文件
func (f *LogFile) removeFiles() {
	fs, err := filepath.Glob(fmt.Sprintf("%s.*", f.option.Path))
	if err != nil {
		return
	}
	sort.Strings(fs)
	x := len(fs) - (f.option.MaxFileCount - 1)
	if f.option.MaxFileCount > 0 && x > 0 {
		dels := fs[:x]
		for _, v := range dels {
			os.Remove(v)
		}
	}
}

// 分割
func (f *LogFile) rotate() error {
	f.removeFiles()
	if f.fd != nil {
		f.fd.Sync()
		f.fd.Close()
		os.Rename(f.name, f.name+time.Now().Format(".20060102150405"))
	}
	// 创建最新的日志文件
	fd, err := os.OpenFile(f.name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	fi, err := fd.Stat()
	if err != nil {
		return err
	}
	f.fd = fd
	f.rsize = uint32(fi.Size())
	return nil
}
