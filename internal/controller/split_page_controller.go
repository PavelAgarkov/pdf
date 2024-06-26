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

type SplitPageController struct {
	bc *BaseController
}

type SplitResponse struct {
	str string
	err error
}

func NewSplitPageController(bc *BaseController) *SplitPageController {
	return &SplitPageController{
		bc: bc,
	}
}

func (r *SplitResponse) GetStr() string {
	return r.str
}

func (r *SplitResponse) GetErr() error {
	return r.err
}

func (spc *SplitPageController) Handle(
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "split page controller")

		overAuthenticatedErr := spc.bc.isOverAuthenticated(operationStorage, c)
		if overAuthenticatedErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": overAuthenticatedErr.Error(),
			})
		}

		authToken := service.GenerateBearerToken()
		ctx := c.UserContext()
		responseCh := make(chan ResponseInterface, 2)

		spc.splitHandler(
			c,
			ctx,
			responseCh,
			operationStorage,
			operationFactory,
			adapterLocator,
			authToken,
			loggerFactory,
		)
		result := spc.bc.Select(ctx, responseCh)

		if result == nil {
			defer loggerFactory.PanicLog("split page controller: context expired", "")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "split controller: context expired",
			})
		}

		if result.GetStr() == internal.ErrorForm {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": result.GetErr().Error(),
			})
		}

		if result.GetStr() == internal.OperationError {
			loggerFactory.ErrorLog(result.GetErr().Error(), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "operation_version_is_not_suitable",
			})
		}

		if result.GetStr() != ChannelResponseOK {
			loggerFactory.ErrorLog(result.GetErr().Error(), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": result.GetErr().Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"hash":  authToken,
			"error": "",
		})
	}
}

func (spc *SplitPageController) splitHandler(
	c *fiber.Ctx,
	ctx context.Context,
	cr chan ResponseInterface,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	authToken string,
	loggerFactory *logger.Factory,
) {
	defer RestoreController(loggerFactory, "split page controller")

	secondLevelHash := hash.GenerateNextLevelHashByPrevious(internal.Hash1lvl(authToken), true)
	pathAdapter := adapterLocator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
	inDir := pathAdapter.GenerateInDirPath(secondLevelHash)
	rootDir := pathAdapter.GenerateRootDir(secondLevelHash)
	outDir := pathAdapter.GenerateOutDirPath(secondLevelHash)
	archiveDir := pathAdapter.GenerateArchiveDirPath(secondLevelHash)
	splitDir := pathAdapter.GenerateDirPathToSplitFiles(secondLevelHash)

	fileAdapter := adapterLocator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	err := fileAdapter.CreateDir(string(rootDir), 0777)
	err = fileAdapter.CreateDir(string(inDir), 0777)
	err = fileAdapter.CreateDir(string(outDir), 0777)
	err = fileAdapter.CreateDir(string(archiveDir), 0777)
	err = fileAdapter.CreateDir(string(splitDir), 0777)
	defer func() {
		_ = os.RemoveAll(string(inDir))
		_ = os.RemoveAll(string(outDir))
		_ = os.RemoveAll(string(splitDir))
		cr <- &SplitResponse{
			str: "panic_context",
			err: errors.New("handle cancel"),
		}
	}()

	if err != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &SplitResponse{
			str: "cant_create_dir",
			err: fmt.Errorf("cant_read_form: %w", err),
		}
		return
	}

	form, errRead := c.MultipartForm()
	if errRead != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &SplitResponse{
			str: "cant_read_form",
			err: fmt.Errorf("cant_read_form: %w", errRead),
		}
		return
	}

	errForm := spc.formValidation(form)
	if errForm != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &SplitResponse{
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
				cr <- &SplitResponse{
					str: "cant_save_file_from_form",
					err: fmt.Errorf("cant_save_file_from_form: %w", errSave),
				}
				return
			}
			files = append(files, string(pathToFile))
		}
	}

	splitPageIntervals := form.Value[internal.SplitPageIntervals]

	userData := entity.NewUserData(
		internal.Hash1lvl(authToken),
		secondLevelHash,
		time.Now().Add(internal.Timer5*internal.Minute),
	)
	splitPageOperation := operationFactory.CreateNewOperation(
		pdf_operation.NewConfiguration(splitPageIntervals, nil),
		userData,
		files,
		rootDir,
		inDir,
		outDir,
		archiveDir,
		splitDir,
		pdf_operation.DestinationSplit,
	).(*pdf_operation.SplitOperation)

	archivePath, operationErr := splitPageOperation.Execute(
		ctx,
		adapterLocator,
		form.Value[internal.ArchiveFormatKeyForRequest][0],
	)
	if operationErr != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &SplitResponse{
			str: internal.OperationError,
			err: fmt.Errorf("can't make opearation: %w", operationErr),
		}
		return
	}

	data := pdf_operation.NewOperationData(
		userData,
		internal.ArchiveDir(archivePath),
		splitPageOperation.GetBaseOperation().GetStatus(),
		splitPageOperation.GetBaseOperation().GetStoppedReason(),
	)

	operationStorage.Insert(secondLevelHash, data)

	cr <- &SplitResponse{str: ChannelResponseOK, err: nil}
	return
}

func (spc *SplitPageController) formValidation(form *multipart.Form) error {
	if _, ok := form.Value[internal.SplitPageIntervals]; !ok {
		return errors.New("form must contain the files split intervals")
	}

	err := validation.NumberFilesValidation(form, 1)
	if err != nil {
		return err
	}

	err = validation.AlphaSymbolValidation(form, internal.SplitPageIntervals)
	if err != nil {
		return err
	}

	err = validation.OrderIntervalValidation(form, internal.SplitPageIntervals)
	if err != nil {
		return err
	}

	err = validation.FormFileValidation(form)
	if err != nil {
		return err
	}

	return nil
}
