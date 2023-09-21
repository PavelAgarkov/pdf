package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"pdf/internal/controller"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
)

func Middleware(app *fiber.App, operationStorage *pdf_operation.OperationStorage, factory logger.Logger) {
	faviconMiddleware(app)
	recoveryHandleRequestMiddleware(app, factory)
	routs404RedirectMiddleware(app)
}

func faviconMiddleware(app *fiber.App) {
	app.Use(favicon.New(favicon.Config{
		File: FaviconFile,
		URL:  "/favicon.ico",
	}))
}

func recoveryHandleRequestMiddleware(app *fiber.App, loggerFactory logger.Logger) {
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
