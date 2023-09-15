package route

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"pdf/internal/logger"
)

func Middleware(app *fiber.App, factory logger.Logger) {
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

func recoveryHandleRequestMiddleware(app *fiber.App, factory logger.Logger) {
	app.Use(func(c *fiber.Ctx) error {
		c.Context()
		defer func() {
			if r := recover(); r != nil {
				errStr := fmt.Sprintf("Recovered. Error: %s\n", r)
				fmt.Println(errStr)
				factory.GetLogger(logger.ErrorName).Error(errStr)
				err := c.RedirectToRoute("root", map[string]interface{}{})
				if err != nil {
					return
				}
			}
		}()
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
