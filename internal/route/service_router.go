package route

import "github.com/gofiber/fiber/v2"

func ServiceRouter(app *fiber.App) {
	// root render vue
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	}).Name("root")

	// assets for vue
	app.Static("/assets/", FrontendAssets).Name("assets")
}
