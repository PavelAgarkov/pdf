package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"pdf/internal/adapter"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/storage"
)

type CancelController struct {
	bc *BaseController
}

type CancelResponse struct {
	str string
	err error
}

func (r *CancelResponse) GetStr() string {
	return r.str
}

func (r *CancelResponse) GetErr() error {
	return r.err
}

func NewCancelController(bc *BaseController) *CancelController {
	return &CancelController{
		bc: bc,
	}
}

func (cc *CancelController) Handle(
	operationStorage *storage.OperationStorage,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "cancel controller")

		operationData, authenticatedErr := cc.bc.isAuthenticated(operationStorage, c, loggerFactory)
		if authenticatedErr != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": authenticatedErr.Error(),
			})
		}

		pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
		rootDir := pathAdapter.GenerateRootDir(operationData.GetUserData().GetHash2Lvl())
		err := os.RemoveAll(string(rootDir))
		if err != nil {
			errMsg := fmt.Sprintf("cancel controller: can't remove dir %s", rootDir)
			loggerFactory.ErrorLog(fmt.Sprintf(errMsg+" %s", err.Error()), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		operationStorage.Delete(operationData.GetUserData().GetHash2Lvl())

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "no",
		})
	}
}
