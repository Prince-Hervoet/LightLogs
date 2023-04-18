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

	go func() {
		for i := 0; i < 10000; i++ {
			test.Info(strconv.FormatInt(int64(i), 10), "这是一条日志")
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			test.Info(strconv.FormatInt(int64(i), 10), "这是一条日志")
		}
	}()

	for {
	}
	time.Sleep(5000)
	test.Close()
}
