package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"os"
	"os/signal"
	"pdf/internal/route"
	"sync"
	"syscall"
	"time"
)

const FrontendDist = "./pdf-frontend/dist"
const address = ":3000"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
		}
	}()
	runServer()
}

func runServer() {
	engine := html.New(FrontendDist, ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	var serverShutdown sync.WaitGroup

	go func() {
		_ = <-sig
		fmt.Println("Gracefully shutting down...")
		serverShutdown.Add(1)
		fmt.Println("Gracefully shutting down...")
		_ = app.ShutdownWithTimeout(1 * time.Second)
		fmt.Println("Wait timeout")
		serverShutdown.Done()
		fmt.Println("Wait timeout")
	}()

	route.ServiceRouter(app)
	route.Router(app)
	route.Middleware(app)

	if err := app.Listen(address); err != nil {
		log.Panic(err)
		errStr := fmt.Sprintf("server is stopped by error %s", err.Error())
		fmt.Println(errStr)
		cleanupTasks()
		panic(errStr)
		return
	}

	serverShutdown.Wait()

	cleanupTasks()
}

func cleanupTasks() {
	fmt.Println("Running cleanup tasks...")
}
