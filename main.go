package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"os"
	"os/signal"
	"pdf/internal/logger"
	"pdf/internal/route"
	"sync"
	"syscall"
	"time"
)

const FrontendDist = "./pdf-frontend/dist"
const address = ":3000"

func main() {
	runServer()
}

func runServer() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	defer cleanupTasks()

	engine := html.New(FrontendDist, ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	var serverShutdown sync.WaitGroup
	go func() {
		_ = <-sig
		serverShutdown.Add(1)
		fmt.Println("Gracefully shutting down...")
		_ = app.ShutdownWithTimeout(1 * time.Second)
		fmt.Println("Wait timeout")
		serverShutdown.Done()
	}()

	loggerFactory := logger.GetLoggerFactory(
		logger.PanicLog,
		logger.ErrLog,
		logger.WarningLog,
		logger.InfoLog,
		logger.FrontendLog,
	)
	defer loggerFactory.FlushLogs(loggerFactory)

	route.ServiceRouter(app)
	route.Router(app, loggerFactory)
	route.Middleware(app, loggerFactory)

	if err := app.Listen(address); err != nil {
		log.Panic(err)
		errStr := fmt.Sprintf("server is stopped by error %s", err.Error())
		fmt.Println(errStr)
		cleanupTasks()
		panic(errStr)
		return
	}

	serverShutdown.Wait()

	return
}

func cleanupTasks() {
	fmt.Println("Running cleanup tasks...")
}
