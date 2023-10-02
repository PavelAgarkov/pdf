package route

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"pdf/internal/controller"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/storage"
)

func Router(
	ctx context.Context,
	app *fiber.App,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) {
	bc := controller.NewBaseController()
	app.Get("/download/", controller.NewDownloadController(bc).Handle(
		operationStorage,
		adapterLocator,
		loggerFactory,
	)).Name("download")

	app.Get("/cancel/", controller.NewCancelController(bc).Handle(
		operationStorage,
		adapterLocator,
		loggerFactory,
	)).Name("cancel")

	app.Post("/merge/", controller.NewMergeController(bc).Handle(
		ctx,
		operationStorage,
		operationFactory,
		adapterLocator,
		loggerFactory,
	)).Name("merge")
}
