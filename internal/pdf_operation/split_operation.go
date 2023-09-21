package pdf_operation

import (
	"errors"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"pdf/internal/adapter"
	"slices"
	"strconv"
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

func (so *SplitOperation) GetSplitDir() adapter.SplitDir {
	return so.splitDir
}

// делать это в контроллере после выполнения операции и вставлять в хранилище эту структуру
//operationData := NewOperationData(bo.GetUserData(), bo.archiveDir, bo.status, bo.stoppedReason)

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

	inFile := string(bo.GetInDir()) + file

	many, intervals := bo.GetConfiguration().parseIntervals(splitIntervals)
	pageCount, err := api.PageCountFile(inFile)
	maxValue := slices.Max(many)

	if err != nil || pageCount < maxValue {
		wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: page coun less interval %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	pdfAdapter := locator.Locate(adapter.PdfAlias).(*adapter.PdfAdapter)
	err = pdfAdapter.SplitFile(inFile, string(so.GetSplitDir()))

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	fileAdapter := locator.Locate(adapter.FileAlias).(*adapter.FileAdapter)
	splitEntries, err := fileAdapter.GetAllEntriesFromDir(string(so.GetSplitDir()), ".pdf")

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: cant read split dir  %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	err = so.mergeFiles(pathAdapter, pdfAdapter, splitEntries, intervals, splitIntervals, inFile)

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: cant read out dir  %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	outEntries, err := fileAdapter.GetAllEntriesFromDir(string(bo.GetOutDir()), ".pdf")

	if err != nil {
		wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: cant read out dir  %w", inFile, err)
		bo.SetStatus(StatusCanceled).SetStoppedReason(StoppedReason(wrapErr.Error()))
		return wrapErr
	}

	fmt.Println(outEntries)

	// работа по архивации

	// поле этого кода в контроллере нужно будет провести работу с архивом, данных из so хватает
	// после окончания операции, нужно внести данные в хранилище операций. Объекты будут вноситься
	// в хранилище при создании операции, затем обновляться во время выполнения (processed, canceled,

	bo.SetStatus(StatusAwaitingDownload)
	return nil
}

func (so *SplitOperation) mergeFiles(
	pathAdapter *adapter.PathAdapter,
	pdfAdapter *adapter.PdfAdapter,
	splitEntries map[string]string,
	intervals [][]int,
	splitIntervals []string,
	inFile string,
) error {
	bo := so.GetBaseOperation()
	for k, interval := range intervals {
		outFile := string(bo.GetOutDir()) + string(bo.GetUserData().GetHash1Lvl()) + "_" + splitIntervals[k] + ".pdf"
		forMerge := make([]string, 0)
		fileIndex := 0

		for index := interval[0]; index <= interval[1]; index++ {
			cast := strconv.Itoa(index)
			find, ok := splitEntries[cast]
			if !ok {
				wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: can't get from map", inFile)
				return wrapErr
			}

			_, newPath, err := pathAdapter.StepForward(adapter.Path(so.GetSplitDir()), find)
			if err != nil {
				wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: can't build filepath  %w", inFile, err)
				return wrapErr
			}

			forMerge = slices.Insert(forMerge, fileIndex, string(newPath))
			fileIndex++
		}

		err := pdfAdapter.MergeFiles(forMerge, outFile)
		if err != nil {
			wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: can't MERGE  %w", inFile, err)
			return wrapErr
		}

		err = pdfAdapter.Optimize(outFile, outFile)
		if err != nil {
			wrapErr := fmt.Errorf("can't execute operation SPLIT to file %s: can't MERGE  %w", inFile, err)
			return wrapErr
		}
	}

	return nil
}
