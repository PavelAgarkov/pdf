package logger

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

const (
	logDir      = "./log"
	PanicLog    = "./log/panic.log"
	ErrLog      = "./log/error.log"
	WarningLog  = "./log/warning.log"
	InfoLog     = "./log/info.log"
	FrontendLog = "./log/frontend.log"
)

const (
	PanicName    = "panic"
	ErrorName    = "error"
	WarningName  = "warning"
	InfoName     = "info"
	FrontendName = "frontend"
)

type Interface interface {
	FlushLogs(logger *Factory)
	GetLogger(name string) *zap.SugaredLogger
}

type Factory struct {
	panicLogger    *zap.SugaredLogger
	errLogger      *zap.SugaredLogger
	warningLogger  *zap.SugaredLogger
	infoLogger     *zap.SugaredLogger
	frontendLogger *zap.SugaredLogger
}

func GetMapLogger() map[string]string {
	return map[string]string{
		PanicName:    PanicLog,
		ErrorName:    ErrLog,
		WarningName:  WarningLog,
		InfoName:     InfoLog,
		FrontendName: FrontendLog,
	}
}

func GetLoggerFactory(name map[string]string) *Factory {
	l := Factory{}
	return l.initLogger(name)
}

func (l *Factory) GetLogger(name string) *zap.SugaredLogger {
	switch name {
	case PanicName:
		return l.panicLogger
	case ErrorName:
		return l.errLogger
	case WarningName:
		return l.warningLogger
	case InfoName:
		return l.infoLogger
	case FrontendName:
		return l.frontendLogger
	default:
		return nil
	}
}

func (l *Factory) FlushLogs(logger *Factory) {
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

func (l *Factory) initLogger(name map[string]string) *Factory {
	return &Factory{
		panicLogger:    l.initLoggerLevel(false, l.openLogFile(name[PanicName], logDir), zap.PanicLevel),
		errLogger:      l.initLoggerLevel(false, l.openLogFile(name[ErrorName], logDir), zap.ErrorLevel),
		warningLogger:  l.initLoggerLevel(false, l.openLogFile(name[WarningName], logDir), zap.WarnLevel),
		infoLogger:     l.initLoggerLevel(false, l.openLogFile(name[InfoName], logDir), zap.InfoLevel),
		frontendLogger: l.initLoggerLevel(false, l.openLogFile(name[FrontendName], logDir), zap.ErrorLevel),
	}
}

func (*Factory) initLoggerLevel(d bool, f *os.File, level zapcore.Level) *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

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
