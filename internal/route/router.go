package route

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"pdf/internal/adapter"
	"pdf/internal/controller"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/storage"
)

const FilesPath = "./files/"
const FaviconFile = "./pdf-frontend/dist/favicon.ico"
const FrontendAssets = "./pdf-frontend/dist/assets/"

func Router(
	ctx context.Context,
	app *fiber.App,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *adapter.Locator,
	loggerFactory *logger.Factory,
) {
	bc := controller.NewBaseController()
	app.Get("/download/", controller.NewFileController(bc).Handle(ctx, filepath.FromSlash(FilesPath), loggerFactory)).
		Name("download")

	app.Get("/cancel/", controller.NewCancelController(bc).Handle(ctx, filepath.FromSlash(FilesPath), loggerFactory)).
		Name("cancel")

	app.Post("/merge/", controller.NewMergeController(bc).Handle(ctx, loggerFactory)).
		Name("merge")
}
