package lightdouble

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Logger interface {
	Info(msg ...string)
	Warn(msg ...string)
	Error(msg ...string)
	Debug(msg ...string)
	Start()
	Stop()
}

var singleLogger Logger = nil

const ID_FORMAT = 'i'
const LEVEL_FORMAT = 'l'
const DATE_FORMAT = 'd'
const STRING_FORMAT = 's'

const INFO = "INFO"
const WARN = "WARN"
const ERROR = "ERROR"
const DEBUG = "DEBUG"

type DoubleLogger struct {
	bufferLimit    int
	waitingTime    int64
	maxBufferCount int
	format         string
	writeBuffer    *builderBuffer
	flushBuffer    *builderBuffer
	flushingList   chan *builderBuffer
	logsBuffer     chan string
	mu             *sync.Mutex
	filePointer    *os.File
}

// limit the max size of a buffer
// waitingTime
// path the file path
func NewDoubleLogger(limit int, maxBufferCount int, waitingTime int64, path string) (*DoubleLogger, error) {
	if limit <= 0 || maxBufferCount <= 0 || waitingTime <= 0 {
		return nil, errors.New("init error")
	}
	fileName := LookupFileName(path)
	file := openFile(path + "/" + fileName)
	return &DoubleLogger{
		bufferLimit:    limit,
		waitingTime:    waitingTime,
		maxBufferCount: maxBufferCount,
		format:         "%i: [%d] [%l] %s",
		writeBuffer:    newBuffer(limit),
		flushBuffer:    newBuffer(limit),
		flushingList:   make(chan *builderBuffer, maxBufferCount),
		logsBuffer:     nil,
		filePointer:    file,
		mu:             &sync.Mutex{},
	}, nil
}

// format like "%i: %d [%l] %s"
func (logger *DoubleLogger) SetFormat(format string) {
	logger.format = format
}

func (logger *DoubleLogger) Info(msg ...string) {
	logger.goToBuffer(INFO, msg)
}

func (logger *DoubleLogger) Warn(msg ...string) {
	logger.goToBuffer(WARN, msg)
}

func (logger *DoubleLogger) Error(msg ...string) {
	logger.goToBuffer(ERROR, msg)
}

func (logger *DoubleLogger) Debug(msg ...string) {
	logger.goToBuffer(DEBUG, msg)
}

func (logger *DoubleLogger) Start() {
	if logger.logsBuffer != nil {
		return
	}
	logger.logsBuffer = make(chan string, 1024)
	if logger.filePointer != nil {
		go logger.write()
		go flushTaskFunc(logger)
	}
}

func (logger *DoubleLogger) goToBuffer(level string, msg []string) {
	if logger.logsBuffer == nil {
		return
	}
	content := jointMessage(level, logger.format, msg)
	logger.logsBuffer <- content
	fmt.Println(content)
}

func (logger *DoubleLogger) write() {
	fmt.Println("start write to buffer")
	for mp := range logger.logsBuffer {
		func() {
			if logger.writeBuffer.write(mp) == 1 {
				return
			}
			logger.mu.Lock()
			defer logger.mu.Unlock()
			logger.flushingList <- logger.writeBuffer
			if logger.flushBuffer == nil {
				logger.writeBuffer = newBuffer(logger.bufferLimit)
			} else {
				logger.writeBuffer = logger.flushBuffer
				logger.flushBuffer = nil
			}
			logger.writeBuffer.write(mp)
		}()
	}

}

func flushTaskFunc(dl *DoubleLogger) {
	fmt.Println("start wait to flush")
	for buffer := range dl.flushingList {
		writeStringToFile(dl.filePointer, buffer.getString())
		buffer.reset()
		func() {
			dl.mu.Lock()
			defer dl.mu.Unlock()
			if dl.flushBuffer == nil {
				dl.flushBuffer = buffer
			} else {
				buffer = nil
			}
		}()
	}
}

func jointMessage(level string, format string, msg []string) string {
	var builder strings.Builder
	meaning := false
	msgIndex := 0
	for i, c := range format {
		if meaning {
			meaning = false
			continue
		}
		if c == '%' {
			if i+1 < len(format) {
				switch format[i+1] {
				case ID_FORMAT:
					u, _ := uuid.NewUUID()
					builder.WriteString(u.String())
					fmt.Println(u.String())
					meaning = true
					break
				case DATE_FORMAT:
					now := time.Now()
					builder.WriteString(now.Format("2006-01-02 15:04:05"))
					meaning = true
					break
				case LEVEL_FORMAT:
					builder.WriteString(level)
					meaning = true
					break
				case STRING_FORMAT:
					if msgIndex < len(msg) {
						fmt.Println(msg[0])
						builder.WriteString(msg[msgIndex])
						msgIndex++
					}
					meaning = true
					break
				default:
					meaning = false
					break
				}
			}
		} else {
			builder.WriteRune(c)
		}
	}
	builder.WriteByte('\n')
	return builder.String()
}
