package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"os"
	"os/signal"
	"pdf/internal/route"
	"sync"
	"syscall"
)

const FrontendDist = "./pdf-frontend/dist"
const FrontendAssets = "./pdf-frontend/dist/assets/"
const FaviconFile = "./pdf-frontend/dist/favicon.ico"
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
	app.Use(func(c *fiber.Ctx) error {
		c.Context()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered. Error:\n", r)
				err := c.RedirectToRoute("root", map[string]interface{}{})
				if err != nil {
					return
				}
			}
		}()
		return c.Next()
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL)

	var serverShutdown sync.WaitGroup

	go func() {
		_ = <-sig
		fmt.Println("Gracefully shutting down...")
		serverShutdown.Add(1)
		_ = app.Shutdown()
		fmt.Println("Wait timeout")
		serverShutdown.Done()
	}()

	service(app)

	if err := app.Listen(address); err != nil {
		log.Panic(err)
	}

	serverShutdown.Wait()

	fmt.Println("Running cleanup tasks...")
}

func service(app *fiber.App) {
	// favicon middleware
	app.Use(favicon.New(favicon.Config{
		File: FaviconFile,
		URL:  "/favicon.ico",
	}))

	// root render vue
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	}).Name("root")

	// assets for vue
	app.Static("/assets/", FrontendAssets).Name("assets")

	route.Router(app)

	app.Use(func(c *fiber.Ctx) error {
		rout := c.Route()
		routName := c.Route().Name
		routes := app.GetRoutes()
		nameSet := make(map[string]struct{})

		for _, v := range routes {
			nameSet[v.Name] = struct{}{}
		}

		_, ok := nameSet[routName]
		if !ok || rout.Name == "" {
			return c.RedirectToRoute("root", map[string]interface{}{})
		}
		return c.Next()
	})
}
