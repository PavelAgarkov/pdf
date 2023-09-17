package route

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"pdf/internal/controller"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/service"
	"pdf/internal/storage"
)

const FilesPath = "./files/"
const FaviconFile = "./pdf-frontend/dist/favicon.ico"
const FrontendAssets = "./pdf-frontend/dist/assets/"

func Router(
	ctx context.Context,
	app *fiber.App,
	userStorage *storage.UserStorage,
	pdfAdapter *service.PdfAdapter,
	operationFactory *pdf_operation.OperationsFactory,
	factory logger.Logger,
) {
	app.Get("/download/:filename", controller.GetFC().FileController(ctx, FilesPath, factory)).
		Name("file-download")
}
