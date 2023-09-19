package pdf_operation

import (
	"errors"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"pdf/internal/adapter"
	"slices"
)

const (
	DestinationSplit = "split"
)

type SplitOperation struct {
	baseOperation *BaseOperation
	splitDir      adapter.SplitDir
}

func NewSplitOperation(
	bo *BaseOperation,
	splitDir adapter.SplitDir,
) *SplitOperation {
	return &SplitOperation{
		baseOperation: bo,
		splitDir:      splitDir,
	}
}

func (so *SplitOperation) GetDestination() string {
	return DestinationSplit
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
		err := errors.New("can't execute operation SPLIT, no intervals: %w")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	allPaths := bo.GetAllPaths()

	if len(allPaths) > 1 {
		err := errors.New("operation SPLIT can't have more 1 file")
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(err.Error()))
		return err
	}

	firstFile := allPaths[0]

	pathAdapter := locator.Locate(adapter.PathAlias).(*adapter.PathAdapter)
	_, file, err := pathAdapter.StepBack(adapter.Path(firstFile))

	inFile := string(bo.inDir) + file

	many, intervals := bo.GetConfiguration().parseIntervals(splitIntervals)
	pageCount, err := api.PageCountFile(inFile)
	maxValue := slices.Max(many)

	if err != nil || pageCount < maxValue {
		wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: page coun less interval %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	fmt.Println(intervals)

	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err = pdfAdapter.SplitFile(inFile, string(so.splitDir))

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	// собрать по intervals файлы из директории

	// работа по архивации

	// поле этого кода в контроллере нужно будет провести работу с архивом, данных из so хватает
	// после окончания операции, нужно внести данные в хранилище операций. Объекты будут вноситься
	// в хранилище при создании операции, затем обновляться во время выполнения (processed, canceled,

	bo.SetStatus(StatusAwaitingDownload)
	return nil
}
