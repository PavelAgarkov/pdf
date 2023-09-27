package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"pdf/internal/adapter"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/route"
	"pdf/internal/storage"
	"sync"
	"syscall"
)

const (
	address = ":3000"
)

func main() {
	runServer()
}

func runServer() {
	pathAdapter := adapter.NewPathAdapter()
	adapterLocator := adapter.NewAdapterLocator(
		adapter.NewFileAdapter(),
		pathAdapter,
		adapter.NewPdfAdapter(),
		adapter.NewArchiveAdapter(pathAdapter),
	)
	engine := html.New(adapter.GenerateFrontendDist(), ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())

	loggerFactory := logger.GetLoggerFactory()

	defer recoveryFunction(loggerFactory)
	defer cleanupTasks(loggerFactory)

	operationStorage := storage.NewInMemoryOperationStorage()
	operationStorage.Run(ctx, pdf_operation.Timer5)

	operationFactory := pdf_operation.NewOperationFactory()

	route.ServiceRouter(app)
	route.Router(
		ctx,
		app,
		operationStorage,
		operationFactory,
		adapterLocator,
		loggerFactory,
	)
	route.Middleware(app, operationStorage, loggerFactory)

	var serverShutdown sync.WaitGroup
	go func() {
		_ = <-sig
		serverShutdown.Add(1)
		_ = app.ShutdownWithContext(ctx)
		cancel()
		serverShutdown.Done()
		loggerFactory.ErrorLog("Gracefully shutting down... Server STOPPED", "")
		return
	}()

	if err := app.Listen(address); err != nil {
		loggerFactory.PanicLog(
			fmt.Sprintf("server is stopped by error %s", err.Error()),
			zap.Stack("stackTrace").String,
		)
		return
	}

	serverShutdown.Wait()

	return
}

func recoveryFunction(loggerFactory *logger.Factory) {
	if r := recover(); r != nil {
		loggerFactory.ErrorLog(fmt.Sprintf("Recovered. Error:\n", r), "")
	}
}

func cleanupTasks(loggerFactory *logger.Factory) {
	loggerFactory.ErrorLog("Running cleanup tasks...", "")
}
