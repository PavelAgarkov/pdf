package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"pdf/internal"
	"pdf/internal/adapter"
	"pdf/internal/locator"
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
	runHTTPServer()
}

func runHTTPServer() {
	pathAdapter := adapter.NewPathAdapter()
	adapterLocator := locator.NewAdapterLocator(
		adapter.NewFileAdapter(),
		pathAdapter,
		adapter.NewPdfAdapter(),
		adapter.NewArchiveAdapter(pathAdapter),
	)
	engine := html.New(adapter.GenerateFrontendDist(), ".html")
	app := fiber.New(fiber.Config{Views: engine})

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())
	loggerFactory := logger.NewLoggerFactory()
	defer recoveryFunction(loggerFactory)

	operationStorage := storage.NewInMemoryOperationStorage()
	operationStorage.Run(ctx, internal.Timer5*internal.Minute, adapterLocator, loggerFactory)

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
	route.Middleware(app, loggerFactory)

	go func() {
		_ = <-sig
		cancel()
	}()

	var serverShutdown sync.WaitGroup
	go func() {
		defer recoveryFunction(loggerFactory)
		<-ctx.Done()
		cleanupTasks(operationStorage, adapterLocator, loggerFactory)
		serverShutdown.Add(1)
		_ = app.ShutdownWithContext(ctx)
		serverShutdown.Done()
		loggerFactory.InfoLog("Gracefully shutting down... Server STOPPED")
	}()

	if err := app.Listen(address); err != nil {
		loggerFactory.PanicLog(
			fmt.Sprintf("server is stopped by error %s", err.Error()),
			zap.Stack("").String,
		)
		return
	}

	serverShutdown.Wait()

	return
}

func recoveryFunction(loggerFactory *logger.Factory) {
	if r := recover(); r != nil {
		loggerFactory.PanicLog(fmt.Sprintf("Recovered. Error:\n", r), "")
	}
}

func cleanupTasks(
	operationStorage *storage.OperationStorage,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) {
	loggerFactory.InfoLog("Running cleanup tasks...")
	operationStorage.ClearStorageAndFilesystem(adapterLocator, loggerFactory)
	loggerFactory.InfoLog("Running cleanup tasks done")
}
