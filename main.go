package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"pdf/internal/logger"
	"pdf/internal/route"
	"sync"
	"syscall"
)

const FrontendDist = "./pdf-frontend/dist"
const address = ":3000"

func main() {
	runServer()
}

func runServer() {
	engine := html.New(FrontendDist, ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	loggerFactory := logger.GetLoggerFactory(logger.GetMapLogger())
	defer recoveryFunction(loggerFactory)
	defer cleanupTasks(loggerFactory)
	defer loggerFactory.FlushLogs(loggerFactory)

	route.ServiceRouter(app)
	route.Router(ctx, app, loggerFactory)
	route.Middleware(app, loggerFactory)

	var serverShutdown sync.WaitGroup
	go func() {
		_ = <-sig
		serverShutdown.Add(1)
		_ = app.ShutdownWithContext(ctx)
		cancel()
		serverShutdown.Done()
		loggerFactory.GetLogger(logger.ErrorName).Error("Gracefully shutting down... Server STOPPED")
		return
	}()

	if err := app.Listen(address); err != nil {
		loggerFactory.
			GetLogger(logger.PanicName).
			With(zap.Stack("stackTrace")).
			Panic(fmt.Sprintf("server is stopped by error %s", err.Error()))
		return
	}

	serverShutdown.Wait()

	return
}

func recoveryFunction(loggerFactory logger.Logger) {
	if r := recover(); r != nil {
		loggerFactory.GetLogger(logger.ErrorName).Error("Recovered. Error:\n", r)
	}
}

func cleanupTasks(loggerFactory logger.Logger) {
	loggerFactory.GetLogger(logger.ErrorName).Error("Running cleanup tasks...")
}
