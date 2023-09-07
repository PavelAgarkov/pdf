package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"pdf/internal/route"
)

const FrontendDist = "./frontend/dist"
const FrontendAssets = "./frontend/dist/assets/"
const FaviconFile = "./frontend/dist/favicon.ico"
const address = ":3000"

func main() {
	engine := html.New(FrontendDist, ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

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

	// not found routs redirect to root
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

	log.Fatal(app.Listen(address))
}
