package log

import (
	"crypto/rand"
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

// Level 日志级别类型
type Level int32

const (
	// 日志级别 兼容Log4j
	// 它假设级别是有序的。对于标准级别，其顺序为：DEBUG < INFO < WARN < ERROR < FATAL < OFF

	// DEBUG 指明细致的事件信息，对调试应用最有用
	DEBUG Level = iota
	// INFO 指明描述信息，从粗粒度上描述了应用运行过程
	INFO
	// WARN 指明潜在的有害状况
	WARN
	// ERROR 指明错误事件，但应用可能还能继续运行
	ERROR
	// FATAL 指明非常严重的错误事件，可能会导致应用终止执行
	FATAL
	// OFF 最高级别，用于关闭日志
	OFF
)

// String 返回日志等级的描述
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERRO"
	case FATAL:
		return "FATA"
	default:
		return "UNKN"
	}
}

// MarshalJSON 保证json编码level是不是显示数字而是string
func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

// TimestampLayout 日志时间格式化模板
var TimestampLayout = "2006-01-02 15:04:05"

// FormatPattern 扁平化格式化日志事件
var FormatPattern = func(ev *Event) []byte {
	data := fmt.Sprintf(
		"%s [%s] %s %s-%s → %s",
		ev.Timestamp, ev.Level, ev.File, ev.ID, ev.Action, ev.Message,
	)
	d := []byte(data)
	if length := len(d); d[length-1] == '\n' {
		return d
	}
	return append(d, '\n')
}

// Formater 格式化日志事件到字符串
type Formater interface {
	Format(*Event) []byte
}

// FormatWriter 格式化输入
// 用于定制日志输出内容的样式
type FormatWriter interface {
	io.Writer
	Formater
}

// Event 日志事件 记录了日志的必要信息
// 可以通过Formater接口输入格式化样式后的数据
type Event struct {
	// 日志产生时的时间
	Timestamp string `json:"timestamp"`
	// 日志等级
	Level Level `json:"level"`
	// 所在文件行数file:line
	// main:20
	File string `json:"file,omitempty"`
	// 日志id 只有Ctx方式才会使用
	// 主要用于上下文关联
	ID string `json:"logid,omitempty"`
	// 日志动作名称 描述干什么的 如 login,callback...
	Action string `json:"action,omitempty"`
	// 日志内容
	Message string `json:"message"`
}

// CreateID 简单的返回一个随机字符串id
func CreateID() string {
	x := make([]byte, 16)
	io.ReadFull(rand.Reader, x)
	return fmt.Sprintf("%x", x)
}

// Logger 日志对象
type Logger struct {
	// 用于并发安全的锁
	mu sync.Mutex
	// 实现FormatWriter接口的输出对象
	// 可以时多个输出对象 为了效率非必要尽量不要太多
	outs []FormatWriter
	// 日志等级限制
	// 输入的等级>=限制才能输出
	lvl Level
}

func New(lvl Level, outs ...FormatWriter) *Logger {
	return &Logger{
		outs: outs,
		lvl:  lvl,
	}
}

func (o *Logger) lockCall(f func()) {
	o.mu.Lock()
	f()
	o.mu.Unlock()
}

var eventObjPool = &sync.Pool{New: func() interface{} { return new(Event) }}

// Output 输出日志消息
// 核心方法 所有日志输出全部以及此方法
func (o *Logger) Output(calldept int, level Level, acname, id, msg string) error {
	// 等级不足以输出
	if o.lvl > level {
		return nil
	}

	if level == FATAL {
		// 追加调用堆栈
		for i := calldept; i < calldept+5; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg += fmt.Sprintf("\n%s:%d", file, line)
		}
	}

	// 从对象池中取出一个style对象并赋值
	ev := eventObjPool.Get().(*Event)
	ev.Timestamp = time.Now().Format(TimestampLayout)
	ev.ID = id
	ev.Level = level
	ev.Action = acname
	ev.Message = msg

	// 获取所在文件以及行数
	_, file, line, ok := runtime.Caller(calldept)
	if !ok {
		file = "???"
		line = 0
	}
	length := len(file) - 1
	for i := length; i > 0; i-- {
		if file[i] == '/' {
			file = file[i+1 : length-2]
			break
		}
	}
	ev.File = fmt.Sprintf("%s:%d", file, line)

	// 保证并发安全
	var err error
	for _, out := range o.outs {
		// Write 方法不要堵塞会
		// 如何高性能杜绝或防止这里堵塞呢？todo
		o.lockCall(func() {
			out.Write(out.Format(ev))
		})
		if err != nil {
			continue
		}
	}
	// 放入对象池中
	eventObjPool.Put(ev)
	return err
}

