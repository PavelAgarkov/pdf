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
	ctx context.Context,
	operationStorage *storage.OperationStorage,
	operationFactory *pdf_operation.OperationsFactory,
	adapterLocator *locator.Locator,
	loggerFactory *logger.Factory,
) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer RestoreController(loggerFactory, "split page controller")

		overAuthenticatedErr := spc.bc.isOverAuthenticated(operationStorage, c, loggerFactory)
		if overAuthenticatedErr != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": overAuthenticatedErr.Error(),
			})
		}

		authToken := service.GenerateBearerToken()

		ctxC, cancel := context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
		cr := make(chan ResponseInterface)
		start := make(chan struct{})

		go spc.realHandler(
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
		res := spc.bc.SelectResult(ctxC, cr, start)
		if res == nil {
			defer loggerFactory.PanicLog("split page controller: context expired", "")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "merge controller: context expired",
			})
		}

		if res.GetStr() != ChannelResponseOK {
			loggerFactory.ErrorLog(res.GetErr().Error(), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": res.GetErr().Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"hash":     authToken,
			"duration": internal.Timer5 * 2,
			"error":    "no",
		})
	}
}

func (spc *SplitPageController) realHandler(
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

	errForm := spc.formValidation(form)
	if errForm != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &MergeResponse{
			str: "error_form",
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

	archivePath, errArch := splitPageOperation.Execute(
		ctx,
		adapterLocator,
		form.Value[internal.ArchiveFormatKeyForRequest][0],
	)
	if errArch != nil {
		_ = os.RemoveAll(string(rootDir))
		cr <- &MergeResponse{
			str: "cant_create_archive",
			err: fmt.Errorf("cant_create_archive: %w", errArch),
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

	cr <- &MergeResponse{str: ChannelResponseOK, err: nil}
	return
}

func (spc *SplitPageController) formValidation(form *multipart.Form) error {
	if _, ok := form.Value[internal.SplitPageIntervals]; !ok {
		return errors.New("form must contain the files split intervals")
	}

	err := spc.bc.numberFilesValidation(form, 1)
	if err != nil {
		return err
	}

	err = spc.bc.alphaSymbolValidation(form, internal.SplitPageIntervals)
	if err != nil {
		return err
	}

	err = spc.bc.orderIntervalValidation(form, internal.SplitPageIntervals)
	if err != nil {
		return err
	}

	err = spc.bc.formValidation(form)
	if err != nil {
		return err
	}

	return nil
}
