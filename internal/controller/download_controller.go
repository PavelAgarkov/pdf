package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"pdf/internal"
	"pdf/internal/adapter"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/storage"
)

type DownloadController struct {
	bc *BaseController
}

type DownloadResponse struct {
	str string
	err error
}

func (dr *DownloadResponse) GetStr() string {
	return dr.str
}

func (dr *DownloadResponse) GetErr() error {
	return dr.err
}

func NewDownloadController(bc *BaseController) *DownloadController {
	return &DownloadController{
		bc: bc,
	}
}

func (dc *DownloadController) Handle(
	operationStorage *storage.OperationStorage,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "download controller")

		c.Accepts("application/zip")
		c.Accepts("application/x-bzip")
		c.Accepts("application/x-tar")

		operationData, authenticatedErr := dc.bc.isAuthenticated(operationStorage, c)
		if authenticatedErr != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": authenticatedErr.Error(),
			})
		}

		archivePath := string(operationData.(*pdf_operation.OperationData).GetArchivePath())
		if archivePath == "" {
			errMsg := fmt.Sprintf("download controller: can't find archive")
			loggerFactory.ErrorLog(errMsg, "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
		_, file, _ := pathAdapter.StepBack(internal.Path(archivePath))
		rootDir := pathAdapter.GenerateRootDir(operationData.GetUserData().GetHash2Lvl())

		defer os.RemoveAll(string(rootDir))
		defer operationStorage.Delete(operationData.GetUserData().GetHash2Lvl())
		return c.Download(archivePath, file)
	}
}