// Debug
func (o *Logger) Debug(s ...interface{}) {
	o.Output(2, DEBUG, "", "", fmt.Sprint(s...))
}

func (o *Logger) Info(s ...interface{}) {
	o.Output(2, INFO, "", "", fmt.Sprint(s...))
}

func (o *Logger) Warn(s ...interface{}) {
	o.Output(2, WARN, "", "", fmt.Sprint(s...))
}

func (o *Logger) Error(s ...interface{}) {
	o.Output(2, ERROR, "", "", fmt.Sprint(s...))
}

func (o *Logger) Fatal(s ...interface{}) {
	o.Output(2, FATAL, "", "", fmt.Sprint(s...))
}

// format

func (o *Logger) Debugf(s string, args ...interface{}) {
	o.Output(2, DEBUG, "", "", fmt.Sprintf(s, args...))
}

func (o *Logger) Infof(s string, args ...interface{}) {
	o.Output(2, INFO, "", "", fmt.Sprintf(s, args...))
}

func (o *Logger) Warnf(s string, args ...interface{}) {
	o.Output(2, WARN, "", "", fmt.Sprintf(s, args...))
}

func (o *Logger) Errorf(s string, args ...interface{}) {
	o.Output(2, ERROR, "", "", fmt.Sprintf(s, args...))
}

func (o *Logger) Fatalf(s string, args ...interface{}) {
	o.Output(2, FATAL, "", "", fmt.Sprintf(s, args...))
}

// Ctx 携带日志id和事件名的日志对象
// 主要用于通过id串联一些日志 起到查询方便
type Ctx struct {
	o       *Logger
	id, tag string
}

var ctxPool = sync.Pool{New: func() interface{} { return new(Ctx) }}

// Ctx 创建一个包含指定id的ctx对象
func (o *Logger) Ctx(id string) *Ctx {
	ctx := ctxPool.Get().(*Ctx)
	ctx.id = id
	ctx.tag = ""
	ctx.o = o
	return ctx
}

// Tag 设置标签名
func (ctx *Ctx) Tag(tag string) *Ctx {
	ctx.tag = tag
	return ctx
}

// Free 释放
func (ctx *Ctx) Free() {
	ctxPool.Put(ctx)
}

func (ctx *Ctx) Debug(s ...interface{}) *Ctx {
	ctx.o.Output(2, DEBUG, ctx.tag, ctx.id, fmt.Sprint(s...))
	return ctx
}

func (ctx *Ctx) Info(s ...interface{}) *Ctx {
	ctx.o.Output(2, INFO, ctx.tag, ctx.id, fmt.Sprint(s...))
	return ctx
}

func (ctx *Ctx) Warn(s ...interface{}) *Ctx {
	ctx.o.Output(2, WARN, ctx.tag, ctx.id, fmt.Sprint(s...))
	return ctx
}

func (ctx *Ctx) Error(s ...interface{}) *Ctx {
	ctx.o.Output(2, ERROR, ctx.tag, ctx.id, fmt.Sprint(s...))
	return ctx
}

func (ctx *Ctx) Fatal(s ...interface{}) *Ctx {
	ctx.o.Output(2, FATAL, ctx.tag, ctx.id, fmt.Sprint(s...))
	return ctx
}

//
func (ctx *Ctx) Debugf(s string, args ...interface{}) *Ctx {
	ctx.o.Output(2, DEBUG, ctx.tag, ctx.id, fmt.Sprintf(s, args...))
	return ctx
}

func (ctx *Ctx) Infof(s string, args ...interface{}) *Ctx {
	ctx.o.Output(2, INFO, ctx.tag, ctx.id, fmt.Sprintf(s, args...))
	return ctx
}

func (ctx *Ctx) Warnf(s string, args ...interface{}) *Ctx {
	ctx.o.Output(2, WARN, ctx.tag, ctx.id, fmt.Sprintf(s, args...))
	return ctx
}

func (ctx *Ctx) Errorf(s string, args ...interface{}) *Ctx {
	ctx.o.Output(2, ERROR, ctx.tag, ctx.id, fmt.Sprintf(s, args...))
	return ctx
}

func (ctx *Ctx) Fatalf(s string, args ...interface{}) *Ctx {
	ctx.o.Output(2, FATAL, ctx.tag, ctx.id, fmt.Sprintf(s, args...))
	return ctx
}
