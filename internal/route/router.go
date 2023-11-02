package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"pdf/internal/controller"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/storage"
	"time"
)

func Router(
	app *fiber.App,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) {
	v1group := app.Group("/v1")
	bc := controller.NewBaseController()

	v1Router(
		operationStorage,
		operationFactory,
		adapterLocator,
		loggerFactory,
		bc,
		v1group,
	)
}

func v1Router(
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
	bc *controller.BaseController,
	v1group fiber.Router,
) {
	v1group.Post("/frontend-log-write/", controller.NewFrontendLogController(bc).Handle(
		loggerFactory,
	)).Name("frontend-log-write")

	v1group.Get("/download/", controller.NewDownloadController(bc).Handle(
		operationStorage,
		adapterLocator,
		loggerFactory,
	)).Name("download")

	v1group.Get("/cancel/", controller.NewCancelController(bc).Handle(
		operationStorage,
		adapterLocator,
		loggerFactory,
	)).Name("cancel")

	v1group.Post("/merge/",
		timeout.NewWithContext(
			controller.NewMergeController(bc).Handle(
				//ctx,
				operationStorage,
				operationFactory,
				adapterLocator,
				loggerFactory,
			), 200*time.Second,
		),
	).Name("merge")

	v1group.Post("/split-page/",
		timeout.NewWithContext(
			controller.NewSplitPageController(bc).Handle(
				//ctx,
				operationStorage,
				operationFactory,
				adapterLocator,
				loggerFactory,
			),
			200*time.Second,
		),
	).Name("split-page")

	v1group.Post("/remove-pages/",
		timeout.NewWithContext(
			controller.NewRemovePageController(bc).Handle(
				//ctx,
				operationStorage,
				operationFactory,
				adapterLocator,
				loggerFactory,
			),
			200*time.Second,
		),
	).Name("remove-pages")
}
