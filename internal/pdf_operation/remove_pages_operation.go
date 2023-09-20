package pdf_operation

import (
	"errors"
	"fmt"
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

func (rpo *RemovePagesOperation) Execute(locator *adapter.Locator) error {
	bo := rpo.GetBaseOperation()
	bo.SetStatus(StatusProcessed)

	removeIntervals := bo.GetConfiguration().GetRemovePagesIntervals()
	if removeIntervals == nil {
		err := errors.New("can't execute operation REMOVE_PAGES, no intervals: %w")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	allPaths := bo.GetAllPaths()

	if len(allPaths) > 1 {
		err := errors.New("operation REMOVE_PAGES can't have more 1 file")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	firstFile := allPaths[0]

	pathAdapter := locator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
	_, file, err := pathAdapter.StepBack(adapter.Path(firstFile))

	inFile := string(bo.inDir) + file

	outFile := string(bo.outDir) + string(bo.GetUserData().GetHash1Lvl()) + ".pdf"
	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err = pdfAdapter.RemovePagesFile(inFile, outFile, removeIntervals)

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation REMOVE_PAGES to file %s: %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	err = pdfAdapter.Optimize(outFile, outFile)
	if err != nil {
		wrapErr := fmt.Errorf("can't optimize operation REMOVE_PAGES to file %s: %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	// работа по архивации

	bo.SetStatus(StatusAwaitingDownload)
	return nil
}
