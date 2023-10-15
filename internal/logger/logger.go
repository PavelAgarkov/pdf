package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pdf/internal"
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
	logFile := l.openLogFile(filepath.FromSlash(internal.PanicLog), filepath.FromSlash(internal.LogDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error() + " panic")
		}
		<-l.panicSem
	}(logFile)
	log.SetOutput(logFile)

	l.withStackPrint(withStack, logText)
}

func (l *Factory) ErrorLog(logText string, withStack string) {
	l.errSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(internal.ErrLog), filepath.FromSlash(internal.LogDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error() + " error")
		}
		<-l.errSem
	}(logFile)
	log.SetOutput(logFile)

	l.withStackPrint(withStack, logText)
}

func (l *Factory) WarningLog(logText string) {
	l.warnSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(internal.WarningLog), filepath.FromSlash(internal.LogDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error() + " warn")
		}
		<-l.warnSem
	}(logFile)
	log.SetOutput(logFile)
	log.Println(logText)
}

func (l *Factory) InfoLog(logText string) {
	l.infoSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(internal.InfoLog), filepath.FromSlash(internal.LogDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error() + " info")
		}
		<-l.infoSem
	}(logFile)
	log.SetOutput(logFile)
	log.Println(logText)
}

func (l *Factory) FrontendLog(logText string) {
	l.frontSem <- 1
	logFile := l.openLogFile(filepath.FromSlash(internal.FrontendLog), filepath.FromSlash(internal.LogDir))
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println(err.Error() + " front")
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
