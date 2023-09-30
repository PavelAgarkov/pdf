package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
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
		defer RestoreController(loggerFactory, c, "merge controller")

		authToken := service.ParseBearerHeader(c.GetReqHeaders()[internal.AuthenticationHeader])
		_, hit := operationStorage.Get(internal.Hash2lvl(authToken))
		if hit {
			errMsg := fmt.Sprintf("merge controller: can't process %s from storage", authToken)
			loggerFactory.ErrorLog(errMsg, "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": errMsg,
			})
		}

		authToken = service.GenerateBearerToken()

		ctxC, cancel := context.WithTimeout(ctx, 3*time.Second)
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
		)
		res := mc.bc.SelectResult(ctxC, cr, start)
		if res == nil {
			loggerFactory.PanicLog("merge controller: context expired", "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "merge controller: context expired",
			})
		}

		if res.GetStr() != "ok" {
			loggerFactory.ErrorLog(res.GetErr().Error(), "")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": res.GetErr().Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"hash":     authToken,
			"duration": internal.Timer5,
			"error":    "no",
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
) {
	<-start

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
		_ = os.RemoveAll(string(rootDir))
	}()

	if err != nil {
		cr <- &MergeResponse{str: "cant_create_dir", err: err}
		return
	}

	filesOrder := make([]string, 0)
	form, errRead := c.MultipartForm()
	if errRead != nil {
		cr <- &MergeResponse{str: "cant_read_form", err: errRead}
		return
	}

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			nameWithoutSpace := strings.ReplaceAll(fileHeader.Filename, " ", "_")
			_, pathToFile, _ := pathAdapter.StepForward(internal.Path(inDir), nameWithoutSpace)
			errSave := c.SaveFile(fileHeader, string(pathToFile))
			if errSave != nil {
				cr <- &MergeResponse{str: "cant_save_file_from_form", err: errSave}
				return
			}
			filesOrder = append(filesOrder, string(pathToFile))
		}
	}

	userData := entity.NewUserData(
		internal.Hash1lvl(authToken),
		secondLevelHash,
		time.Now().Add(internal.Timer5*internal.Minute),
	)
	mergePagesOperation := operationFactory.CreateNewOperation(
		pdf_operation.NewConfiguration(nil, nil, nil),
		userData,
		filesOrder,
		rootDir,
		inDir,
		outDir,
		archiveDir,
		"",
		pdf_operation.DestinationMerge,
	).(*pdf_operation.MergeOperation)

	archivePath, errArch := mergePagesOperation.Execute(ctx, adapterLocator, internal.ZipFormat)
	if errArch != nil {
		cr <- &MergeResponse{str: "cant_create_archive", err: errArch}
		return
	}

	data := pdf_operation.NewOperationData(
		userData,
		internal.ArchiveDir(archivePath),
		mergePagesOperation.GetBaseOperation().GetStatus(),
		mergePagesOperation.GetBaseOperation().GetStoppedReason(),
	)

	operationStorage.Insert(secondLevelHash, data)

	cr <- &MergeResponse{str: "ok", err: nil}
	return
}
