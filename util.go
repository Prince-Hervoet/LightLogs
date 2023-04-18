package lightdouble

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

const logFileLimit = 209715200

func openFile(path string) *os.File {
	// 打开文件，如果不存在则创建，可读写
	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return file
}

func closeFile(file *os.File) error {
	if file == nil {
		return nil
	}
	err := file.Close()
	return err
}

func Int32ToStringAndZero(number int64, count int) string {
	str := strconv.FormatInt(number, 10)
	if len(str) >= count {
		return str
	}
	remain := count - len(str)
	var builder strings.Builder
	for i := 0; i < remain; i++ {
		builder.WriteByte('0')
	}
	builder.WriteString(str)
	return builder.String()
}

func writeStringToFile(file *os.File, data string) error {
	if file == nil || data == "" {
		return nil
	}
	_, err := file.WriteString(data)
	return err
}

func stringToBytes(str string) []byte {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&str))
	byteSlice := *(*[]byte)(unsafe.Pointer(sliceHeader))
	return byteSlice
}

func LookupFileName(path string) string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("无法读取目录:", err)
		os.Exit(1)
	}
	var fileInfos []fileInfo
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "log-") {
			fileInfo := fileInfo{
				name: file.Name(),
				size: file.Size(),
			}
			fileInfos = append(fileInfos, fileInfo)
		}
	}

	if len(fileInfos) == 0 {
		return "log-0000000000"
	}
	// 按照文件名字进行降序排序
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].name > fileInfos[j].name
	})
	prevName := fileInfos[0].name
	numStr := strings.Split(prevName, "-")[1]
	num, err2 := strconv.ParseInt(numStr, 10, 64)
	if err2 == nil {
		return "log-" + Int32ToStringAndZero(num+1, 10)
	}
	return ""
}

type fileInfo struct {
	name string
	size int64
}
