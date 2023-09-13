package logger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

const logDir = "./log"
const PanicLog = "./log/panic.log"
const ErrLog = "./log/error.log"
const WarningLog = "./log/warning.log"
const InfoLog = "./log/info.log"
const FrontendLog = "./log/frontend.log"

type LoggerFactory struct {
	panicLogger    *zap.SugaredLogger
	errLogger      *zap.SugaredLogger
	warningLogger  *zap.SugaredLogger
	infoLogger     *zap.SugaredLogger
	frontendLogger *zap.SugaredLogger
}

func GetLoggerFactory(panicFile, errorFile, warningFile, infoFile, frontendFile string) *LoggerFactory {
	l := LoggerFactory{}
	return l.initLogger(panicFile, errorFile, warningFile, infoFile, frontendFile)
}

func (l *LoggerFactory) GetPanicLogger() *zap.SugaredLogger {
	return l.panicLogger
}

func (l *LoggerFactory) GetErrorLogger() *zap.SugaredLogger {
	return l.errLogger
}

func (l *LoggerFactory) GetWarningLogger() *zap.SugaredLogger {
	return l.warningLogger
}

func (l *LoggerFactory) GetInfoLogger() *zap.SugaredLogger {
	return l.infoLogger
}

func (l *LoggerFactory) GetFrontendLogger() *zap.SugaredLogger {
	return l.frontendLogger
}

func (*LoggerFactory) FlushLogs(logger *LoggerFactory) {
	fmt.Println("Flush logger")
	err := logger.panicLogger.Sync()
	if err != nil {
		fmt.Println("failed buffered panic-logger flash " + err.Error())
	}
	err = logger.errLogger.Sync()
	if err != nil {
		fmt.Println("failed buffered err-logger flash " + err.Error())
	}
	err = logger.warningLogger.Sync()
	if err != nil {
		fmt.Println("failed buffered warning-logger flash " + err.Error())
	}
	err = logger.infoLogger.Sync()
	if err != nil {
		fmt.Println("failed buffered info-logger flash " + err.Error())
	}
	err = logger.frontendLogger.Sync()
	if err != nil {
		fmt.Println("failed buffered frontend-logger flash " + err.Error())
	}
}

func (l *LoggerFactory) initLogger(panicFile, errorFile, warningFile, infoFile, frontendFile string) *LoggerFactory {
	return &LoggerFactory{
		panicLogger:    l.initLoggerLevel(false, l.openLogFile(panicFile, logDir), zap.PanicLevel),
		errLogger:      l.initLoggerLevel(false, l.openLogFile(errorFile, logDir), zap.ErrorLevel),
		warningLogger:  l.initLoggerLevel(false, l.openLogFile(warningFile, logDir), zap.WarnLevel),
		infoLogger:     l.initLoggerLevel(false, l.openLogFile(infoFile, logDir), zap.InfoLevel),
		frontendLogger: l.initLoggerLevel(false, l.openLogFile(frontendFile, logDir), zap.ErrorLevel),
	}
}

func (*LoggerFactory) initLoggerLevel(d bool, f *os.File, level zapcore.Level) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	fileEncoder := zapcore.NewJSONEncoder(pe)
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	if d {
		level = zap.DebugLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)
	logger := zap.New(core)

	return logger.Sugar()
}

func (*LoggerFactory) openLogFile(fname string, dirname string) *os.File {
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
