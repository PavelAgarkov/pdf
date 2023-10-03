package controller

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/storage"
)

type RemovePageController struct {
	bc *BaseController
}

type RemovePageResponse struct {
	str string
	err error
}

func NewRemovePageController(bc *BaseController) *RemovePageController {
	return &RemovePageController{
		bc: bc,
	}
}

func (r *RemovePageResponse) GetStr() string {
	return r.str
}

func (r *RemovePageResponse) GetErr() error {
	return r.err
}

func (sc *RemovePageController) Handle(
	ctx context.Context,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "remove page controller")
		return nil
	}
}
