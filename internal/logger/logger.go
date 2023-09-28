package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	logDir      = "./log"
	PanicLog    = "./log/panic.log"
	ErrLog      = "./log/error.log"
	WarningLog  = "./log/warning.log"
	InfoLog     = "./log/info.log"
	FrontendLog = "./log/frontend.log"
)

type Factory struct {
	errSem   chan int
	warnSem  chan int
	panicSem chan int
	infoSem  chan int
	frontSem chan int
}

func NewLoggerFactory() *Factory {
	return &Factory{
		errSem:   make(chan int, 1),
		warnSem:  make(chan int, 1),
		panicSem: make(chan int, 1),
		infoSem:  make(chan int, 1),
		frontSem: make(chan int, 1),
	}
}

func (l *Factory) PanicLog(logText string, withStack string) {
	l.panicSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(PanicLog), filepath.FromSlash(logDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		<-l.panicSem
	}(logFile)
	log.SetOutput(logFile)

	l.withStackPrint(withStack, logText)
}

func (l *Factory) ErrorLog(logText string, withStack string) {
	l.errSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(ErrLog), filepath.FromSlash(logDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		<-l.errSem
	}(logFile)
	log.SetOutput(logFile)

	l.withStackPrint(withStack, logText)
}

func (l *Factory) WarningLog(logText string) {
	l.warnSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(WarningLog), filepath.FromSlash(logDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		<-l.warnSem
	}(logFile)
	log.SetOutput(logFile)
	log.Println(logText)
}

func (l *Factory) InfoLog(logText string) {
	l.infoSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(InfoLog), filepath.FromSlash(logDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		<-l.infoSem
	}(logFile)
	log.SetOutput(logFile)
	log.Println(logText)
}

func (l *Factory) FrontendLog(logText string) {
	l.frontSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(FrontendLog), filepath.FromSlash(logDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		<-l.frontSem
	}(logFile)
	log.SetOutput(logFile)
	log.Println(logText)
}

func (l *Factory) openLogFile(fname string, dirname string) *os.File {
	if _, err := os.Stat(dirname); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dirname, 0777)
		if err != nil {
			log.Println(err)
		}
	}

	if _, err := os.Stat(fname); errors.Is(err, os.ErrNotExist) || err != nil {
		file, err := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
		}

		return file
	}
	file, _ := os.OpenFile(fname, os.O_WRONLY|os.O_APPEND, 0666)

	return file
}

func (l *Factory) withStackPrint(withStack, logText string) {
	if withStack != "" {
		log.Println(logText + " ;; STACK ;; " + withStack)
	} else {
		log.Println(logText)
	}
}
