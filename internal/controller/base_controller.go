package controller

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"pdf/internal"
	"pdf/internal/logger"
	"slices"
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
