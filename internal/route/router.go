package route

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"pdf/internal/controller"
	"pdf/internal/logger"
)

const FilesPath = "./files/"
const FaviconFile = "./pdf-frontend/dist/favicon.ico"
const FrontendAssets = "./pdf-frontend/dist/assets/"

func Router(app *fiber.App, loggerFactory *logger.LoggerFactory) {
	app.Get("/download/:filename", controller.GetFC().FileController(FilesPath, loggerFactory)).
		Name("file-download")
}

func ServiceRouter(app *fiber.App) {
	// root render vue
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	}).Name("root")

	// assets for vue
	app.Static("/assets/", FrontendAssets).Name("assets")
}

func Middleware(app *fiber.App, loggerFactory *logger.LoggerFactory) {
	// favicon middleware
	app.Use(favicon.New(favicon.Config{
		File: FaviconFile,
		URL:  "/favicon.ico",
	}))

	app.Use(func(c *fiber.Ctx) error {
		c.Context()
		defer func() {
			if r := recover(); r != nil {
				errStr := fmt.Sprintf("Recovered. Error: %s\n", r)
				fmt.Println(errStr)
				log.Panic(errStr)
				err := c.RedirectToRoute("root", map[string]interface{}{})
				if err != nil {
					return
				}
			}
		}()
		return c.Next()
	})

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
