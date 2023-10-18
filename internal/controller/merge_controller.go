package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"os"
	"pdf/internal"
	"pdf/internal/adapter"
	"pdf/internal/entity"
	"pdf/internal/hash"
	"pdf/internal/locator"
	"pdf/internal/logger"
	"pdf/internal/pdf_operation"
	"pdf/internal/service"
	"pdf/internal/storage"
	"pdf/internal/validation"
	"strings"
	"time"
)

const ChannelResponseOK = "ok"

type MergeController struct {
	bc *BaseController
}

type MergeResponse struct {
	str string
	err error
}

func NewMergeController(bc *BaseController) *MergeController {
	return &MergeController{
		bc: bc,
	}
}

func (r *MergeResponse) GetStr() string {
	return r.str
}

func (r *MergeResponse) GetErr() error {
	return r.err
}

func (mc *MergeController) Handle(
	ctx context.Context,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "merge controller")

		overAuthenticatedErr := mc.bc.isOverAuthenticated(operationStorage, c)
		if overAuthenticatedErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": overAuthenticatedErr.Error(),
			})
		}

		authToken := service.GenerateBearerToken()

		ctxC, cancel := context.WithTimeout(ctx, 75*time.Second)
		defer cancel()
		cr := make(chan ResponseInterface)
		start := make(chan struct{})

		go mc.realHandler(
			c,
			ctxC,
			start,
			cr,
			operationStorage,
			operationFactory,
			adapterLocator,
			authToken,
			loggerFactory,
		)
		res := mc.bc.SelectResponse(ctxC, cr, start)
		if res == nil {
			loggerFactory.PanicLog("merge controller: context expired", "")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "merge controller: context expired",
			})
		}

		if res.GetStr() == internal.ErrorForm {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": res.GetErr().Error(),
			})
		}

		if res.GetStr() == internal.OperationError {
			loggerFactory.ErrorLog(res.GetErr().Error(), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "operation_version_is_not_suitable",
			})
		}

		if res.GetStr() != ChannelResponseOK {
			loggerFactory.ErrorLog(res.GetErr().Error(), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": res.GetErr().Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"hash":  authToken,
			"error": "",
		})
	}
}

func (mc *MergeController) realHandler(
	c *fiber.Ctx,
	ctx context.Context,
	start chan struct{},
	cr chan ResponseInterface,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	authToken string,
	loggerFactory *logger.Factory,
) {
	<-start

	defer RestoreController(loggerFactory, "merge controller")

	secondLevelHash := hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	rootDir := pathAdapter.GenerateRootDir(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	archiveDir := pathAdapter.GenerateArchiveDirPath(secondLevelHash)

	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	err := fileAdapter.CreateDir(string(rootDir), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)
	err = fileAdapter.CreateDir(string(archiveDir), 0777)
	defer func() {
		_ = os.RemoveAll(string(inDir))
		_ = os.RemoveAll(string(outDir))
		cr <- &MergeResponse{
			str: "panic_context",
			err: errors.New("handle cancel"),
		}
	}()

	if err != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &MergeResponse{
			str: "cant_create_dir",
			err: fmt.Errorf("cant_read_form: %w", err),
		}
		return
	}

	form, errRead := c.MultipartForm()
	if errRead != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &MergeResponse{
			str: "cant_read_form",
			err: fmt.Errorf("cant_read_form: %w", errRead),
		}
		return
	}

	errForm := mc.formValidation(form)
	if errForm != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &MergeResponse{
			str: internal.ErrorForm,
			err: errForm,
		}
		return
	}

	// обработка файлов из формы
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			nameWithoutSpace := strings.ReplaceAll(fileHeader.Filename, " ", "_")
			_, pathToFile, _ := pathAdapter.StepForward(internal.Path(inDir), nameWithoutSpace)
			errSave := c.SaveFile(fileHeader, string(pathToFile))
			if errSave != nil {
				_ = os.RemoveAll(string(rootDir))
				cr <- &MergeResponse{
					str: "cant_save_file_from_form",
					err: fmt.Errorf("cant_save_file_from_form: %w", errSave),
				}
				return
			}
		}
	}

	//получения списка, с порядком файлов указанным пользователем
	orderFiles, _ := form.Value[internal.OrderFilesRequestKey]

	//переопределение путей до файловой системы
	for k, v := range orderFiles {
		nameWithoutSpace := strings.ReplaceAll(v, " ", "_")
		_, pathToFile, _ := pathAdapter.StepForward(internal.Path(inDir), nameWithoutSpace)
		orderFiles[k] = string(pathToFile)
	}

	userData := entity.NewUserData(
		internal.Hash1lvl(authToken),
		secondLevelHash,
		time.Now().Add(internal.Timer5*internal.Minute),
	)
	mergePagesOperation := operationFactory.CreateNewOperation(
		pdf_operation.NewConfiguration(nil, nil),
		userData,
		orderFiles,
		rootDir,
		inDir,
		outDir,
		archiveDir,
		"",
		pdf_operation.DestinationMerge,
	).(*pdf_operation.MergeOperation)

	archivePath, operationErr := mergePagesOperation.Execute(
		ctx,
		adapterLocator,
		form.Value[internal.ArchiveFormatKeyForRequest][0],
	)
	if operationErr != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &MergeResponse{
			str: internal.OperationError,
			err: fmt.Errorf("can't make opearation: %w", operationErr),
		}
		return
	}

	data := pdf_operation.NewOperationData(
		userData,
		internal.ArchiveDir(archivePath),
		mergePagesOperation.GetBaseOperation().GetStatus(),
		mergePagesOperation.GetBaseOperation().GetStoppedReason(),
	)

	operationStorage.Insert(secondLevelHash, data)

	cr <- &MergeResponse{str: ChannelResponseOK, err: nil}
	return
}

func (mc *MergeController) formValidation(form *multipart.Form) error {
	if _, ok := form.Value[internal.OrderFilesRequestKey]; !ok {
		return errors.New("form must contain the order in which the files will be merged")
	}

	if len(form.File) == 0 {
		return errors.New("for merge operation must than 1 pdf file")
	}

	number := 0
	for _, fileHeaders := range form.File {
		number = len(fileHeaders)
		break
	}

	if len(form.Value[internal.OrderFilesRequestKey]) != number {
		return errors.New("form must be contain merge order for all files")
	}

	err := validation.FormFileValidation(form)
	if err != nil {
		return err
	}

	return nil
}
