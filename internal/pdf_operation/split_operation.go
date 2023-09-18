package pdf_operation

import (
	"errors"
	"fmt"
	"pdf/internal/adapter"
)

const (
	DestinationSplit = "split"
)

type SplitOperation struct {
	baseOperation *BaseOperation
}

func NewSplitOperation(
	bo *BaseOperation,
) *SplitOperation {
	return &SplitOperation{baseOperation: bo}
}

func (so *SplitOperation) GetDestination() string {
	return DestinationDownload
}

func (so *SplitOperation) GetBaseOperation() *BaseOperation {
	return so.baseOperation
}

// возможно придется добавить структуру с параметрами в интерфейс, для операций

func (so *SplitOperation) Execute(locator *adapter.Locator) error {
	bo := so.GetBaseOperation()
	bo.SetStatus(StatusProcessed)

	splitIntervals := bo.GetConfiguration().GetSplitIntervals()
	if splitIntervals == nil {
		err := errors.New("can't execute operation split, no intervals: %w")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	allPaths := bo.GetAllPaths()

	if len(allPaths) > 1 {
		err := errors.New("operation split can't have more 1 file")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	inFile := allPaths[0]

	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err := pdfAdapter.SplitFile(inFile, string(bo.outDir))

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation split to file %s: %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	// поле этого кода в контроллере нужно будет провести работу с архивом, данных из so хватает
	// после окончания операции, нужно внести данные в хранилище операций. Объекты будут вноситься
	// в хранилище при создании операции, затем обновляться во время выполнения (processed, canceled,

	bo.SetStatus(StatusAwaitingDownload)
	return nil
}
