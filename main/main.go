package main

import (
	lightdouble "lightDouble/lightDouble"
	"strconv"
	"time"
)

func main() {
	test, _ := lightdouble.NewDoubleLogger(1000, 100, 200, "D:\\myprojects\\LightDouble\\test")
	test.Start()
	test.SetFormat("%s [%d] [%l] %s")
	for i := 0; i < 10000; i++ {
		test.Info(strconv.FormatInt(int64(i), 10), "这是一条日志")
	}
	time.Sleep(3000)
	test.Close()
}
