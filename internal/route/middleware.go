package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"path/filepath"
	"pdf/internal/controller"
	"pdf/internal/logger"
	"pdf/internal/storage"
)

func Middleware(app *fiber.App, operationStorage *storage.OperationStorage, loggerFactory *logger.Factory) {
	faviconMiddleware(app)
	corsMiddleware(app)
	recoveryHandleRequestMiddleware(app, loggerFactory)
	routs404RedirectMiddleware(app)
}

func faviconMiddleware(app *fiber.App) {
	app.Use(favicon.New(favicon.Config{
		File: FaviconFile,
		URL:  filepath.FromSlash("/favicon.ico"),
	}))
}

func corsMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "localhost",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
}

func recoveryHandleRequestMiddleware(app *fiber.App, loggerFactory *logger.Factory) {
	app.Use(func(c *fiber.Ctx) error {
		defer controller.RestoreController(loggerFactory, c)
		return c.Next()
	})
}

func routs404RedirectMiddleware(app *fiber.App) {
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
