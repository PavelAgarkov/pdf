package pdf_operation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"pdf/internal/adapter"
)

const (
	DestinationRemovePages = "remove_pages"
)

type RemovePagesOperation struct {
	baseOperation *BaseOperation
}

func NewRemovePagesOperation(
	bo *BaseOperation,
) *RemovePagesOperation {
	return &RemovePagesOperation{baseOperation: bo}
}

func (rpo *RemovePagesOperation) GetDestination() string {
	return DestinationRemovePages
}

func (rpo *RemovePagesOperation) GetBaseOperation() *BaseOperation {
	return rpo.baseOperation
}

func (rpo *RemovePagesOperation) Execute(ctx context.Context, locator *adapter.Locator, format string) (string, error) {
	defer func() {
		_ = os.RemoveAll(string(rpo.GetBaseOperation().GetInDir()))
		_ = os.RemoveAll(string(rpo.GetBaseOperation().GetOutDir()))
	}()

	bo := rpo.GetBaseOperation()
	bo.SetStatus(StatusProcessed)

	removeIntervals := bo.GetConfiguration().GetRemovePagesIntervals()
	if removeIntervals == nil {
		err := errors.New("can't execute operation REMOVE_PAGES, no intervals: %w")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return "", err
	}

	allPaths := bo.GetAllPaths()

	if len(allPaths) > 1 {
		err := errors.New("operation REMOVE_PAGES can't have more 1 file")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return "", err
	}

	firstFile := allPaths[0]

	pathAdapter := locator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
	_, file, err := pathAdapter.StepBack(adapter.Path(firstFile))

	inFile := string(bo.GetInDir()) + file

	outFile := string(bo.GetOutDir()) + string(bo.GetUserData().GetHash1Lvl()) + ".pdf"
	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err = pdfAdapter.RemovePagesFile(inFile, outFile, removeIntervals)

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation REMOVE_PAGES to file %s: %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return "", wrapErr
	}

	err = pdfAdapter.Optimize(outFile, outFile)
	if err != nil {
		wrapErr := fmt.Errorf("can't optimize operation REMOVE_PAGES to file %s: %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return "", wrapErr
	}

	fileAdapter := locator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	outEntries, err := fileAdapter.GetAllEntriesFromDir(string(bo.GetOutDir()), ".pdf")

	associationPath := pathAdapter.BuildOutPathFilesMap(outEntries, rpo.GetBaseOperation().GetUserData().GetHash2Lvl())
	archiveAdapter := locator.Locate(adapter.ArchiveAlias).(*adapter.ArchiveAdapter)
	compressor, _ := archiveAdapter.CreateCompressor(format)
	archivePath, err := archiveAdapter.Archive(
		ctx,
		compressor,
		associationPath,
		rpo.GetBaseOperation().GetUserData().GetHash2Lvl(),
		rpo.GetBaseOperation().GetArchiveDir(),
	)

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation MERGE : can't achivation %s:  %w", archivePath, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return "", wrapErr
	}

	bo.SetStatus(StatusAwaitingDownload)
	return archivePath, nil
}
