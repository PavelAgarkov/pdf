package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"pdf/internal"
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
		defer RestoreController(loggerFactory, c, "cancel controller")

		operationData, hit := operationStorage.Get(internal.Hash2lvl(c.Cookies(internal.AuthenticationKey)))
		if !hit {
			errMsg := fmt.Sprintf("cancel controller: can't find hit %s from storage", internal.Hash2lvl(c.Cookies(internal.AuthenticationKey)))
			loggerFactory.ErrorLog(errMsg, "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		cookiesAdapter := adapterLocator.Locate(adapter.CookiesAlias).(*adapter.CookiesAdapter)
		ok, err := cookiesAdapter.IsAuthenticated(
			operationData.GetUserData().GetHash2Lvl(),
			internal.Hash1lvl(c.Cookies(internal.AuthenticationKey)),
			hit,
		)
		if err != nil {
			errMsg := fmt.Sprintf("cancel controller: can't delete %s from storage", internal.Hash2lvl(c.Cookies(internal.AuthenticationKey)))
			loggerFactory.ErrorLog(fmt.Sprintf(errMsg+" %s", err.Error()), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}
		if !ok {
			errMsg := fmt.Sprintf("cancel controller: can't acces to %s files by hash", internal.Hash2lvl(c.Cookies(internal.AuthenticationKey)))
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
