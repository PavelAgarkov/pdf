package route

import (
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"pdf/internal"
)

func ServiceRouter(app *fiber.App) {
	// root render vue
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	}).Name("root")

	// assets for vue
	app.Static("/assets/", filepath.FromSlash(internal.FrontendAssets)).Name("assets")
}
