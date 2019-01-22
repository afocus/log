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

	g := o.Ctx("aaaaa").Tag("进车")
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

func TestConsoleJSON(t *testing.T) {
	c := New()
	c.UseJSON(true)
	o := log.New(log.DEBUG, c)
	o.Debug("Debug")
	o.Info("Info")
	o.Warn("Warn")
	o.Error("Error")
	o.Fatal("fatal")

	g := o.Ctx("aaaaa").Tag("进车")

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
