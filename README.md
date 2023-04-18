# LightLogs
This is a logger,and it has two buffers to increase throughput.
```
  package main

  import (
    "fmt"
    lightdouble "lightDouble/lightDouble"
    "strconv"
    "sync"
    "time"
  )

  func main() {
    test, _ := lightdouble.NewDoubleLogger(32000, 100, 200, "./test")
    test.Start()
    test.SetFormat("%s [%d] [%l] %s")
    wg := sync.WaitGroup{}

    singleSize := len("1 [2023-04-18 23:46:32] [INFO] 这是一条日志")
    count := 100000
    sum := singleSize * count * 3
    wg.Add(3)
    start := time.Now().UnixMilli()
    go func() {
      for i := 0; i < count; i++ {
        test.Info(strconv.FormatInt(int64(i), 10), "这是一条日志")
      }
      wg.Done()
    }()

    go func() {
      for i := 0; i < count; i++ {
        test.Info(strconv.FormatInt(int64(i), 10), "这是一条日志")
      }
      wg.Done()
    }()

    go func() {
      for i := 0; i < count; i++ {
        test.Info(strconv.FormatInt(int64(i), 10), "这是一条日志")
      }
      wg.Done()
    }()
    wg.Wait()
    end := time.Now().UnixMilli()
    fmt.Print("耗时: ")
    fmt.Print((end - start))
    fmt.Println("ms")

    fmt.Print("写入速度: ")
    fmt.Print((int64(sum) / (end - start)))
    fmt.Println(" byte/s")
    test.Close()
  }
```
