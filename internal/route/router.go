package route

import (
	"github.com/gofiber/fiber/v2"
	"pdf/internal/controller"
)

const FilesPath = "./files/"

func Router(app *fiber.App) {
	app.Get("/download/:filename", controller.GetFC().FileController(FilesPath)).
		Name("file-download")
}
