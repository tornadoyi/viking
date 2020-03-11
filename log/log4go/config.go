package log4go

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type IConsoleConfig interface {
	Enable() bool
	Level() string
	Pattern() string
}

type IFileConfig interface {
	Enable()   bool
	Category() string
	Level()    string
	Filename() string

	// %T - Time (15:04:05 MST)
	// %t - Time (15:04)
	// %D - Date (2006/01/02)
	// %d - Date (01/02/06)
	// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
	// %S - Source
	// %M - Message
	// %C - Category
	// It ignores unknown format strings (and removes them)
	// Recommended: "[%D %T] [%C] [%L] (%S) %M"//
	Pattern() string

	Rotate()   bool
	Maxsize()  string
	Maxlines() string
	Daily()    bool
	Sanitize() bool
}

type ISocketConfig interface {
	Enable()   bool
	Category() string
	Level()    string
	Pattern()  string

	Addr()     string
	Protocol() string
}

type ILogConfig interface {
	Console() IConsoleConfig
	Files()   []IFileConfig
	Sockets() []ISocketConfig
}

func (log Logger) InitWithConfig(lc ILogConfig){
	if lc.Console().Enable() {
		filt, _ := configToConsoleLogWriter(lc.Console())
		log["stdout"] = &Filter{getLogLevel(lc.Console().Level()), filt, "DEFAULT"}
	}

	files := lc.Files()
	for _, fc := range files {
		if !fc.Enable() {
			continue
		}
		if len(fc.Category()) == 0 {
			panic(errors.New("Logger init failed, file category can not be empty "))
		}

		filt, _ := configToFileLogWriter(fc)
		log[fc.Category()] = &Filter{getLogLevel(fc.Level()), filt, fc.Category()}
	}

	sockets := lc.Sockets()
	for _, sc := range sockets {
		if !sc.Enable() {
			continue
		}
		if len(sc.Category()) == 0 {
			panic(errors.New("Logger init failed, file category can not be empty"))
		}

		filt, _ := configToSocketLogWriter(sc)
		log[sc.Category()] = &Filter{getLogLevel(sc.Level()), filt, sc.Category()}
	}
}


func configToConsoleLogWriter(cf IConsoleConfig) (*ConsoleLogWriter, bool) {
	format := "[%D %T] [%C] [%L] (%S) %M"

	if len(cf.Pattern()) > 0 {
		format = strings.Trim(cf.Pattern(), " \r\n")
	}

	if !cf.Enable() {
		return nil, true
	}

	clw := NewConsoleLogWriter()
	clw.SetFormat(format)

	return clw, true
}

func configToFileLogWriter(ff IFileConfig) (*FileLogWriter, bool) {
	file := "app.log"
	format := "[%D %T] [%C] [%L] (%S) %M"
	maxlines := 0
	maxsize := 0
	daily := false
	rotate := false
	sanitize := false

	if len(ff.Filename()) > 0 {
		file = ff.Filename()
	}
	if len(ff.Pattern()) > 0 {
		format = strings.Trim(ff.Pattern(), " \r\n")
	}
	if len(ff.Maxlines()) > 0 {
		maxlines = strToNumSuffix(strings.Trim(ff.Maxlines(), " \r\n"), 1000)
	}
	if len(ff.Maxsize()) > 0 {
		maxsize = strToNumSuffix(strings.Trim(ff.Maxsize(), " \r\n"), 1024)
	}
	daily = ff.Daily()
	rotate = ff.Rotate()
	sanitize = ff.Sanitize()

	if !ff.Enable() {
		return nil, true
	}

	flw := NewFileLogWriter(file, rotate, daily)
	flw.SetFormat(format)
	flw.SetRotateLines(maxlines)
	flw.SetRotateSize(maxsize)
	flw.SetSanitize(sanitize)
	return flw, true
}

func configToSocketLogWriter(sf ISocketConfig) (SocketLogWriter, bool) {
	endpoint := ""
	protocol := "tcp"

	if len(sf.Addr()) == 0 {
		panic(errors.New("Error: Required property \"addr\" for file filter missing"))
	}
	endpoint = sf.Addr()

	// set socket protocol
	if len(sf.Protocol()) > 0 {
		if sf.Protocol() != "tcp" && sf.Protocol() != "udp" {
			fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Required property \"%s\" for file filter wrong type, use default tcp instead.\n", "protocol")
		} else {
			protocol = sf.Protocol()
		}
	}

	if !sf.Enable() {
		return nil, true
	}

	return NewSocketLogWriter(protocol, endpoint), true
}