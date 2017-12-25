package net

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/afocus/log"
)

func createHttpServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, _ := httputil.DumpRequest(r, true)
		fmt.Println("msg:", string(b))
		w.Write([]byte("ok"))
	})
	http.ListenAndServe(":6546", nil)
}

// func TestHttp(t *testing.T) {
// 	go createHttpServer()

// 	opt := &Option{
// 		Addr:   "http://127.0.0.1:6546",
// 		Method: "POST",
// 	}
// 	loger := log.New(log.DEBUG, New(opt))
// 	loger.Debug("test...")
// 	loger.Info("xxxxx")
// 	time.Sleep(time.Second * 3)
// }

type customTestHttp struct {
	*Http
}

func (c *customTestHttp) Format(ev *log.Event) []byte {
	b, _ := json.Marshal(ev)
	return b
}

func TestHttpCustom(t *testing.T) {
	go createHttpServer()
	opt := &Option{
		Addr:   "http://127.0.0.1:6546",
		Method: "POST",
		Header: http.Header{"Content-type": {"appliction/json"}},
	}
	loger := log.New(log.DEBUG, &customTestHttp{New(opt)})
	loger.Debug("test...")
	loger.Info("xxxxx")
	time.Sleep(time.Second * 30)
}
