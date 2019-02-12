package file

import (
	"testing"

	"github.com/afocus/log"
)

func TestFile(t *testing.T) {
	f, _ := New(&Option{Path: "test.log"})
	o := log.New(log.DEBUG, f)
	o.Debug("Debug")
	o.Info("Info")
	o.Warn("Warn")
	o.Error("Error")
	o.Fatal("fatal")

	g := o.Ctx(log.CreateID()).Tag("进车")
	g.Warnf("此数据已处理过")
	g.Free()

}

func TestFileJSON(t *testing.T) {
	f, _ := New(&Option{Path: "test.log", UseJSON: true})
	o := log.New(log.DEBUG, f)

	g := o.Ctx(log.CreateID()).Tag("进车")

	datamap := map[string]interface{}{
		"park_code":  7100000001,
		"vpl_number": "陕A74110",
		"arm_code":   "0101",
		"status":     1,
		"in_time":    "2018-01-02 12:30:42",
	}
	g.Fields(datamap).Info("重要进车数据哈哈哈哈")
	g.Warnf("此数据已处理过 %v", datamap["vpl_number"])
	g.Free()

}
