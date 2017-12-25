package net

import (
	"bytes"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/afocus/log"
)

type Http struct {
	// 配置信息
	option *Option

	client *http.Client
	reqs   chan *http.Request
	// 自动扩容缩放任务池大小
	workSize int32
}

type Option struct {
	// 请求地址
	Addr string
	// 请求的方法 GET POST PUT...
	Method string
	// 携带的header头信息
	Header http.Header
	// 请求超时时间 默认无限制
	Timeout time.Duration
	// 最大并发处理数量
	MaxWorks int32
}

func New(option *Option) *Http {
	if option.MaxWorks < 4 {
		option.MaxWorks = 4
	}
	ins := &Http{
		option: option,
		reqs:   make(chan *http.Request, 20),
		client: &http.Client{
			Timeout: option.Timeout,
		},
	}
	ins.addWork()
	return ins
}

// Format 格式化内容
// 如果要换别的格式请继承Http并重新实现Format方法
func (h *Http) Format(ev *log.Event) []byte {
	return log.FormatPattern(ev)
}

func (h *Http) Write(data []byte) (int, error) {
	if data == nil {
		return 0, nil
	}
	req, err := http.NewRequest(h.option.Method, h.option.Addr, bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	if h.option.Header != nil {
		req.Header = h.option.Header
	}
	select {
	case h.reqs <- req:
	default:
		if !h.addWork() {
			// 处理不过来且超过最大任务池限制
			// 为了保证不堵塞任务 先放弃吧
			return 0, errors.New("log http task busy")
		}
	}
	return len(data), nil
}

func (h *Http) worker() {
	atomic.AddInt32(&h.workSize, 1)
	for {
		select {
		case req := <-h.reqs:
			resp, err := h.client.Do(req)
			if err != nil {
				continue
			}
			resp.Body.Close()
		case <-time.After(time.Minute):
			// 过了一分钟了 还是木有接受任何任务
			// 则说明此work时空闲的,需要结束掉
			// 任务数-1
			if atomic.AddInt32(&h.workSize, -1) == 0 {
				// 保证至少有一个任务再进行中
				atomic.AddInt32(&h.workSize, 1)
			} else {
				return
			}
		}
	}
}

func (h *Http) addWork() bool {
	if atomic.LoadInt32(&h.workSize) == h.option.MaxWorks {
		return false
	}
	go h.worker()
	return true
}
