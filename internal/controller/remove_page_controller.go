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

func (rpc *RemovePageController) Handle(
	ctx context.Context,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "remove page controller")

		overAuthenticatedErr := rpc.bc.isOverAuthenticated(operationStorage, c)
		if overAuthenticatedErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": overAuthenticatedErr.Error(),
			})
		}

		authToken := service.GenerateBearerToken()

		ctxC, cancel := context.WithTimeout(ctx, 200*time.Second)
		defer cancel()
		cr := make(chan ResponseInterface)
		start := make(chan struct{})

		go rpc.realHandler(
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
		res := rpc.bc.SelectResponse(ctxC, cr, start)
		if res == nil {
			loggerFactory.PanicLog("remove page controller: context expired", "")
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

func (rpc *RemovePageController) realHandler(
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

	errForm := rpc.formValidation(form)
	if errForm != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &MergeResponse{
			str: internal.ErrorForm,
			err: errForm,
		}
		return
	}

	// обработка файлов из формы
	files := make([]string, 0)
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
			files = append(files, string(pathToFile))
		}
	}

	//получения списка, с порядком файлов указанным пользователем
	removePagesIntervals, _ := form.Value[internal.RemovePagesIntervals]

	userData := entity.NewUserData(
		internal.Hash1lvl(authToken),
		secondLevelHash,
		time.Now().Add(internal.Timer5*internal.Minute),
	)
	removePagesOperation := operationFactory.CreateNewOperation(
		pdf_operation.NewConfiguration(nil, removePagesIntervals),
		userData,
		files,
		rootDir,
		inDir,
		outDir,
		archiveDir,
		"",
		pdf_operation.DestinationRemovePages,
	).(*pdf_operation.RemovePagesOperation)

	archivePath, operationErr := removePagesOperation.Execute(
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
		removePagesOperation.GetBaseOperation().GetStatus(),
		removePagesOperation.GetBaseOperation().GetStoppedReason(),
	)

	operationStorage.Insert(secondLevelHash, data)

	cr <- &MergeResponse{str: ChannelResponseOK, err: nil}
	return
}

func (rpc *RemovePageController) formValidation(form *multipart.Form) error {
	if _, ok := form.Value[internal.RemovePagesIntervals]; !ok {
		return errors.New("form must contain the remove pages intervals")
	}

	err := validation.NumberFilesValidation(form, 1)
	if err != nil {
		return err
	}

	err = validation.AlphaSymbolValidation(form, internal.RemovePagesIntervals)
	if err != nil {
		return err
	}

	err = validation.OrderIntervalValidation(form, internal.RemovePagesIntervals)
	if err != nil {
		return err
	}

	err = validation.FormFileValidation(form)
	if err != nil {
		return err
	}

	return nil
}
