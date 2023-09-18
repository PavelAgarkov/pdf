package route

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"pdf/internal/adapter"
	"pdf/internal/controller"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
)

const FilesPath = "./files/"
const FaviconFile = "./pdf-frontend/dist/favicon.ico"
const FrontendAssets = "./pdf-frontend/dist/assets/"

func Router(
	ctx context.Context,
	app *fiber.App,
	operationStorage *pdf_operation.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *adapter.Locator,
	loggerFactory logger.Logger,
) {
	app.Get("/download/:filename", controller.GetFC().FileController(ctx, FilesPath, loggerFactory)).
		Name("file-download")
}
