package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"pdf/internal"
	"pdf/internal/hash"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/service"
	"pdf/internal/storage"
	"slices"
	"strconv"
	"strings"
)

type BaseController struct{}

type ResponseInterface interface {
	GetStr() string
	GetErr() error
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (bc *BaseController) SelectResult(
	ctx context.Context,
	ch chan ResponseInterface,
	start chan struct{},
) ResponseInterface {
	start <- struct{}{}
	select {
	case <-ctx.Done():
		res, ok := <-ch
		if !ok {
			return nil
		}
		return res
	case res := <-ch:
		return res
	}
}

func RestoreController(loggerFactory *logger.Factory, destination string) {
	if r := recover(); r != nil {
		panicStr := fmt.Sprintf(destination+" : Recovered. Panic: %s\n", r)
		loggerFactory.PanicLog(panicStr, "")
	}
}

func (bc *BaseController) isAuthenticated(
	operationStorage *storage.OperationStorage,
	c *fiber.Ctx,
	loggerFactory *logger.Factory,
) (pdf_operation.OperationDataInterface, error) {
	authToken := service.ParseBearerHeader(c.GetReqHeaders()[internal.AuthenticationHeader])
	operationData, hit := operationStorage.Get(hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true))
	if !hit {
		errMsg := fmt.Sprintf("can't find hit %s from storage", authToken)
		loggerFactory.ErrorLog(errMsg, "")
		return nil, errors.New(errMsg)
	}

	ok, err := service.IsAuthenticated(operationData.GetUserData().GetHash2Lvl(), internal.Hash1lvl(authToken))
	if err != nil {
		errMsg := fmt.Sprintf("can't acces %s to storage", authToken)
		loggerFactory.ErrorLog(fmt.Sprintf(errMsg+" %s", err.Error()), "")
		return nil, errors.New(errMsg)
	}
	if !ok {
		errMsg := fmt.Sprintf("can't acces to %s files by hash", authToken)
		loggerFactory.ErrorLog(errMsg, "")
		return nil, errors.New(errMsg)
	}

	return operationData, nil
}

func (bc *BaseController) isOverAuthenticated(
	operationStorage *storage.OperationStorage,
	c *fiber.Ctx,
	loggerFactory *logger.Factory,
) error {
	authToken := service.ParseBearerHeader(c.GetReqHeaders()[internal.AuthenticationHeader])
	_, hit := operationStorage.Get(hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true))
	if hit {
		errMsg := fmt.Sprintf("can't process %s already in storage", authToken)
		loggerFactory.ErrorLog(errMsg, "")
		return errors.New(errMsg)
	}

	return nil
}

func (bc *BaseController) formValidation(form *multipart.Form) error {
	var sumSize int64 = 0
	countFiles := 0
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			sumSize += fileHeader.Size
			countFiles++
		}
	}

	archiveFormatSlice, ok := form.Value[internal.ArchiveFormatKeyForRequest]
	archiveFormat := ""
	if len(archiveFormatSlice) > 0 {
		archiveFormat = archiveFormatSlice[0]
	}

	formats := []string{
		internal.ZipFormat,
		internal.ZipZstFormat,
		internal.TarFormat,
		internal.TarGzFormat,
	}

	if !ok || !slices.Contains(formats, archiveFormat) {
		return errors.New("archive format must be selected from the list")
	}

	if sumSize > internal.MaxSumUploadFilesSizeByte {
		return errors.New("upload files must be less 100Mb")
	}

	if countFiles > internal.MaxNumberUploadFiles {
		return errors.New("number upload files must be less 100")
	}

	return nil
}

func (bc *BaseController) alphaSymbolValidation(form *multipart.Form, key string) error {
	const alpha = "1234567890-"
	for _, interval := range form.Value[key] {
		for _, char := range interval {
			if !strings.Contains(alpha, strings.ToLower(string(char))) {
				return errors.New(fmt.Sprintf("invalid symbol, !%s!", string(char)))
			}
		}

		chunks := strings.Split(interval, "-")
		if len(chunks) > 2 {
			return errors.New("format must be some '2-5' or 5")
		}
	}
	return nil
}

func (bc *BaseController) orderIntervalValidation(form *multipart.Form, key string) error {
	_, intervals := internal.ParseIntervals(form.Value[key])
	for _, interval := range intervals {
		if len(interval) == 2 {
			if interval[0] > interval[1] {
				return errors.New("interval format 'n-n' must be written in ascending order")
			}
		}
	}
	return nil
}

func (bc *BaseController) numberFilesValidation(form *multipart.Form, must int) error {
	number := 0
	for _, fileHeaders := range form.File {
		number = len(fileHeaders)
		break
	}
	if number != must {
		return errors.New(fmt.Sprintf("for split operation must %s pdf files", strconv.Itoa(must)))
	}
	return nil
}
