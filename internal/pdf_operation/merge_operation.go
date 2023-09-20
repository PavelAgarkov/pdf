package pdf_operation

import (
	"errors"
	"fmt"
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

func (mo *MergeOperation) Execute(locator *adapter.Locator) error {
	bo := mo.GetBaseOperation()
	bo.SetStatus(StatusProcessed)

	inFiles := bo.GetAllPaths()

	if len(inFiles) <= 1 {
		err := errors.New("operation MERGE can't have less 1 file")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	outFile := string(bo.outDir) + string(bo.GetUserData().GetHash1Lvl()) + ".pdf"

	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err := pdfAdapter.MergeFiles(inFiles, outFile)

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation MERGE to files %s: %w", inFiles, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	err = pdfAdapter.Optimize(outFile, outFile)
	if err != nil {
		wrapErr := fmt.Errorf("can't optimize operation MERGE to file %s: %w", outFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	// работа по архивации

	bo.SetStatus(StatusAwaitingDownload)
	return nil
}
