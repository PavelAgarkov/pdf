package pdf_operation

import (
	"errors"
	"fmt"
	"pdf/internal/service"
)

type SplitOperation struct {
	baseOperation *BaseOperation
}

func NewSplitOperation(
	bo *BaseOperation,
) *SplitOperation {
	return &SplitOperation{baseOperation: bo}
}

func (so *SplitOperation) GetBaseOperation() *BaseOperation {
	return so.baseOperation
}

func (so *SplitOperation) Execute(pdfAdapter *service.PdfAdapter) error {
	bo := so.GetBaseOperation()
	allPaths := bo.GetAllPaths()
	if len(allPaths) > 1 {
		return errors.New("operation split can't have more 1 file")
	}
	inFile := allPaths[0]
	err := pdfAdapter.SplitFile(inFile, string(bo.outDir))
	if err != nil {
		return fmt.Errorf("can't execute operation split to file %s: %w", inFile, err)
	}
	return nil
}
