package mylog

import (
	lightdouble "lightDouble"
	"strconv"
	"testing"
)

func TestGggg(t *testing.T) {
	BenchmarkMyLog(&testing.B{})
}

func BenchmarkMyLog(b *testing.B) {
	logger, _ := lightdouble.NewDoubleLogger(4096, 50, 200, "D:\\myprojects\\LightDouble\\test")
	logger.Start()
	logger.SetFormat("%s [%d] [%l] %s")
	for n := 0; n < 100000; n++ {
		logger.Info(strconv.FormatInt(int64(n), 10), "这是一条日志")
	}
	logger.Close()
}
