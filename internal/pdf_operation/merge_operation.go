package pdf_operation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"pdf/internal/adapter"
)

const (
	DestinationMerge = "merge"
)

type MergeOperation struct {
	baseOperation *BaseOperation
}

func NewMergeOperation(
	bo *BaseOperation,
) *MergeOperation {
	return &MergeOperation{baseOperation: bo}
}

func (mo *MergeOperation) GetDestination() string {
	return DestinationMerge
}

func (mo *MergeOperation) GetBaseOperation() *BaseOperation {
	return mo.baseOperation
}

func (mo *MergeOperation) Execute(ctx context.Context, locator *adapter.Locator, format string) (string, error) {
	defer func() {
		_ = os.RemoveAll(string(mo.GetBaseOperation().GetInDir()))
		_ = os.RemoveAll(string(mo.GetBaseOperation().GetOutDir()))
	}()

	bo := mo.GetBaseOperation()
	bo.SetStatus(StatusProcessed)

	inFiles := bo.GetAllPaths()

	if len(inFiles) <= 1 {
		err := errors.New("operation MERGE can't have less 1 file")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return "", err
	}

	outFile := string(bo.GetOutDir()) + string(bo.GetUserData().GetHash1Lvl()) + ".pdf"

	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err := pdfAdapter.MergeFiles(inFiles, outFile)

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation MERGE to files %s: %w", inFiles, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return "", wrapErr
	}

	err = pdfAdapter.Optimize(outFile, outFile)
	if err != nil {
		wrapErr := fmt.Errorf("can't optimize operation MERGE to file %s: %w", outFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return "", wrapErr
	}

	fileAdapter := locator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	outEntries, err := fileAdapter.GetAllEntriesFromDir(string(bo.GetOutDir()), ".pdf")

	pathAdapter := locator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
	associationPath := pathAdapter.BuildOutPathFilesMap(outEntries, mo.GetBaseOperation().GetUserData().GetHash2Lvl())
	archiveAdapter := locator.Locate(adapter.ArchiveAlias).(*adapter.ArchiveAdapter)
	compressor, _ := archiveAdapter.CreateCompressor(format)
	archivePath, err := archiveAdapter.Archive(
		ctx,
		compressor,
		associationPath,
		mo.GetBaseOperation().GetUserData().GetHash2Lvl(),
		mo.GetBaseOperation().GetArchiveDir(),
	)

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation MERGE : can't achivation %s:  %w", archivePath, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return "", wrapErr
	}

	bo.SetStatus(StatusAwaitingDownload)
	return archivePath, nil
}
