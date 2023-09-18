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

func (do *MergeOperation) GetDestination() string {
	return DestinationDownload
}

func (mo *MergeOperation) GetBaseOperation() *BaseOperation {
	return mo.baseOperation
}

func (mo *MergeOperation) Execute(locator *adapter.Locator) error {
	bo := mo.GetBaseOperation()
	bo.SetStatus(StatusProcessed)

	mergeOrder := bo.GetConfiguration().GetMergeOrder()
	if mergeOrder == nil {
		err := errors.New("can't execute operation merge, no merge order: %w")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	allPaths := bo.GetAllPaths()

	if len(allPaths) <= 1 {
		err := errors.New("operation merge can't have less 1 file")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err := pdfAdapter.MergeFiles(allPaths, string(bo.outDir))

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation merge to files %s: %w", allPaths, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	bo.SetStatus(StatusAwaitingDownload)
	return nil
}
