package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"pdf/internal"
	"pdf/internal/adapter"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/service"
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

		authToken := service.ParseBearerHeader(c.GetReqHeaders()[internal.AuthenticationHeader])
		operationData, hit := operationStorage.Get(internal.Hash2lvl(authToken))
		if !hit {
			errMsg := fmt.Sprintf("cancel controller: can't find hit %s from storage", authToken)
			loggerFactory.ErrorLog(errMsg, "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		ok, err := service.IsAuthenticated(operationData.GetUserData().GetHash2Lvl(), internal.Hash1lvl(authToken))
		if err != nil {
			errMsg := fmt.Sprintf("cancel controller: can't delete %s from storage", authToken)
			loggerFactory.ErrorLog(fmt.Sprintf(errMsg+" %s", err.Error()), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		if !ok {
			errMsg := fmt.Sprintf("cancel controller: can't acces to %s files by hash", authToken)
			loggerFactory.ErrorLog(errMsg, "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
		rootDir := pathAdapter.GenerateRootDir(operationData.GetUserData().GetHash2Lvl())
		err = os.RemoveAll(string(rootDir))
		if err != nil {
			errMsg := fmt.Sprintf("cancel controller: can't remove dir %s", rootDir)
			loggerFactory.ErrorLog(fmt.Sprintf(errMsg+" %s", err.Error()), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		operationStorage.Delete(operationData.GetUserData().GetHash2Lvl())

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "ok",
		})
	}
}
