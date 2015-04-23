package log

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/zhyueh/figo/toolkit"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	TruncateImmediately = 0
	TruncateDay         = 1
	TruncateHour        = 2
	TruncateMinute      = 3
	TruncateTenMinute   = 4
)

const (
	LevelDebug = 0
	LevelInfo  = 1
	LevelWarn  = 2
	LevelError = 3
	LevelFatal = 4
)

func levelModeToString(levelMode int8) string {
	switch levelMode {
	case LevelInfo:
		return "I"
	case LevelDebug:
		return "D"
	case LevelWarn:
		return "W"
	case LevelError:
		return "E"
	case LevelFatal:
		return "F"
	}
	return "level undefined"
}

type DataLogger struct {
	dirPath        string
	name           string
	writingName    string
	truncateMode   int8
	levelMode      int8
	levelString    string
	cacheLineCount int
	lines          chan string
	quit           chan bool
	lock           *sync.Mutex
}

func NewDataLogger(dirPath string, name string, truncateMode int8, levelMode int8, cacheLineCount int) (*DataLogger, error) {
	this := DataLogger{}
	this.dirPath = fmt.Sprintf("%s/%s", dirPath, name)
	this.name = name
	this.truncateMode = truncateMode
	this.levelMode = levelMode
	this.levelString = levelModeToString(levelMode)
	this.cacheLineCount = cacheLineCount

	toolkit.EnsureDirExists(this.dirPath)
	this.writingName = this.getFileName()
	this.lines = make(chan string, 1000000)
	this.quit = make(chan bool, 1)
	this.lock = new(sync.Mutex)
	if this.truncateMode != TruncateImmediately {
		this.autoFlush()
	}
	return &this, nil
}

func (this DataLogger) getFileName() (fileName string) {
	t := time.Now()
	switch this.truncateMode {
	case TruncateImmediately:
		fileName = fmt.Sprintf("%s_%d%02d%02d", this.name, t.Year(), t.Month(), t.Day())
	case TruncateDay:
		fileName = fmt.Sprintf("%s_%d%02d%02d0000", this.name, t.Year(), t.Month(), t.Day())
	case TruncateHour:
		fileName = fmt.Sprintf("%s_%d%02d%02d%02d00", this.name, t.Year(), t.Month(), t.Day(), t.Hour())
	case TruncateMinute:
		fileName = fmt.Sprintf("%s_%d%02d%02d%02d%02d", this.name, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
	case TruncateTenMinute:
		fileName = fmt.Sprintf("%s_%d%02d%02d%02d%02d", this.name, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()/10*10)
	default:
		fileName = this.name
	}
	return
}

func (this *DataLogger) Debug(format string, a ...interface{}) {
	this.log(LevelDebug, format, a...)
}

func (this *DataLogger) Info(format string, a ...interface{}) {
	this.log(LevelInfo, format, a...)
}

func (this *DataLogger) Warn(format string, a ...interface{}) {
	this.log(LevelWarn, format, a...)
}

func (this *DataLogger) Error(format string, a ...interface{}) {
	this.log(LevelError, format, a...)
}

func (this *DataLogger) Fatal(format string, a ...interface{}) {
	this.log(LevelFatal, format, a...)
	os.Exit(-1)
}

func (this *DataLogger) log(levelMode int8, format string, a ...interface{}) {
	if levelMode < this.levelMode {
		return
	}

	this.defaultFormatLog(levelMode, format, a...)
}

func (this *DataLogger) defaultFormatLog(level int8, format string, a ...interface{}) {
	now := time.Now()
	log := fmt.Sprintf(format, a...)
	this.writelnf("[%s %02d%02d%02d %02d:%02d:%02d] %s", this.levelString, now.Year()%100, now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), log)
}

func (this *DataLogger) writelnf(format string, a ...interface{}) {
	this.writeln(fmt.Sprintf(format, a...))
}

func (this *DataLogger) writeln(line string) {
	if this.truncateMode == TruncateImmediately {
		this.flushLine(line)
	} else {
		this.lines <- line
	}
}

func (this *DataLogger) getFileInfo() (*os.File, error) {
	return os.OpenFile(fmt.Sprintf("%s/%s.log", this.dirPath, this.getFileName()), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
}

func (this *DataLogger) flushLine(line string) error {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	file, err := this.getFileInfo()
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(line)
	return nil
}
func (this *DataLogger) Flush() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	lineCount := len(this.lines)

	if lineCount == 0 {
		return nil
	}

	file, err := this.getFileInfo()
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := bytes.Buffer{}
	for i := 1; i <= lineCount; i++ {
		str := <-this.lines
		if !strings.HasSuffix(str, "\n") {
			str += "\n"
		}

		buffer.WriteString(str)

		if i%this.cacheLineCount == 0 {
			file.WriteString(buffer.String())
			buffer = bytes.Buffer{}
		}
	}
	if buffer.Len() > 0 {
		file.WriteString(buffer.String())
	}

	this.writingName = this.getFileName()

	return nil
}

func (this *DataLogger) Close() {
	fmt.Println("closed log", this.name)
	this.Flush()
	this.quit <- true
	delFlushLog(this)
}

func (this *DataLogger) autoFlush() {
	go func() {
		timer := time.NewTicker(time.Second)
		for {
			select {
			case <-timer.C:
				func() {
					lineCount := len(this.lines)
					if (lineCount >= this.cacheLineCount) || (lineCount > 0 && this.writingName != this.getFileName()) {
						if err := this.Flush(); err != nil {
							panic(err)
						}
					}
				}()
			case <-this.quit:
				break
			}
		}
	}()
	addFlushLog(this)
}

var flushLogMaps = make(map[string]*DataLogger, 16)

func addFlushLog(logger *DataLogger) error {
	if _, exists := flushLogMaps[logger.name]; exists {
		return errors.New("exists flush log")
	}

	flushLogMaps[logger.name] = logger
	return nil
}

func delFlushLog(logger *DataLogger) error {
	delete(flushLogMaps, logger.name)
	return nil
}

func CloseAllLogs() {
	for _, logger := range flushLogMaps {
		defer logger.Close()
	}
}
